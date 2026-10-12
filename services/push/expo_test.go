package push

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type expoStub struct {
	server  *httptest.Server
	last    []map[string]any
	headers http.Header
	status  int
	body    string
}

func newExpoStub(t *testing.T) *expoStub {
	s := &expoStub{status: http.StatusOK}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		s.headers = r.Header.Clone()
		s.last = nil
		_ = json.Unmarshal(raw, &s.last)
		w.WriteHeader(s.status)
		_, _ = w.Write([]byte(s.body))
	}))
	t.Cleanup(s.server.Close)
	return s
}

func (s *expoStub) provider(token string) *Expo {
	return NewExpoProvider(ExpoConfig{SendURL: s.server.URL, ReceiptsURL: s.server.URL, AccessToken: token})
}

func TestExpoSendFormatsMessagesAndReadsTickets(t *testing.T) {
	s := newExpoStub(t)
	s.body = `{"data":[
		{"status":"ok","id":"t1"},
		{"status":"error","message":"\"ExponentPushToken[dead]\" is not a registered push notification recipient","details":{"error":"DeviceNotRegistered"}},
		{"status":"error","message":"slow","details":{"error":"MessageRateExceeded"}},
		{"status":"error","message":"too big","details":{"error":"MessageTooBig"}},
		{"status":"error","message":"\"x\" is not a valid Expo push token"},
		{"status":"error","message":"what"}]}`
	e := s.provider("secret")
	res, err := e.Send(context.Background(), []Outbound{
		{Token: "ExponentPushToken[a]", Platform: "ios", Title: "T", Body: "B", Data: map[string]any{"k": "v"}, TTL: 90 * time.Second, DedupeTag: "m7"},
		{Token: "ExponentPushToken[b]", Platform: "android", Title: "T", DedupeTag: "m7"},
		{Token: "ExponentPushToken[c]", Platform: "unknown", Title: "T", DedupeTag: "m7"},
		{Token: "ExponentPushToken[d]", Title: "T"}, {Token: "ExponentPushToken[e]", Title: "T"}, {Token: "ExponentPushToken[f]", Title: "T"},
	})
	if err != nil {
		t.Fatal(err)
	}

	want := []struct {
		kind ResultKind
		code string
	}{
		{ResultAccepted, ""}, {ResultInvalidDevice, "DeviceNotRegistered"}, {ResultRetry, "MessageRateExceeded"},
		{ResultRejected, "MessageTooBig"}, {ResultRejected, "InvalidToken"}, {ResultRetry, "unknown"},
	}
	for i, w := range want {
		if res[i].Kind != w.kind || res[i].Code != w.code {
			t.Fatalf("ticket %d = %+v, want %v/%s", i, res[i], w.kind, w.code)
		}
	}
	if res[0].TicketID != "t1" {
		t.Fatalf("ticket id = %q", res[0].TicketID)
	}
	if strings.Contains(res[1].Detail, "PushToken") {
		t.Fatalf("token in detail: %q", res[1].Detail)
	}

	if s.headers.Get("Authorization") != "Bearer secret" || s.headers.Get("Content-Type") != "application/json" {
		t.Fatalf("headers = %v", s.headers)
	}
	ios, android, other := s.last[0], s.last[1], s.last[2]
	if ios["to"] != "ExponentPushToken[a]" || ios["title"] != "T" || ios["body"] != "B" || ios["ttl"] != float64(90) ||
		ios["collapseId"] != "m7" || ios["tag"] != nil || ios["priority"] != "high" || ios["channelId"] != "default" {
		t.Fatalf("ios message = %v", ios)
	}
	if android["tag"] != "m7" || android["collapseId"] != nil {
		t.Fatalf("android message = %v", android)
	}
	if other["tag"] != nil || other["collapseId"] != nil {
		t.Fatalf("a device of unknown platform got a collapse key: %v", other)
	}
	if s.last[3]["body"] != nil || s.last[3]["ttl"] != nil {
		t.Fatalf("optional fields were sent empty: %v", s.last[3])
	}
}

func TestExpoSendsAnImageAsRichContent(t *testing.T) {
	s := newExpoStub(t)
	s.body = `{"data":[{"status":"ok","id":"a"},{"status":"ok","id":"b"},{"status":"ok","id":"c"},{"status":"ok","id":"d"}]}`
	const image = "https://cdn.x.test/poster.jpg"
	_, err := s.provider("").Send(context.Background(), []Outbound{
		{Token: "ExponentPushToken[a]", Platform: "ios", Title: "T", ImageURL: image},
		{Token: "ExponentPushToken[b]", Platform: "android", Title: "T", ImageURL: image},
		{Token: "ExponentPushToken[c]", Platform: "ios", Title: "T"},
		{Token: "ExponentPushToken[d]", Platform: "unknown", Title: "T", ImageURL: image},
	})
	if err != nil {
		t.Fatal(err)
	}
	rich := func(m map[string]any) any {
		if r, ok := m["richContent"].(map[string]any); ok {
			return r["image"]
		}
		return nil
	}
	ios, android, plain, unknown := s.last[0], s.last[1], s.last[2], s.last[3]
	if rich(ios) != image || ios["mutableContent"] != true {
		t.Fatalf("ios message = %v", ios)
	}
	if rich(android) != image || android["mutableContent"] != nil {
		t.Fatalf("android message = %v", android)
	}
	// A plain notification never asks iOS to run the service extension.
	if plain["richContent"] != nil || plain["mutableContent"] != nil {
		t.Fatalf("a text notification asks for rich content: %v", plain)
	}
	if rich(unknown) != image {
		t.Fatalf("unknown platform message = %v", unknown)
	}
}

