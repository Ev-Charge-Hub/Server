package usecase

import (
	"Ev-Charge-Hub/Server/internal/dto/request"
	"Ev-Charge-Hub/Server/internal/dto/response"
	"Ev-Charge-Hub/Server/internal/repository"
	"Ev-Charge-Hub/Server/internal/repository/models"
	"context"
	"fmt"
)

type EVStationUsecase interface {
	FilterStations(ctx context.Context, filter request.StationFilterRequest) ([]response.EVStationResponse, error)
	ShowAllStations(ctx context.Context) ([]response.EVStationResponse, error)
	GetStationByID(ctx context.Context, id string) (*response.EVStationResponse, error)
}

type evStationUsecase struct {
	stationRepo repository.EVStationRepository
}

func NewEVStationUsecase(repo repository.EVStationRepository) EVStationUsecase {
	return &evStationUsecase{stationRepo: repo}
}

func (u *evStationUsecase) FilterStations(ctx context.Context, filter request.StationFilterRequest) ([]response.EVStationResponse, error) {
	var isOpen *bool

	// Convert status string to boolean
    if filter.Status != "" {
        switch filter.Status {
        case "open":
            isOpen = new(bool)
            *isOpen = true
        case "closed":
            isOpen = new(bool)
            *isOpen = false
        default:
            return nil, fmt.Errorf("invalid status value: %s", filter.Status)
        }
    }
	
	stations, err := u.stationRepo.FindStations(ctx, filter.Company, filter.Type, filter.Search, filter.PlugName, isOpen)
	if err != nil {
		return nil, err
	}

	var stationResponses []response.EVStationResponse
	for _, station := range stations {
		stationResponses = append(stationResponses, mapStationDBToResponse(station))
	}
	return stationResponses, nil
}

func (u *evStationUsecase) ShowAllStations(ctx context.Context) ([]response.EVStationResponse, error) {
	stations, err := u.stationRepo.FindAllStations(ctx)
	if err != nil {
		return nil, err
	}

	var stationResponses []response.EVStationResponse
	for _, station := range stations {
		stationResponses = append(stationResponses, mapStationDBToResponse(station))
	}
	return stationResponses, nil
}

func (u *evStationUsecase) GetStationByID(ctx context.Context, id string) (*response.EVStationResponse, error) {
	station, err := u.stationRepo.FindStationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := mapStationDBToResponse(*station)
	return &response, nil
}

func mapStationDBToResponse(station models.EVStationDB) response.EVStationResponse {
	var connectors []response.Connector
	for _, c := range station.Connectors {
		connectors = append(connectors, response.Connector{
			ConnectorID:  c.ConnectorID,
			Type:         c.Type,
			PlugName:     c.PlugName,
			PricePerUnit: c.PricePerUnit,
			PowerOutput:  c.PowerOutput,
			IsAvailable:  c.IsAvailable,
		})
	}

	return response.EVStationResponse{
		ID:        station.ID.Hex(),
		StationID: station.StationID,
		Name:      station.Name,
		Latitude:  station.Latitude,
		Longitude: station.Longitude,
		Company:   station.Company,
		Status: response.StationStatus{
			OpenHours:  station.Status.OpenHours,
			CloseHours: station.Status.CloseHours,
			IsOpen:     station.Status.IsOpen,
		},
		Connectors: connectors,
	}
}
