package services

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services/netex"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

type NaolibService struct {
	db *sql.DB
}

func NewNaolibService(db *sql.DB) *NaolibService {
	return &NaolibService{
		db: db,
	}
}

func (s *NaolibService) GetDepartures(stopPlaceId string) (map[string]models.Departures, error) {

	rows, err := s.db.Query("SELECT id FROM NETEX_Quay WHERE site_ref_stopplace_id = $1", stopPlaceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quays []string
	for rows.Next() {
		var quay string
		err = rows.Scan(&quay)
		if err != nil {
			return nil, err
		}
		quays = append(quays, quay)
	}

	siriResponse, err := netex.CallStopMonitoringRequest(quays)
	if err != nil {
		return nil, err
	}

	departures := siriResponse.ServiceDelivery.StopMonitoringDelivery.MonitoredStopVisits

	departuresMap := make(map[string]models.Departures)

	for _, departure := range departures {
		lineRef := departure.MonitoredVehicleJourney.LineRef

		lineDepartures, ok := departuresMap[lineRef]
		if !ok {
			lineDepartures = models.Departures{
				DepartureDirectionAller: models.DepartureDirection{
					Direction:  "",
					Departures: []models.Departure{},
				},
				DepartureDirectionRetour: models.DepartureDirection{
					Direction:  "",
					Departures: []models.Departure{},
				},
			}
		}

		if departure.MonitoredVehicleJourney.DirectionName == "A" {
			// si la destination n'est pas "" (valeur par défaut), ça signifie qu'on l'a déjà changé. on doit aller vérifier
			// si c'est la même destination, sinon on rajoute la destination avec "/ <destination>"
			if lineDepartures.DepartureDirectionAller.Direction == "" {
				lineDepartures.DepartureDirectionAller.Direction = departure.MonitoredVehicleJourney.DestinationName
			} else {
				if !strings.Contains(lineDepartures.DepartureDirectionAller.Direction, departure.MonitoredVehicleJourney.DestinationName) {
					lineDepartures.DepartureDirectionAller.Direction += " / " + departure.MonitoredVehicleJourney.DestinationName
				}
			}

			departure := models.Departure{
				LineRef:         lineRef,
				Direction:       lineDepartures.DepartureDirectionAller.Direction,
				DestinationName: departure.MonitoredVehicleJourney.DestinationName,
				DepartureTime:   departure.MonitoredVehicleJourney.MonitoredCall.ExpectedDepartureTime,
				ArrivalTime:     departure.MonitoredVehicleJourney.MonitoredCall.ExpectedArrivalTime,
				VehicleMode:     departure.MonitoredVehicleJourney.VehicleMode,
			}

			lineDepartures.DepartureDirectionAller.Departures = append(lineDepartures.DepartureDirectionAller.Departures, departure)
		} else {
			if lineDepartures.DepartureDirectionRetour.Direction == "" {
				lineDepartures.DepartureDirectionRetour.Direction = departure.MonitoredVehicleJourney.DestinationName
			} else {
				if !strings.Contains(lineDepartures.DepartureDirectionRetour.Direction, departure.MonitoredVehicleJourney.DestinationName) {
					lineDepartures.DepartureDirectionRetour.Direction += " / " + departure.MonitoredVehicleJourney.DestinationName
				}
			}

			departure := models.Departure{
				LineRef:         lineRef,
				Direction:       lineDepartures.DepartureDirectionRetour.Direction,
				DestinationName: departure.MonitoredVehicleJourney.DestinationName,
				DepartureTime:   departure.MonitoredVehicleJourney.MonitoredCall.ExpectedDepartureTime,
				ArrivalTime:     departure.MonitoredVehicleJourney.MonitoredCall.ExpectedArrivalTime,
				VehicleMode:     departure.MonitoredVehicleJourney.VehicleMode,
			}

			lineDepartures.DepartureDirectionRetour.Departures = append(lineDepartures.DepartureDirectionRetour.Departures, departure)
		}

		departuresMap[lineRef] = lineDepartures
	}

	return departuresMap, nil
}
