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
	FilterStations(ctx context.Context, request request.StationFilterRequest) ([]response.EVStationResponse, error)
	ShowAllStations(ctx context.Context) ([]response.EVStationResponse, error)
	GetStationByID(ctx context.Context, id string) (*response.EVStationResponse, error)
	CreateStation(ctx context.Context, station models.EVStationDB) error
	EditStation(ctx context.Context, id string, station models.EVStationDB) error
	RemoveStation(ctx context.Context, request request.RemoveStationRequest) error
	SetBooking(ctx context.Context, request request.SetBookingRequest) error
	GetBookingByUserName(ctx context.Context, request request.GetBookingRequest) (*response.BookingResponse, error)
	GetBookingsByUserName(ctx context.Context, request request.GetBookingsRequest) ([]response.BookingResponse, error)
	GetStationByConnectorID(ctx context.Context, request request.GetStationByConnectorIDRequest) (*response.EVStationResponse, error)
	GetStationByUserName(ctx context.Context, request request.GetStationByUsernameRequest) (*response.EVStationResponse, error)
}

// Create Class
type evStationUsecase struct {
	stationRepo repository.EVStationRepository
}

// Init class && imprement EVStationUsecase interface
func NewEVStationUsecase(repo repository.EVStationRepository) EVStationUsecase {
	return &evStationUsecase{stationRepo: repo}
}

