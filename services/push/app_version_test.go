package push

import (
	"sort"
	"testing"

	"github.com/plugimt/transat-backend/models"
)

func TestTheAppVersionAudienceReachesOnlyOutdatedDevices(t *testing.T) {
	f := newFx(t)
	f.user("old@x.test", 1)
	f.user("both@x.test", 1)
	f.user("fresh@x.test", 1)
	f.device(1, "old@x.test", "ios", true)
	f.device(2, "both@x.test", "android", true)
	f.device(3, "both@x.test", "ios", true)
	f.device(4, "fresh@x.test", "ios", true)
	f.device(5, "old@x.test", "android", true) // never reported its version
	for id, version := range map[int]string{1: "2.3.9", 2: "2.10.0", 3: "2.4.0-beta", 4: "2.10.1"} {
		f.exec(`UPDATE devices SET app_version = $2 WHERE installation_id = $1`, uuid(id), version)
	}

	resolve := func(platform string) []string {
		t.Helper()
		aud := Audience{NotificationAudience: models.NotificationAudience{Type: models.AudienceAppVersion, BelowVersion: "2.10.0", Platform: platform}}
		if err := aud.Validate(); err != nil {
			t.Fatal(err)
		}
		res, err := f.p.Resolve(f.ctx, aud, "")
		if err != nil {
			t.Fatal(err)
		}
		var ids []string
		for _, d := range res.Targets {
			ids = append(ids, d.DeviceID)
		}
		sort.Strings(ids)
		return ids
	}

	// 2.3.9 < 2.10.0 (numeric, not text, order); 2.4.0-beta < 2.10.0; a device without version is outdated;
	// 2.10.0 itself and newer are left alone, even for a user who also has an outdated device.
	got := resolve("")
	want := []string{uuid(1), uuid(3), uuid(5)}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Fatalf("all platforms: %v, want %v", got, want)
	}
	if got := resolve("android"); len(got) != 1 || got[0] != uuid(5) {
		t.Fatalf("android only: %v", got)
	}

	if err := (models.NotificationAudience{Type: models.AudienceAppVersion, BelowVersion: "two"}).Validate(); err == nil {
		t.Fatal("a version that is not numbers was accepted")
	}
}
