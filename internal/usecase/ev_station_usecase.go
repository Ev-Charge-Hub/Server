package usecase

import (
	"Ev-Charge-Hub/Server/internal/dto/request"
	"Ev-Charge-Hub/Server/internal/dto/response"
	"Ev-Charge-Hub/Server/internal/repository"
	"Ev-Charge-Hub/Server/internal/repository/models"
	"context"
	"fmt"
	"time"
)

type EVStationUsecase interface {
	FilterStations(ctx context.Context, filter request.StationFilterRequest) ([]response.EVStationResponse, error)
	ShowAllStations(ctx context.Context) ([]response.EVStationResponse, error)
	GetStationByID(ctx context.Context, id string) (*response.EVStationResponse, error)
	CreateStation(ctx context.Context, station models.EVStationDB) error
	EditStation(ctx context.Context, id string, station models.EVStationDB) error
	RemoveStation(ctx context.Context, id string) error
	SetBooking(ctx context.Context, booking request.SetBookingRequest) error
}

// Create Class
type evStationUsecase struct {
	stationRepo repository.EVStationRepository
}

// Init class && imprement EVStationUsecase interface
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

func (u *evStationUsecase) CreateStation(ctx context.Context, station models.EVStationDB) error {
	return u.stationRepo.CreateStation(ctx, station)
}

func (u *evStationUsecase) EditStation(ctx context.Context, id string, station models.EVStationDB) error {
	return u.stationRepo.EditStation(ctx, id, station)
}

func (u *evStationUsecase) RemoveStation(ctx context.Context, id string) error {
	return u.stationRepo.RemoveStation(ctx, id)
}

func (u *evStationUsecase) SetBooking(ctx context.Context, booking request.SetBookingRequest) error {
	// Validate Date Format
	_, err := time.Parse("2006-01-02T15:04:05", booking.BookingEndTime)
	if err != nil {
		return fmt.Errorf("invalid booking_end_time format")
	}

	bookingDB := models.BookingDB{
		Username:       booking.Username,
		BookingEndTime: booking.BookingEndTime,
	}

	return u.stationRepo.SetBooking(ctx, booking.ConnectorId, bookingDB)
}

func mapStationDBToResponse(station models.EVStationDB) response.EVStationResponse {
	var connectors []response.ConnectorResponse
	for _, c := range station.Connectors {
		var booking *response.BookingResponse = nil

		// ตรวจสอบว่ามี Booking อยู่หรือไม่
		if c.Booking != nil {
			booking = &response.BookingResponse{
				Username:       c.Booking.Username,
				BookingEndTime: c.Booking.BookingEndTime,
			}
		}

		connectors = append(connectors, response.ConnectorResponse{
			ConnectorID:  c.ConnectorID,
			Type:         c.Type,
			PlugName:     c.PlugName,
			PricePerUnit: c.PricePerUnit,
			PowerOutput:  c.PowerOutput,
			Booking:      booking,
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
