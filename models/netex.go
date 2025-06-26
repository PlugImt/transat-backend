package models

import (
	"encoding/xml"
	"os"
)

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

type NETEXCommonFile struct {
	XMLName              xml.Name          `xml:"PublicationDelivery"`
	PublicationTimestamp string            `xml:"PublicationTimestamp"`
	ParticipantRef       string            `xml:"ParticipantRef"`
	DataObjects          CommonDataObjects `xml:"dataObjects"`
}

type CommonDataObjects struct {
	GeneralFrame GeneralFrameOffer `xml:"GeneralFrame"`
}

type GeneralFrameOffer struct {
	ID             string         `xml:"id,attr"`
	Version        string         `xml:"version,attr"`
	TypeOfFrameRef TypeOfFrameRef `xml:"TypeOfFrameRef"`
	Members        struct {
		Network     Network     `xml:"Network"`
		Lines       []SIRILine  `xml:"Line"`
		Operators   []Operator  `xml:"Operator"`
		Authorities []Authority `xml:"Authority"`
	} `xml:"members"`
}

type Network struct {
	ID           string       `xml:"id,attr"`
	Version      string       `xml:"version,attr"`
	Name         string       `xml:"Name"`
	Members      []LineRef    `xml:"members>LineRef"`
	AuthorityRef AuthorityRef `xml:"AuthorityRef"`
}

type LineRef struct {
	Ref     string `xml:"ref,attr"`
	Version string `xml:"version,attr,omitempty"`
}

type AuthorityRef struct {
	Ref     string `xml:"ref,attr"`
	Version string `xml:"version,attr,omitempty"`
}

type SIRILine struct {
	ID                      string                  `xml:"id,attr"`
	Version                 string                  `xml:"version,attr"`
	KeyList                 []KeyValue              `xml:"keyList>KeyValue"`
	Name                    string                  `xml:"Name"`
	ShortName               string                  `xml:"ShortName"`
	TransportMode           string                  `xml:"TransportMode"`
	PublicCode              string                  `xml:"PublicCode"`
	PrivateCode             string                  `xml:"PrivateCode"`
	OperatorRef             LineRef                 `xml:"OperatorRef"`
	Presentation            Presentation            `xml:"Presentation"`
	AccessibilityAssessment AccessibilityAssessment `xml:"AccessibilityAssessment"`
}

type KeyValue struct {
	Key   string `xml:"Key"`
	Value string `xml:"Value"`
}

type Presentation struct {
	Colour     string `xml:"Colour"`
	TextColour string `xml:"TextColour"`
}

type AccessibilityAssessment struct {
	ID                     string      `xml:"id,attr"`
	Version                string      `xml:"version,attr"`
	MobilityImpairedAccess string      `xml:"MobilityImpairedAccess"`
	Limitations            Limitations `xml:"limitations"`
}

type Limitations struct {
	AccessibilityLimitation AccessibilityLimitation `xml:"AccessibilityLimitation"`
}

type AccessibilityLimitation struct {
	ID               string `xml:"id,attr"`
	Version          string `xml:"version,attr"`
	WheelchairAccess string `xml:"WheelchairAccess"`
}

type Operator struct {
	ID               string         `xml:"id,attr"`
	Version          string         `xml:"version,attr"`
	CompanyNumber    string         `xml:"CompanyNumber"`
	Name             string         `xml:"Name"`
	ContactDetails   ContactDetails `xml:"ContactDetails"`
	OrganisationType string         `xml:"OrganisationType"`
}

type Authority struct {
	ID               string         `xml:"id,attr"`
	Version          string         `xml:"version,attr"`
	CompanyNumber    string         `xml:"CompanyNumber"`
	Name             string         `xml:"Name"`
	ContactDetails   ContactDetails `xml:"ContactDetails"`
	OrganisationType string         `xml:"OrganisationType"`
}

type ContactDetails struct {
	Phone string `xml:"Phone"`
	Url   string `xml:"Url"`
}

func DecodeCommonNetexData(file string) (*NETEXCommonFile, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var netexData NETEXCommonFile
	err = xml.Unmarshal(data, &netexData)
	if err != nil {
		return nil, err
	}

	return &netexData, nil
}
