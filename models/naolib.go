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
	Line            Line      `json:"line"`
	LineRef         string    `json:"lineRef"`
	Direction       string    `json:"direction"`
	DestinationName string    `json:"destinationName"`
	DepartureTime   time.Time `json:"departureTime"`
	ArrivalTime     time.Time `json:"arrivalTime"`
	VehicleMode     string    `json:"vehicleMode"`
}

type Line struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	TransportMode    string `json:"transportMode"`
	Number           string `json:"number"`
	BackgroundColour string `json:"backgroundColour"`
	ForegroundColour string `json:"foregroundColour"`
}
