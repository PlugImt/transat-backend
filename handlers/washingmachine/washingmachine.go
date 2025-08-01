package washingmachine

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
)

type WashingMachineHandler struct{}

// NewWashingMachineHandler creates a new instance of WashingMachineHandler
func NewWashingMachineHandler() *WashingMachineHandler {
	return &WashingMachineHandler{}
}

// GetWashingMachines returns the status of all washing machines
func (h *WashingMachineHandler) GetWashingMachines() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Prepare the request to the external API
		requestData := strings.NewReader("action=READ_LIST_STATUS&serial_centrale=65e4444c3471550a789e2138a9e28eff")
		req, err := http.NewRequest("POST", "https://status.wi-line.fr/update_machine_ext.php", requestData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.WashingMachineResponse{
				Success: false,
				Message: "Failed to create request to external API",
			})
		}

		// Set headers
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.WashingMachineResponse{
				Success: false,
				Message: "Failed to fetch data from external API",
			})
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.WashingMachineResponse{
				Success: false,
				Message: "Failed to read response from external API",
			})
		}

		// Parse the response
		var externalResponse struct {
			MachineInfoStatus struct {
				MachineList []models.MachineData `json:"machine_list"`
			} `json:"machine_info_status"`
		}

		if err := json.Unmarshal(body, &externalResponse); err != nil {
			// Try to read the raw response for debugging
			return c.Status(fiber.StatusInternalServerError).JSON(models.WashingMachineResponse{
				Success: false,
				Message: fmt.Sprintf("Failed to parse response from external API: %v", err),
			})
		}

		// Transform the data into the desired format
		formattedData := transformMachineData(externalResponse.MachineInfoStatus.MachineList)

		fmt.Println(formattedData)

		// Return the formatted data
		return c.JSON(models.WashingMachineResponse{
			Success: true,
			Data:    formattedData,
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
				{Number: 13, Available: true, TimeLeft: 69},
				{Number: 14, Available: true, TimeLeft: 1},
			},
			Dryers: []models.FormattedMachine{
				{Number: 15, Available: true, TimeLeft: 7},
				{Number: 16, Available: true, TimeLeft: 0},
				{Number: 17, Available: true, TimeLeft: 234},
				{Number: 18, Available: true, TimeLeft: 24},
			},
		}

		return c.JSON(models.WashingMachineResponse{
			Success: true,
			Data:    testData,
		})

	}
}

// transformMachineData categorizes and transforms the machine data into washing machines and dryers
func transformMachineData(machineList []models.MachineData) models.FormattedMachineData {
	var result models.FormattedMachineData

	// Initialize slices
	result.WashingMachines = []models.FormattedMachine{}
	result.Dryers = []models.FormattedMachine{}

	for _, machine := range machineList {
		// Convert selecteur_machine to number
		number, _ := strconv.Atoi(machine.SelecteurMachine)

		// Check if machine is available (status 1 appears to mean available)
		isAvailable := machine.Status == 1

		// Create formatted machine
		formattedMachine := models.FormattedMachine{
			Number:    number,
			Available: isAvailable,
			TimeLeft:  machine.TimeBeforeOff,
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
