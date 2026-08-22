package models

type GTFSLineConfig struct {
	Name           string   `json:"name"`
	RouteShortName string   `json:"routeShortName"`
	DirectionID    string   `json:"directionId"`
	StopIDs        []string `json:"stopIds"`
}

type BusDeparturesResponse struct {
	LastRefresh string              `json:"lastRefresh"`
	Lines       []BusLineDepartures `json:"lines"`
}

type BusLineDepartures struct {
	Name       string         `json:"name"`
	Departures []BusDeparture `json:"departures"`
}

type BusDeparture struct {
	Time string `json:"time"`
}
