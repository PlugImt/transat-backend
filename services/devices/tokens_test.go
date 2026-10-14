package devices

import (
	"testing"
)

func TestTokenFingerprintIsStableAndDoesNotRevealTheToken(t *testing.T) {
	a, b := TokenFingerprint(tok1), TokenFingerprint(tok1)
	if a != b || a == TokenFingerprint(tok2) || len(a) != 16 {
		t.Fatalf("fingerprints = %q %q", a, TokenFingerprint(tok2))
	}
	if a == tok1 || len(a) >= len(tok1) {
		t.Fatal("the fingerprint must not be the token")
	}
}

func TestInvalidateTokenDisablesOnlyTheTokenThatWasReportedDead(t *testing.T) {
	f := newFx(t)
	f.register(idA, "alice@x.test", tok1)

	// News about a token the device no longer holds changes nothing.
	if ok, err := f.svc.InvalidateToken(f.ctx, idA, TokenFingerprint(tok2)); err != nil || ok {
		t.Fatalf("InvalidateToken(other) = %v, %v", ok, err)
	}
	if f.state("alice@x.test", idA) != StateConnected {
		t.Fatal("a good token was disabled")
	}

	ok, err := f.svc.InvalidateToken(f.ctx, idA, TokenFingerprint(tok1))
	if err != nil || !ok {
		t.Fatalf("InvalidateToken = %v, %v", ok, err)
	}
	d, _ := f.device("alice@x.test", idA)
	if d.State != StateDisabled || d.DisabledReason != "invalid_token" || d.HasPushToken {
		t.Fatalf("device = %+v", d)
	}
	if got := f.legacy("alice@x.test"); len(got) != 0 {
		t.Fatalf("legacy tokens = %v", got)
	}
	if ok, _ := f.svc.InvalidateToken(f.ctx, idA, TokenFingerprint(tok1)); ok {
		t.Fatal("a second call must be a no-op")
	}

	// The next registration of the app, with a new token, brings the device back.
	f.register(idA, "alice@x.test", tok3)
	if f.state("alice@x.test", idA) != StateConnected {
		t.Fatal("the device did not come back after registering again")
	}
}

func TestLegacyDevicesAreAsRecentAsTheirAccount(t *testing.T) {
	f := newFx(t)
	if err := RegisterLegacyToken(f.ctx, f.db, "alice@x.test", tok1); err != nil {
		t.Fatal(err)
	}
	list := func() Device {
		devices, err := f.svc.List(f.ctx, "alice@x.test")
		if err != nil || len(devices) != 1 {
			t.Fatalf("devices = %v, %v", devices, err)
		}
		return devices[0]
	}
	if d := list(); d.State != StateConnected || d.Platform != "unknown" {
		t.Fatalf("a legacy token registered just now: %+v", d)
	}

	f.exec(`UPDATE devices SET last_seen_at = now() - interval '90 days'`)
	if d := list(); d.State != StateExpired {
		t.Fatalf("no sign of life at all: %s", d.State)
	}
	f.exec(`UPDATE newf SET last_activity = now() - interval '3 days' WHERE email = 'alice@x.test'`)
	if d := list(); d.State != StateActive {
		t.Fatalf("an account active 3 days ago: %s", d.State)
	}
	f.exec(`UPDATE newf SET password_updated_date = now() WHERE email = 'alice@x.test'`)
	if d := list(); d.State != StateRevoked {
		t.Fatalf("after a password change: %s", d.State)
	}
}
