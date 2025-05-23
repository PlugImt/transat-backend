package models

import (
	"encoding/xml"
	"time"
)

type SIRI struct {
	XMLName         xml.Name        `xml:"Siri"`
	ServiceDelivery ServiceDelivery `xml:"ServiceDelivery"`
}

type ServiceDelivery struct {
	ResponseTimestamp      time.Time              `xml:"ResponseTimestamp"`
	ProducerRef            string                 `xml:"ProducerRef"`
	RequestMessageRef      string                 `xml:"RequestMessageRef"`
	Status                 bool                   `xml:"Status"`
	MoreData               bool                   `xml:"MoreData"`
	StopMonitoringDelivery StopMonitoringDelivery `xml:"StopMonitoringDelivery"`
}

type StopMonitoringDelivery struct {
	ResponseTimestamp   time.Time            `xml:"ResponseTimestamp"`
	MonitoredStopVisits []MonitoredStopVisit `xml:"MonitoredStopVisit"`
}

type MonitoredStopVisit struct {
	RecordedAtTime          time.Time               `xml:"RecordedAtTime"`
	ItemIdentifier          string                  `xml:"ItemIdentifier"`
	MonitoringRef           string                  `xml:"MonitoringRef"`
	MonitoredVehicleJourney MonitoredVehicleJourney `xml:"MonitoredVehicleJourney"`
}

type MonitoredVehicleJourney struct {
	LineRef                 string                  `xml:"LineRef"`
	FramedVehicleJourneyRef FramedVehicleJourneyRef `xml:"FramedVehicleJourneyRef"`
	VehicleMode             string                  `xml:"VehicleMode"`
	PublishedLineName       string                  `xml:"PublishedLineName"`
	DirectionName           string                  `xml:"DirectionName"`
	DestinationRef          string                  `xml:"DestinationRef"`
	DestinationName         string                  `xml:"DestinationName"`
	FirstOrLastJourney      string                  `xml:"FirstOrLastJourney"`
	Monitored               bool                    `xml:"Monitored"`
	MonitoredCall           MonitoredCall           `xml:"MonitoredCall"`
}

type FramedVehicleJourneyRef struct {
	DataFrameRef           string `xml:"DataFrameRef"`
	DatedVehicleJourneyRef string `xml:"DatedVehicleJourneyRef"`
}

type MonitoredCall struct {
	StopPointRef          string    `xml:"StopPointRef"`
	Order                 int       `xml:"Order"`
	AimedArrivalTime      time.Time `xml:"AimedArrivalTime"`
	ExpectedArrivalTime   time.Time `xml:"ExpectedArrivalTime"`
	ArrivalStatus         string    `xml:"ArrivalStatus"`
	AimedDepartureTime    time.Time `xml:"AimedDepartureTime"`
	ExpectedDepartureTime time.Time `xml:"ExpectedDepartureTime"`
}
