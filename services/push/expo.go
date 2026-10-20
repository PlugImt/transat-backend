package push

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	expoSendURL     = "https://api.expo.dev/v2/push/send"
	expoReceiptsURL = "https://exp.host/--/api/v2/push/getReceipts"
	// Expo accepts 100 messages per request and 1000 tickets per receipt request.
	expoMaxBatch = 100
	// MaxReceiptBatch is how many tickets one receipt request may carry.
	MaxReceiptBatch = 1000
	maxDetailLength = 200
)

// ExpoConfig configures the Expo push service client.
type ExpoConfig struct {
	SendURL     string
	ReceiptsURL string
	// AccessToken is the Expo access token, required once "enhanced push security" is on for the project. It is a
	// secret and only ever sent as a bearer token to Expo.
	AccessToken string
	Client      *http.Client
}

// ExpoConfigFromEnv reads EXPO_ACCESS_TOKEN.
func ExpoConfigFromEnv() ExpoConfig {
	return ExpoConfig{AccessToken: strings.TrimSpace(os.Getenv("EXPO_ACCESS_TOKEN"))}
}

// Expo sends through the Expo push service, which forwards to APNs (iOS) and FCM (Android).
type Expo struct {
	cfg ExpoConfig
}

// NewExpoProvider builds the Expo client; empty fields of cfg get their defaults.
func NewExpoProvider(cfg ExpoConfig) *Expo {
	if cfg.SendURL == "" {
		cfg.SendURL = expoSendURL
	}
	if cfg.ReceiptsURL == "" {
		cfg.ReceiptsURL = expoReceiptsURL
	}
	if cfg.Client == nil {
		cfg.Client = &http.Client{Timeout: 15 * time.Second}
	}
	return &Expo{cfg: cfg}
}

func (e *Expo) Name() string  { return "expo" }
func (e *Expo) MaxBatch() int { return expoMaxBatch }

type expoTicket struct {
	Status  string `json:"status"`
	ID      string `json:"id"`
	Message string `json:"message"`
	Details struct {
		Error string `json:"error"`
	} `json:"details"`
}

type expoErrors struct {
	Errors []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func (e *Expo) Send(ctx context.Context, messages []Outbound) ([]SendResult, error) {
	body := make([]map[string]any, len(messages))
	for i, m := range messages {
		body[i] = expoMessage(m)
	}
	raw, err := e.post(ctx, e.cfg.SendURL, body)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Data []expoTicket `json:"data"`
		expoErrors
	}
	if err := json.Unmarshal(raw, &parsed); err != nil || len(parsed.Data) != len(messages) {
		if len(parsed.Errors) > 0 {
			return nil, &RequestError{Code: parsed.Errors[0].Code, Isolate: true, Detail: redact(parsed.Errors[0].Message)}
		}
		// Expo answered 200 but not with a ticket per message: the messages were most likely taken, and sending
		// them again would duplicate them. They are accepted without a ticket, so no receipt can follow.
		slog.Warn("push.provider_malformed_response", "provider", "expo", "messages", len(messages))
		results := make([]SendResult, len(messages))
		return results, nil
	}

	results := make([]SendResult, len(messages))
	for i, t := range parsed.Data {
		results[i] = expoTicketResult(t)
	}
	return results, nil
}

func expoMessage(m Outbound) map[string]any {
	msg := map[string]any{
		"to":        m.Token,
		"title":     m.Title,
		"sound":     "default",
		"channelId": "default",
		"priority":  "high",
	}
	if m.Body != "" {
		msg["body"] = m.Body
	}
	if m.Data != nil {
		msg["data"] = m.Data
	}
	// Android shows the image on its own; iOS shows it only when the app's notification service extension runs,
	// which needs mutableContent. Plain notifications never ask for the extension.
	if m.ImageURL != "" {
		msg["richContent"] = map[string]string{"image": m.ImageURL}
		if m.Platform == "ios" {
			msg["mutableContent"] = true
		}
	}
	if m.ActionCategory != "" {
		msg["categoryId"] = m.ActionCategory
	}
	if ttl := int(m.TTL.Seconds()); ttl > 0 {
		msg["ttl"] = ttl
	}
	// FCM keeps at most four collapse keys per device, so Android gets a tag (it replaces what is shown) and iOS
	// gets the collapse id (APNs replaces what is shown and drops the older copy in transit).
	switch m.Platform {
	case "ios":
		if m.DedupeTag != "" {
			msg["collapseId"] = m.DedupeTag
		}
	case "android":
		if m.DedupeTag != "" {
			msg["tag"] = m.DedupeTag
		}
	}
	return msg
}

