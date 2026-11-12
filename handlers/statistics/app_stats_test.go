package statistics

import (
	"context"
	"testing"

	"github.com/plugimt/transat-backend/internal/database/dbtest"
)

func TestDeviceStatsCountRecentInstallationsByPlatformAndVersion(t *testing.T) {
	db := dbtest.Open(t)
	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatal(err)
		}
	}
	for _, u := range []string{"a@imt.fr", "b@imt.fr", "c@imt.fr"} {
		exec(`INSERT INTO newf (email, password, first_name, last_name) VALUES ($1, 'x', 'Jo', 'Doe')`, u)
	}
	device := func(id, email, platform, version, status string, daysAgo int) {
		var reason any
		if status == "disabled" {
			reason = "signed_out"
		}
		exec(`INSERT INTO devices (installation_id, user_email, platform, app_version, status, disabled_reason, last_seen_at)
			VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6, now() - make_interval(days => $7))`, id, email, platform, version, status, reason, daysAgo)
	}
	device("00000000-0000-0000-0000-000000000001", "a@imt.fr", "ios", "3.2.0", "active", 1)
	device("00000000-0000-0000-0000-000000000002", "a@imt.fr", "android", "3.2.0", "active", 2)
	device("00000000-0000-0000-0000-000000000003", "b@imt.fr", "android", "3.1.0", "active", 3)
	device("00000000-0000-0000-0000-000000000004", "c@imt.fr", "android", "", "active", 4)
	device("00000000-0000-0000-0000-000000000005", "c@imt.fr", "ios", "2.0.0", "active", 90)  // not seen for 3 months
	device("00000000-0000-0000-0000-000000000006", "b@imt.fr", "ios", "3.2.0", "disabled", 1) // signed out

	st, err := deviceStats(context.Background(), db)
	if err != nil {
		t.Fatal(err)
	}
	if st.Devices != 4 || st.Users != 3 {
		t.Fatalf("totals: %+v", st)
	}
	if len(st.Platforms) != 2 || st.Platforms[0] != (PlatformCount{"android", 3, 3}) || st.Platforms[1] != (PlatformCount{"ios", 1, 1}) {
		t.Fatalf("platforms: %+v", st.Platforms)
	}
	got := map[string]int{}
	for _, v := range st.Versions {
		got[v.Platform+" "+v.Version] = v.Devices
	}
	want := map[string]int{"ios 3.2.0": 1, "android 3.2.0": 1, "android 3.1.0": 1, "android ": 1}
	if len(got) != len(want) {
		t.Fatalf("versions: %+v", st.Versions)
	}
	for k, n := range want {
		if got[k] != n {
			t.Fatalf("versions: %+v", st.Versions)
		}
	}
}
