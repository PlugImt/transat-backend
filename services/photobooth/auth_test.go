package photobooth

import (
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"
)

const testSecret = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFG" // base64url, 32 bytes once decoded

func testVerifier(t *testing.T, now time.Time) *Verifier {
	t.Helper()
	keys, err := ParseDeviceKeys("booth-1:" + testSecret)
	if err != nil {
		t.Fatal(err)
	}
	v := NewVerifier(keys)
	v.now = func() time.Time { return now }
	return v
}

func signed(secret []byte, method, path string, at time.Time, nonce string, body []byte) (string, string) {
	ts := strconv.FormatInt(at.Unix(), 10)
	return ts, Sign(secret, Canonical(method, path, ts, nonce, body))
}

func TestVerifierAcceptsOnceAndRejectsTampering(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	v := testVerifier(t, now)
	secret := v.keys["booth-1"]
	body := []byte(`{"photo_count":3}`)
	nonce := "nonce-aaaaaaaaaaaaaaaa"
	ts, sig := signed(secret, "POST", "/photobooth/booth/sessions/s1/finalize", now, nonce, body)

	if _, err := v.Verify("POST", "/photobooth/booth/sessions/s1/finalize", "booth-1", ts, nonce, sig, []byte(`{"photo_count":4}`)); err == nil {
		t.Fatal("altered body accepted")
	}
	if _, err := v.Verify("POST", "/photobooth/booth/sessions/s2/finalize", "booth-1", ts, nonce, sig, body); err == nil {
		t.Fatal("altered path accepted")
	}
	if _, err := v.Verify("DELETE", "/photobooth/booth/sessions/s1/finalize", "booth-1", ts, nonce, sig, body); err == nil {
		t.Fatal("altered method accepted")
	}
	if _, err := v.Verify("POST", "/photobooth/booth/sessions/s1/finalize", "other", ts, nonce, sig, body); err == nil {
		t.Fatal("unknown key accepted")
	}
	dev, err := v.Verify("POST", "/photobooth/booth/sessions/s1/finalize", "booth-1", ts, nonce, sig, body)
	if err != nil || dev != "booth-1" {
		t.Fatalf("valid request rejected: %v", err)
	}
	if _, err := v.Verify("POST", "/photobooth/booth/sessions/s1/finalize", "booth-1", ts, nonce, sig, body); err == nil {
		t.Fatal("replayed request accepted")
	}
}

func TestVerifierRejectsStaleRequests(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	v := testVerifier(t, now)
	secret := v.keys["booth-1"]
	for _, at := range []time.Time{now.Add(-6 * time.Minute), now.Add(6 * time.Minute)} {
		ts, sig := signed(secret, "GET", "/photobooth/booth/ping", at, "nonce-bbbbbbbbbbbbbbbb", nil)
		if _, err := v.Verify("GET", "/photobooth/booth/ping", "booth-1", ts, "nonce-bbbbbbbbbbbbbbbb", sig, nil); !errors.Is(err, ErrClockSkew) {
			t.Fatalf("request at %v: got %v, want clock skew", at, err)
		}
	}
}

func TestParseDeviceKeys(t *testing.T) {
	if _, err := ParseDeviceKeys("booth-1:short"); err == nil || strings.Contains(err.Error(), "short") == false {
		t.Fatalf("short secret accepted: %v", err)
	}
	if _, err := ParseDeviceKeys("booth 1:" + testSecret); err == nil {
		t.Fatal("bad id accepted")
	}
	if _, err := ParseDeviceKeys("a:" + testSecret + ",a:" + testSecret); err == nil {
		t.Fatal("duplicate id accepted")
	}
	keys, err := ParseDeviceKeys(" a:" + testSecret + " , b:" + testSecret + "=")
	if err != nil || len(keys) != 2 || len(keys["a"]) != 32 {
		t.Fatalf("keys=%v err=%v", keys, err)
	}
}
