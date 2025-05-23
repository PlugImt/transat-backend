package models

import "time"

type DepartureDirection struct {
	Direction  string      `json:"direction"`
	Departures []Departure `json:"departures"`
}

type Departures struct {
	DepartureDirectionAller  DepartureDirection `json:"aller"`
	DepartureDirectionRetour DepartureDirection `json:"retour"`
}

type Departure struct {
	LineRef         string    `json:"lineRef"`
	Direction       string    `json:"direction"`
	DestinationName string    `json:"destinationName"`
	DepartureTime   time.Time `json:"departureTime"`
	ArrivalTime     time.Time `json:"arrivalTime"`
	VehicleMode     string    `json:"vehicleMode"`
}
