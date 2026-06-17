package models

type BusDeparture struct {
	Name           string `json:"name"`
	NextDeparture  string `json:"nextDeparture"`
	NextDeparture2 string `json:"nextDeparture2"`
}