func TestExpoRequestFailuresAreClassified(t *testing.T) {
	cases := []struct {
		name      string
		status    int
		body      string
		retryable bool
		isolate   bool
		code      string
	}{
		{"rate limit", 429, `{"errors":[{"code":"TOO_MANY_REQUESTS","message":"x"}]}`, true, false, "TOO_MANY_REQUESTS"},
		{"outage", 503, `oops`, true, false, "http_503"},
		{"bad credentials", 401, `{"errors":[{"code":"UNAUTHORIZED","message":"x"}]}`, true, false, "UNAUTHORIZED"},
		{"forbidden", 403, ``, true, false, "http_403"},
		{"two projects", 400, `{"errors":[{"code":"PUSH_TOO_MANY_EXPERIENCE_IDS","message":"mixed ExponentPushToken[a]"}]}`, false, true, "PUSH_TOO_MANY_EXPERIENCE_IDS"},
		{"bad request", 422, ``, false, true, "http_422"},
	}
	for _, c := range cases {
		s := newExpoStub(t)
		s.status, s.body = c.status, c.body
		_, err := s.provider("").Send(context.Background(), []Outbound{{Token: "ExponentPushToken[a]", Title: "T"}})
		re, ok := err.(*RequestError)
		if !ok {
			t.Fatalf("%s: error = %v", c.name, err)
		}
		if re.Retryable != c.retryable || re.Isolate != c.isolate || re.Code != c.code || re.Status != c.status {
			t.Fatalf("%s: %+v", c.name, re)
		}
		if strings.Contains(re.Detail, "PushToken") {
			t.Fatalf("%s: token in %q", c.name, re.Detail)
		}
	}

	// A connection that cannot be made is transient too.
	s := newExpoStub(t)
	e := s.provider("")
	s.server.Close()
	_, err := e.Send(context.Background(), []Outbound{{Token: "ExponentPushToken[a]", Title: "T"}})
	if re, ok := err.(*RequestError); !ok || !re.Retryable || re.Code != "network" {
		t.Fatalf("network error = %v", err)
	}
}

func TestExpoAnUnreadableAnswerDoesNotCauseDuplicates(t *testing.T) {
	s := newExpoStub(t)
	s.body = `{"data":[{"status":"ok","id":"only-one"}]}`
	res, err := s.provider("").Send(context.Background(), []Outbound{
		{Token: "ExponentPushToken[a]", Title: "T"}, {Token: "ExponentPushToken[b]", Title: "T"},
	})
	if err != nil || len(res) != 2 {
		t.Fatalf("%v, %v", res, err)
	}
	for _, r := range res {
		if r.Kind != ResultAccepted || r.TicketID != "" {
			t.Fatalf("result = %+v: the messages were most likely taken, they must not be sent again", r)
		}
	}
}

func TestExpoReceipts(t *testing.T) {
	s := newExpoStub(t)
	s.body = `{"data":{
		"a":{"status":"ok"},
		"b":{"status":"error","message":"gone ExponentPushToken[x]","details":{"error":"DeviceNotRegistered"}},
		"c":{"status":"error","details":{"error":"MessageRateExceeded"}},
		"d":{"status":"error","details":{"error":"InvalidCredentials"}},
		"e":{"status":"error"}}}`
	got, err := s.provider("").Receipts(context.Background(), []string{"a", "b", "c", "d", "e", "missing"})
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]ReceiptKind{"a": ReceiptHandedOff, "b": ReceiptInvalidDevice, "c": ReceiptRetry, "d": ReceiptFailed, "e": ReceiptFailed}
	for id, kind := range want {
		if got[id].Kind != kind {
			t.Fatalf("receipt %s = %+v, want kind %v", id, got[id], kind)
		}
	}
	if _, ok := got["missing"]; ok {
		t.Fatal("a receipt that is not available must be absent")
	}
	if got["d"].Code != "InvalidCredentials" || got["e"].Code != "unknown" || strings.Contains(got["b"].Detail, "PushToken") {
		t.Fatalf("receipts = %+v", got)
	}
	if ids := s.last; ids != nil {
		t.Fatalf("the receipt request is an object, not a list: %v", ids)
	}

	s.status, s.body = 500, ``
	if _, err := s.provider("").Receipts(context.Background(), []string{"a"}); err == nil {
		t.Fatal("an outage was not reported")
	}
}

func TestExpoDefaultsAndLimits(t *testing.T) {
	e := NewExpoProvider(ExpoConfig{})
	if e.cfg.SendURL != "https://api.expo.dev/v2/push/send" || e.cfg.ReceiptsURL != "https://exp.host/--/api/v2/push/getReceipts" {
		t.Fatalf("urls = %q %q", e.cfg.SendURL, e.cfg.ReceiptsURL)
	}
	if e.MaxBatch() != 100 || e.Name() != "expo" {
		t.Fatal("Expo accepts 100 messages per request")
	}
	t.Setenv("EXPO_ACCESS_TOKEN", " abc ")
	if ExpoConfigFromEnv().AccessToken != "abc" {
		t.Fatal("the access token must come from the environment")
	}
	long := redact(strings.Repeat("é", 500))
	if len(long) > maxDetailLength || !strings.HasSuffix(long, "é") {
		t.Fatalf("redact kept %d bytes", len(long))
	}
}
