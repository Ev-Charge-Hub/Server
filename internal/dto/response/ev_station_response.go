package response

import (
	"Ev-Charge-Hub/Server/internal/constants"
)

type EVStationResponse struct {
	ID         string                `json:"id"`
	StationID  string                `json:"station_id"`
	Name       string                `json:"name"`
	Latitude   float64               `json:"latitude"`
	Longitude  float64               `json:"longitude"`
	Company    string                `json:"company"`
	Status     StationStatusResponse `json:"status"`
	Connectors []ConnectorResponse   `json:"connectors"`
}

type StationStatusResponse struct {
	OpenHours  string `json:"open_hours"`
	CloseHours string `json:"close_hours"`
	IsOpen     bool   `json:"is_open"`
}

type ConnectorResponse struct {
	ConnectorID  string                  `json:"connector_id"`
	Type         constants.ConnectorType `json:"type"`
	PlugName     constants.PlugName      `json:"plug_name"`
	PricePerUnit float64                 `json:"price_per_unit"`
	PowerOutput  int                     `json:"power_output"`
	Booking      *BookingResponse        `json:"booking,omitempty"`
}

type BookingResponse struct {
	Username       string `json:"username"`
	BookingEndTime string `json:"booking_end_time"`
}
