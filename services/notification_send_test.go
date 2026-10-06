package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/plugimt/transat-backend/models"
)

func TestSendPushNotificationBatchesAndReportsFailures(t *testing.T) {
	var mu sync.Mutex
	var batches [][]map[string]interface{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var messages []map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&messages); err != nil {
			t.Errorf("invalid body: %v", err)
		}
		mu.Lock()
		batches = append(batches, messages)
		mu.Unlock()

		tickets := make([]map[string]string, len(messages))
		for i, m := range messages {
			tickets[i] = map[string]string{"status": "ok"}
			if m["to"] == "tok-7" {
				tickets[i] = map[string]string{"status": "error", "message": "DeviceNotRegistered"}
			}
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": tickets})
	}))
	defer srv.Close()

	ns := &NotificationService{expoPushURL: srv.URL, client: srv.Client()}

	tokens := make([]string, 0, 251)
	for i := 0; i < 250; i++ {
		tokens = append(tokens, fmt.Sprintf("tok-%d", i))
	}
	tokens = append(tokens, "tok-0", "") // duplicate and empty are dropped

	err := ns.SendPushNotification(models.NotificationPayload{
		NotificationTokens: tokens,
		Title:              "Hello",
		Message:            "World",
		Navigation:         models.EventNavigation(42),
	})
	if err == nil {
		t.Fatal("expected an error for the rejected ticket")
	}

	if len(batches) != 3 || len(batches[0]) != 100 || len(batches[2]) != 50 {
		t.Fatalf("batch sizes = %v, want 100/100/50", []int{len(batches[0]), len(batches[1]), len(batches[2])})
	}

	first := batches[0][0]
	if first["title"] != "Hello" || first["body"] != "World" || first["channelId"] != "default" {
		t.Fatalf("unexpected message: %v", first)
	}
	data, ok := first["data"].(map[string]interface{})
	if !ok || data["screen"] != "Events" {
		t.Fatalf("navigation data missing: %v", first["data"])
	}
	nav, _ := data["navigation"].(map[string]interface{})
	if nav["type"] != "event" || nav["id"] != "42" {
		t.Fatalf("navigation = %v", nav)
	}
}

func TestSendPushNotificationWithoutTokens(t *testing.T) {
	ns := &NotificationService{}
	if err := ns.SendPushNotification(models.NotificationPayload{Title: "x"}); err == nil {
		t.Fatal("expected an error without tokens")
	}
}

func TestSendPushNotificationExpoDown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	ns := &NotificationService{expoPushURL: srv.URL, client: srv.Client()}
	err := ns.SendPushNotification(models.NotificationPayload{NotificationTokens: []string{"a", "b"}, Title: "x"})
	if err == nil {
		t.Fatal("expected an error when Expo answers 500")
	}
}
