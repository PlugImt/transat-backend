package models

// MachineData represents the structure of washing machine data from the external API
type MachineData struct {
	MachineID        string `json:"machine_id"`
	NomType          string `json:"nom_type"`
	SelecteurMachine string `json:"selecteur_machine"`
	Status           int    `json:"status"`
	TimeBeforeOff    int    `json:"time_before_off"`
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
	Success bool                `json:"success"`
	Data    FormattedMachineData `json:"data"`
	Message string              `json:"message,omitempty"`
} 