package models

import (
	"fmt"
	"strconv"
	"strings"
)

// FlexInt reads an integer that the external API sends either as a number or as a string: it changed from one to the
// other once already, which made every response unreadable.
type FlexInt int

func (f *FlexInt) UnmarshalJSON(data []byte) error {
	text := strings.TrimSpace(strings.Trim(string(data), `"`))
	if text == "" || text == "null" {
		*f = 0
		return nil
	}
	number, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return fmt.Errorf("not an integer: %s", data)
	}
	*f = FlexInt(number)
	return nil
}

// MachineData represents the structure of washing machine data from the external API
type MachineData struct {
	MachineID          FlexInt `json:"machine_id"`
	NomType            string  `json:"nom_type"`
	SelecteurMachine   FlexInt `json:"selecteur_machine"`
	Status             FlexInt `json:"status"`
	Temps              FlexInt `json:"temps"`
	TimeBeforeOff      FlexInt `json:"time_before_off"`
	TimeBeforeOffVirtu FlexInt `json:"time_before_off_virtu"`
}

// FormattedMachine represents the transformed machine data
type FormattedMachine struct {
	Number    int  `json:"number"`
	Available bool `json:"available"`
	TimeLeft  int  `json:"time_left"`
}

// FormattedMachineData represents the categorized machines
type FormattedMachineData struct {
	WashingMachines []FormattedMachine `json:"washing_machine"`
	Dryers          []FormattedMachine `json:"dryer"`
}

// WashingMachineResponse represents the response structure for the washing machines endpoint
type WashingMachineResponse struct {
	Success bool                 `json:"success"`
	Data    FormattedMachineData `json:"data"`
	Message string               `json:"message,omitempty"`
}
