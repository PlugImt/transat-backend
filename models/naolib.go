package models

// Departure représente un départ de bus
type Departure struct {
	Sens          int    `json:"sens"`
	Terminus      string `json:"terminus"`
	InfoTrafic    bool   `json:"infotrafic"`
	Temps         string `json:"temps"`
	DernierDepart string `json:"dernierDepart"`
	TempsReel     string `json:"tempsReel"`
	Ligne         Ligne  `json:"ligne"`
	Arret         Arret  `json:"arret"`
}

// Ligne représente les informations sur la ligne de bus
type Ligne struct {
	NumLigne  string `json:"numLigne"`
	TypeLigne int    `json:"typeLigne"`
}

// Arret représente les informations sur l'arrêt de bus
type Arret struct {
	CodeArret string `json:"codeArret"`
}
