package washingmachine

import (
	"io"
	"net/http"
	"testing"
)

// The shape the operator's endpoint returns today: ids and selectors are numbers.
const numericBody = `{"machine_info_status":{"centrale_nom":"1021 EDM DE NANTES","machine_list":[
	{"machine_id":12629,"selecteur_machine":11,"nom":"","nom_type":"LAVE LINGE ","status":12,"temps":0,"time_before_off":157,"time_before_off_virtu":0},
	{"machine_id":16315,"selecteur_machine":12,"nom":"","nom_type":"LAVE LINGE ","status":1,"temps":0,"time_before_off":0,"time_before_off_virtu":0},
	{"machine_id":16316,"selecteur_machine":13,"nom":"","nom_type":"LAVE LINGE ","status":6,"temps":125,"time_before_off":0,"time_before_off_virtu":0},
	{"machine_id":16318,"selecteur_machine":15,"nom":"","nom_type":"SECHE LINGE","status":1,"temps":0,"time_before_off":0,"time_before_off_virtu":0}]}}`

func machinesOf(t *testing.T, body string) (washers, dryers []map[string]any) {
	t.Helper()
	app, _ := newHandler(t, func(w http.ResponseWriter, _ *http.Request) { _, _ = io.WriteString(w, body) })
	status, response, _ := get(t, app, "")
	if status != http.StatusOK {
		t.Fatalf("status %d, body %v", status, response)
	}
	data := response["data"].(map[string]any)
	for _, m := range data["washing_machine"].([]any) {
		washers = append(washers, m.(map[string]any))
	}
	for _, m := range data["dryer"].([]any) {
		dryers = append(dryers, m.(map[string]any))
	}
	return washers, dryers
}

func TestNumericIdsAndSelectorsAreRead(t *testing.T) {
	washers, dryers := machinesOf(t, numericBody)

	if len(washers) != 3 || len(dryers) != 1 {
		t.Fatalf("%d washers and %d dryers, want 3 and 1", len(washers), len(dryers))
	}
	running := washers[0]
	if running["number"] != float64(11) || running["available"] != false || running["time_left"] != float64(157) {
		t.Fatalf("running washer: %v", running)
	}
	if washers[1]["available"] != true || dryers[0]["number"] != float64(15) {
		t.Fatalf("free washer %v, dryer %v", washers[1], dryers[0])
	}
	// Status 6 counts down in "temps", as on the operator's own page.
	if washers[2]["time_left"] != float64(125) {
		t.Fatalf("status 6 washer: %v", washers[2])
	}
}

func TestStringIdsAndSelectorsAreStillRead(t *testing.T) {
	washers, dryers := machinesOf(t, validBody)

	if len(washers) != 1 || len(dryers) != 1 || washers[0]["number"] != float64(11) || dryers[0]["time_left"] != float64(300) {
		t.Fatalf("washers %v, dryers %v", washers, dryers)
	}
}
