package faces

import "testing"

func TestOnlyProductionScansWithoutFolders(t *testing.T) {
	t.Setenv("FACES_ML_URL", "http://face-ml.test:3003")
	t.Setenv("GOOGLE_DRIVE_CREDENTIALS_JSON", "")
	cases := []struct {
		env, scope string
		disabled   bool
	}{
		{"dev", "", true},
		{"pr-12", "", true},
		{"dev", "FOLDER-1, FOLDER-2", false},
		{"production", "", false},
		{"", "", false}, // not on Railway (local runs, tests)
	}
	for _, c := range cases {
		t.Setenv("RAILWAY_ENVIRONMENT_NAME", c.env)
		t.Setenv("FACES_SCOPE_FOLDERS", c.scope)
		s := FromEnv(t.Context(), nil)
		if s.cfg.ScanDisabled != c.disabled {
			t.Errorf("env %q scope %q: scan disabled = %v, want %v", c.env, c.scope, s.cfg.ScanDisabled, c.disabled)
		}
	}
}
