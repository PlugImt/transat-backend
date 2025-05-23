package models

import "encoding/xml"

type NETEXStopsFile struct {
	XMLName             xml.Name            `xml:"PublicationDelivery"`
	PublicationDelivery PublicationDelivery `xml:"PublicationDelivery"`
}

type PublicationDelivery struct {
	XMLName              xml.Name    `xml:"PublicationDelivery"`
	PublicationTimestamp string      `xml:"PublicationTimestamp"`
	ParticipantRef       string      `xml:"ParticipantRef"`
	DataObjects          DataObjects `xml:"dataObjects"`
	Version              string      `xml:"version,attr"`
	Xmlns                string      `xml:"xmlns,attr"`
	XmlnsAcbs            string      `xml:"xmlns:acbs,attr"`
	XmlnsGML             string      `xml:"xmlns:gml,attr"`
	XmlnsIfopt           string      `xml:"xmlns:ifopt,attr"`
	XmlnsSiri            string      `xml:"xmlns:siri,attr"`
	XmlnsXSD             string      `xml:"xmlns:xsd,attr"`
	XmlnsXsi             string      `xml:"xmlns:xsi,attr"`
}

type DataObjects struct {
	GeneralFrame GeneralFrame `xml:"GeneralFrame"`
}

type GeneralFrame struct {
	ID             string         `xml:"id,attr"`
	Modification   string         `xml:"modification,attr"`
	Version        string         `xml:"version,attr"`
	TypeOfFrameRef TypeOfFrameRef `xml:"TypeOfFrameRef"`
	Members        Members        `xml:"members"`
}

type TypeOfFrameRef struct {
	Ref  string `xml:"ref,attr"`
	Text string `xml:",chardata"`
}

type Members struct {
	Quays      []Quay      `xml:"Quay"`
	StopPlaces []StopPlace `xml:"StopPlace"`
}

type Quay struct {
	ID            string           `xml:"id,attr"`
	Modification  string           `xml:"modification,attr"`
	Version       string           `xml:"version,attr"`
	KeyList       []KeyValueHolder `xml:"keyList>KeyValue"`
	Name          string           `xml:"Name"`
	Centroid      Centroid         `xml:"Centroid"`
	PostalAddress PostalAddress    `xml:"PostalAddress"`
	SiteRef       SiteRef          `xml:"SiteRef"`
	TransportMode string           `xml:"TransportMode"`
}

type StopPlace struct {
	ID                  string           `xml:"id,attr"`
	Modification        string           `xml:"modification,attr"`
	Version             string           `xml:"version,attr"`
	Name                string           `xml:"Name"`
	Centroid            Centroid         `xml:"Centroid"`
	TransportMode       string           `xml:"TransportMode"`
	OtherTransportModes string           `xml:"OtherTransportModes"`
	StopPlaceType       string           `xml:"StopPlaceType"`
	Weighting           string           `xml:"Weighting"`
	QuayRefs            []QuayRef        `xml:"quays>QuayRef"`
	KeyList             []KeyValueHolder `xml:"keyList>KeyValue"`
}

type KeyValueHolder struct {
	Key   string `xml:"Key"`
	Value string `xml:"Value"`
}

type QuayRef struct {
	Ref     string `xml:"ref,attr"`
	Version string `xml:"version,attr"`
}

type Centroid struct {
	Location Location `xml:"Location"`
}

type Location struct {
	Longitude string `xml:"Longitude"`
	Latitude  string `xml:"Latitude"`
}

type PostalAddress struct {
	ID           string `xml:"id,attr"`
	Version      string `xml:"version,attr"`
	Name         string `xml:"Name"`
	PostalRegion string `xml:"PostalRegion"`
}

type SiteRef struct {
	Ref string `xml:"ref,attr"`
}
