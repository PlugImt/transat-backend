package washingmachine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/i18n"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

const (
	upstreamURL     = "https://status.wi-line.fr/update_machine_ext.php"
	upstreamTimeout = 8 * time.Second
	// freshFor: a successful answer is reused, so a burst of users makes one upstream call.
	freshFor = 15 * time.Second
	// retryAfter: after a failure the upstream is left alone, so users get the message at once instead of waiting for the timeout.
	retryAfter = 20 * time.Second
	// maxBody caps what is read from the upstream, which is not under our control.
	maxBody = 1 << 20
)

type WashingMachineHandler struct {
	client *http.Client
	url    string

	mu       sync.Mutex
	cached   models.FormattedMachineData
	cachedAt time.Time
	failedAt time.Time
}

// NewWashingMachineHandler creates a new instance of WashingMachineHandler
func NewWashingMachineHandler() *WashingMachineHandler {
	return &WashingMachineHandler{
		client: &http.Client{Timeout: upstreamTimeout},
		url:    upstreamURL,
	}
}

// fetchUpstream returns an error for everything that is not a usable machine list: a paused or broken site answers
// with an error status, an HTML page, or nothing at all.
func (h *WashingMachineHandler) fetchUpstream(ctx context.Context) (models.FormattedMachineData, error) {
	requestData := strings.NewReader("action=READ_LIST_STATUS&serial_centrale=65e4444c3471550a789e2138a9e28eff")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.url, requestData)
	if err != nil {
		return models.FormattedMachineData{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := h.client.Do(req)
	if err != nil {
		return models.FormattedMachineData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.FormattedMachineData{}, fmt.Errorf("upstream answered HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBody))
	if err != nil {
		return models.FormattedMachineData{}, err
	}

	var externalResponse struct {
		MachineInfoStatus struct {
			MachineList []models.MachineData `json:"machine_list"`
		} `json:"machine_info_status"`
	}
	if err := json.Unmarshal(body, &externalResponse); err != nil {
		return models.FormattedMachineData{}, fmt.Errorf("upstream answered something that is not the machine list: %w", err)
	}
	if len(externalResponse.MachineInfoStatus.MachineList) == 0 {
		return models.FormattedMachineData{}, errors.New("upstream returned no machine")
	}

	return transformMachineData(externalResponse.MachineInfoStatus.MachineList), nil
}

// machines serves the recent answer when there is one, otherwise asks the upstream. Concurrent requests wait on the
// lock and then read what the first one found.
func (h *WashingMachineHandler) machines(ctx context.Context) (models.FormattedMachineData, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now()
	if !h.cachedAt.IsZero() && now.Sub(h.cachedAt) < freshFor {
		return h.cached, nil
	}
	if !h.failedAt.IsZero() && now.Sub(h.failedAt) < retryAfter {
		return models.FormattedMachineData{}, errUpstreamUnavailable
	}

	data, err := h.fetchUpstream(ctx)
	if err != nil {
		h.failedAt = now
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Washing machines upstream unavailable: %v", err))
		return models.FormattedMachineData{}, errUpstreamUnavailable
	}

	h.cached, h.cachedAt, h.failedAt = data, now, time.Time{}
	return data, nil
}

var errUpstreamUnavailable = errors.New("washing machines upstream unavailable")

// GetWashingMachines returns the status of all washing machines
func (h *WashingMachineHandler) GetWashingMachines() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Detached from the request: a user leaving must not be recorded as the upstream failing.
		data, err := h.machines(context.Background())
		if err != nil {
			// The app shows the "error" of a failed response as is, so it carries the message meant for the user.
			message := i18n.Localize(i18n.GetLocalizer(c.Get("Accept-Language")), "washing_machines.unavailable",
				"Washing machine info is temporarily unavailable: the website providing it is not responding. Please try again later.", nil)
			c.Set(fiber.HeaderRetryAfter, strconv.Itoa(int(retryAfter.Seconds())))
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success": false,
				"code":    "washing_machines_unavailable",
				"error":   message,
				"message": message,
			})
		}

		return c.JSON(models.WashingMachineResponse{
			Success: true,
			Data:    data,
		})
	}
}

func (h *WashingMachineHandler) GetWashingMachinesTest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Simulate a response for testing purposes

		// Create test data
		testData := models.FormattedMachineData{
			WashingMachines: []models.FormattedMachine{
				{Number: 11, Available: true, TimeLeft: 30},
				{Number: 12, Available: true, TimeLeft: 0},
				{Number: 13, Available: true, TimeLeft: 934},
				{Number: 14, Available: true, TimeLeft: 400},
			},
			Dryers: []models.FormattedMachine{
				{Number: 15, Available: true, TimeLeft: 953},
				{Number: 16, Available: true, TimeLeft: 0},
				{Number: 17, Available: true, TimeLeft: 5235},
				{Number: 18, Available: true, TimeLeft: 2489},
			},
		}

		return c.JSON(models.WashingMachineResponse{
			Success: true,
			Data:    testData,
		})

	}
}

// timeLeft reads the remaining time from the field the operator's own page reads for each status: 4 and 6 count down
// in "temps", 14 in "time_before_off_virtu", the others in "time_before_off".
func timeLeft(machine models.MachineData) int {
	switch machine.Status {
	case 4, 6:
		return int(machine.Temps)
	case 14:
		return int(machine.TimeBeforeOffVirtu)
	default:
		return int(machine.TimeBeforeOff)
	}
}

// transformMachineData categorizes and transforms the machine data into washing machines and dryers
func transformMachineData(machineList []models.MachineData) models.FormattedMachineData {
	var result models.FormattedMachineData

	// Initialize slices
	result.WashingMachines = []models.FormattedMachine{}
	result.Dryers = []models.FormattedMachine{}

	for _, machine := range machineList {
		// Status 1 is "free" on the operator's own page.
		isAvailable := machine.Status == 1

		// Create formatted machine
		formattedMachine := models.FormattedMachine{
			Number:    int(machine.SelecteurMachine),
			Available: isAvailable,
			TimeLeft:  timeLeft(machine),
		}

		// Categorize by machine type
		if strings.Contains(strings.ToUpper(machine.NomType), "LAVE LINGE") {
			result.WashingMachines = append(result.WashingMachines, formattedMachine)
		} else if strings.Contains(strings.ToUpper(machine.NomType), "SECHE LINGE") {
			result.Dryers = append(result.Dryers, formattedMachine)
		}
	}

	return result
}
