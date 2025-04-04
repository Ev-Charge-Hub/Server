package request

import (
	"Ev-Charge-Hub/Server/internal/constants"
)

type EVStationRequest struct {
	Name       string               `json:"name" binding:"required"`
	Latitude   float64              `json:"latitude" binding:"required"`
	Longitude  float64              `json:"longitude" binding:"required"`
	Company    string               `json:"company" binding:"required"`
	Status     StationStatusRequest `json:"status" binding:"required"`
	Connectors []ConnectorRequest   `json:"connectors" binding:"required,dive"`
}

type StationStatusRequest struct {
	OpenHours  string `json:"open_hours" binding:"required"`
	CloseHours string `json:"close_hours" binding:"required"`
	IsOpen     bool   `json:"is_open" binding:"required"`
}

type ConnectorRequest struct {
	Type         constants.ConnectorType `json:"type" binding:"required"`
	PlugName     constants.PlugName      `json:"plug_name" binding:"required"`
	PricePerUnit float64                 `json:"price_per_unit" binding:"required"`
	PowerOutput  int                     `json:"power_output" binding:"required"`
	Booking      *BookingRequest         `json:"booking"`
}

type BookingRequest struct {
	Username       string `json:"username" binding:"required"`
	BookingEndTime string `json:"booking_end_time" binding:"required"`
}

type GetStationByUsernameRequest struct {
	Username string `json:"username" binding:"required"`
}

type GetStationByConnectorIDRequest struct {
	ConnectorId string `json:"connector_id" binding:"required"`
}

type GetStationByIDRequest struct {
	ID string `json:"id" binding:"required"`
}

type EditStationRequest struct {
	ID         string                `json:"id" binding:"required"`
	Name       *string               `json:"name,omitempty"`
	Latitude   *float64              `json:"latitude,omitempty"`
	Longitude  *float64              `json:"longitude,omitempty"`
	Company    *string               `json:"company,omitempty"`
	Status     *StationStatusRequest `json:"status,omitempty"`
	Connectors *[]ConnectorRequest   `json:"connectors,omitempty"`
}

type RemoveStationRequest struct {
	ID string `json:"id" binding:"required"`
}