func expoTicketResult(t expoTicket) SendResult {
	if t.Status == "ok" {
		return SendResult{Kind: ResultAccepted, TicketID: t.ID}
	}
	code, detail := t.Details.Error, redact(t.Message)
	switch code {
	case "DeviceNotRegistered":
		return SendResult{Kind: ResultInvalidDevice, Code: code, Detail: detail}
	case "MessageRateExceeded":
		return SendResult{Kind: ResultRetry, Code: code, Detail: detail}
	case "MessageTooBig", "InvalidCredentials", "MismatchSenderId":
		return SendResult{Kind: ResultRejected, Code: code, Detail: detail}
	case "":
		if strings.Contains(strings.ToLower(t.Message), "not a valid expo push token") {
			return SendResult{Kind: ResultRejected, Code: "InvalidToken", Detail: detail}
		}
		code = "unknown"
	}
	// An error nobody planned for is more likely Expo's trouble than the message's: try again, within the bound.
	return SendResult{Kind: ResultRetry, Code: code, Detail: detail}
}

func (e *Expo) Receipts(ctx context.Context, ticketIDs []string) (map[string]ReceiptResult, error) {
	raw, err := e.post(ctx, e.cfg.ReceiptsURL, map[string][]string{"ids": ticketIDs})
	if err != nil {
		return nil, err
	}
	var parsed struct {
		Data map[string]expoTicket `json:"data"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, &RequestError{Code: "malformed_response", Retryable: true}
	}
	out := make(map[string]ReceiptResult, len(parsed.Data))
	for id, r := range parsed.Data {
		if r.Status == "ok" {
			out[id] = ReceiptResult{Kind: ReceiptHandedOff}
			continue
		}
		code, detail := r.Details.Error, redact(r.Message)
		switch code {
		case "DeviceNotRegistered":
			out[id] = ReceiptResult{Kind: ReceiptInvalidDevice, Code: code, Detail: detail}
		case "MessageRateExceeded":
			out[id] = ReceiptResult{Kind: ReceiptRetry, Code: code, Detail: detail}
		default:
			if code == "" {
				code = "unknown"
			}
			out[id] = ReceiptResult{Kind: ReceiptFailed, Code: code, Detail: detail}
		}
	}
	return out, nil
}

// post sends one JSON request and turns every failure into a RequestError the pipeline knows how to treat.
func (e *Expo) post(ctx context.Context, url string, payload any) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, &RequestError{Code: "encode", Detail: err.Error()}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, &RequestError{Code: "request", Detail: err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if e.cfg.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+e.cfg.AccessToken)
	}

	resp, err := e.cfg.Client.Do(req)
	if err != nil {
		return nil, &RequestError{Code: "network", Retryable: true, Detail: redact(err.Error())}
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if resp.StatusCode == http.StatusOK {
		return raw, nil
	}

	var failure expoErrors
	_ = json.Unmarshal(raw, &failure)
	code, detail := fmt.Sprintf("http_%d", resp.StatusCode), ""
	if len(failure.Errors) > 0 {
		if failure.Errors[0].Code != "" {
			code = failure.Errors[0].Code
		}
		detail = redact(failure.Errors[0].Message)
	}
	re := &RequestError{Code: code, Status: resp.StatusCode, Detail: detail}
	switch {
	case resp.StatusCode == http.StatusTooManyRequests, resp.StatusCode >= 500:
		re.Retryable = true
	case resp.StatusCode == http.StatusUnauthorized, resp.StatusCode == http.StatusForbidden:
		// A missing or wrong access token is a configuration problem: waiting lets it be fixed, giving up would
		// lose the notification.
		re.Retryable = true
		slog.Error("push.provider_auth_failed", "provider", "expo", "status", resp.StatusCode)
	default:
		re.Isolate = true
	}
	return nil, re
}

var expoTokenText = regexp.MustCompile(`Expo(?:nent)?PushToken\[[^\]]*\]`)

// redact keeps device tokens out of what is stored and logged, and bounds its size.
func redact(text string) string {
	text = expoTokenText.ReplaceAllString(text, "[push token]")
	if len(text) > maxDetailLength {
		text = strings.ToValidUTF8(text[:maxDetailLength], "")
	}
	return text
}
