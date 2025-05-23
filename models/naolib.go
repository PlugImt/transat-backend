package models

type DepartureDirection struct {
	Direction  string               `json:"direction"`
	Departures []MonitoredStopVisit `json:"departures"`
}

type Departures struct {
	DepartureDirectionAller  DepartureDirection `json:"aller"`
	DepartureDirectionRetour DepartureDirection `json:"retour"`
}