func (u *evStationUsecase) FilterStations(ctx context.Context, request request.StationFilterRequest) ([]response.EVStationResponse, error) {
	var isOpen *bool

	// Convert status string to boolean
	if request.Status != "" {
		switch request.Status {
		case "open":
			isOpen = new(bool)
			*isOpen = true
		case "closed":
			isOpen = new(bool)
			*isOpen = false
		default:
			return nil, fmt.Errorf("invalid status value: %s", request.Status)
		}
	}

	stations, err := u.stationRepo.FindStations(ctx, request.Company, request.Type, request.Search, request.PlugName, isOpen)
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

func (u *evStationUsecase) RemoveStation(ctx context.Context, request request.RemoveStationRequest) error {
	return u.stationRepo.RemoveStation(ctx, request.ID)
}

// func (u *evStationUsecase) SetBooking(ctx context.Context, booking request.SetBookingRequest) error {
// 	// Validate Date Format
// 	_, err := time.Parse("2006-01-02T15:04:05", booking.BookingEndTime)
// 	if err != nil {
// 		return fmt.Errorf("invalid booking_end_time format")
// 	}
//
// 	bookingDB := models.BookingDB{
// 		Username:       booking.Username,
// 		BookingEndTime: booking.BookingEndTime,
// 	}
//
// 	haveBooking, err := u.stationRepo.FindBookingByUserName(ctx, booking.Username)
//
// 	if err != nil {
// 		return fmt.Errorf("error finding booking: %v", err)
// 	}
//
// 	// Check if the booking end time is in the past
// 	if haveBooking != nil && haveBooking.BookingEndTime < time.Now().Format("2006-01-02T15:04:05") {
// 		return fmt.Errorf("user already has a booking")
// 	}
//
// 	return u.stationRepo.SetBooking(ctx, booking.ConnectorId, bookingDB)
// }

func (u *evStationUsecase) GetStationByConnectorID(ctx context.Context, request request.GetStationByConnectorIDRequest) (*response.EVStationResponse, error) {
	station, err := u.stationRepo.FindStationByConnectorID(ctx, request.ConnectorId)
	if err != nil {
		return nil, err
	}
	response := mapStationDBToResponse(*station)
	return &response, nil
}

func (u *evStationUsecase) SetBooking(ctx context.Context, request request.SetBookingRequest) error {
	// 📥 ผู้ใช้ส่งคำขอจอง (connector_id + username + booking_end_time)

	// 1️⃣ เช็กว่า booking_end_time > เวลาปัจจุบันไหม
	// 	❌ ถ้าไม่ → "จองย้อนหลังไม่ได้"

	// 2️⃣ เช็กว่า user คนนี้ เคยจองอื่นไว้ที่ยังไม่หมดเวลาไหม
	// 	❌ ถ้ามี → "ห้ามจองซ้ำ"

	// 3️⃣ เช็กว่า connector นี้ มีการจองอื่นอยู่ที่ยังไม่หมดเวลาไหม
	// 	❌ ถ้ามี → "มีคนจองไปแล้ว"

	// ✅ ถ้าผ่านทั้ง 3 ข้อ → ทำการจอง

	// Validate booking_end_time format
	endTime, err := time.Parse("2006-01-02T15:04:05", request.BookingEndTime)
	if err != nil {
		return fmt.Errorf("invalid booking_end_time format")
	}

	// 1️⃣ ห้ามจองย้อนหลัง === เช็กว่า booking_end_time > เวลาปัจจุบันไหม
	if !endTime.After(time.Now()) {
		return fmt.Errorf("booking_end_time must be in the future")
	}

	// เช็กว่าผู้ใช้มี booking ซ้ำอยู่หรือไม่
	bookings, err := u.stationRepo.FindBookingsByUserName(ctx, request.Username)
	if err == nil {
		for _, b := range bookings {
			expiredAt, err := time.Parse("2006-01-02T15:04:05", b.BookingEndTime)
			if err != nil {
				continue
			}
			//  2️⃣ ถ้ายังไม่หมดอายุ ห้ามจองใหม่
			if time.Now().Before(expiredAt) {
				return fmt.Errorf("user already has an active booking until %s", b.BookingEndTime)
			}

		}
	}

	// 3️⃣ เช็กว่า connector นี้ มีการจองอื่นอยู่ที่ยังไม่หมดเวลาไหม  ❌ ถ้ามี → "มีคนจองไปแล้ว
	station, err := u.stationRepo.FindStationByConnectorID(ctx, request.ConnectorId)
	if err != nil {
		return fmt.Errorf("error finding connector: %v", err)
	}

	connectorFound := false
	for _, c := range station.Connectors {
		if c.ConnectorID == request.ConnectorId {
			connectorFound = true
			if c.Booking != nil {
				expiredAt, err := time.Parse("2006-01-02T15:04:05", c.Booking.BookingEndTime)
				if err == nil && time.Now().Before(expiredAt) {
					return fmt.Errorf("connector is already booked until %s", c.Booking.BookingEndTime)
				}
			}
			break
		}
	}

	if !connectorFound {
		return fmt.Errorf("connector not found")
	}

	// ✅ Create BookingDB object
	bookingDB := models.BookingDB{
		Username:       request.Username,
		BookingEndTime: request.BookingEndTime,
	}

	// ✅ Save to repository
	return u.stationRepo.SetBooking(ctx, request.ConnectorId, bookingDB)
}

func (u *evStationUsecase) GetBookingByUserName(ctx context.Context, request request.GetBookingRequest) (*response.BookingResponse, error) {
	booking, err := u.stationRepo.FindBookingByUserName(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	// Map BookingDB to BookingResponse
	response := &response.BookingResponse{
		Username:       booking.Username,
		BookingEndTime: booking.BookingEndTime,
	}

	return response, nil
}

func (u *evStationUsecase) GetBookingsByUserName(ctx context.Context, request request.GetBookingsRequest) ([]response.BookingResponse, error) {
	bookings, err := u.stationRepo.FindBookingsByUserName(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	var result []response.BookingResponse
	for _, b := range bookings {
		result = append(result, response.BookingResponse{
			Username:       b.Username,
			BookingEndTime: b.BookingEndTime,
		})
	}

	return result, nil
}

func (u *evStationUsecase) GetStationByUserName(ctx context.Context, request request.GetStationByUsernameRequest) (*response.EVStationResponse, error) {
	station, err := u.stationRepo.FindStationByUserName(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	resp := mapStationDBToResponse(*station)
	return &resp, nil
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
		Status: response.StationStatusResponse{
			OpenHours:  station.Status.OpenHours,
			CloseHours: station.Status.CloseHours,
			IsOpen:     station.Status.IsOpen,
		},
		Connectors: connectors,
	}
}
