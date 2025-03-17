package response

import "Ev-Charge-Hub/Server/internal/constants"

type EVStationResponse struct {
	ID         string        `json:"id"`
	StationID  string        `json:"station_id"`
	Name       string        `json:"name"`
	Latitude   float64       `json:"latitude"`
	Longitude  float64       `json:"longitude"`
	Company    string        `json:"company"`
	Status     StationStatus `json:"status"`
	Connectors []Connector   `json:"connectors"`
}

type StationStatus struct {
	OpenHours  string `json:"open_hours"`
	CloseHours string `json:"close_hours"`
	IsOpen     bool   `json:"is_open"`
}

type Connector struct {
	ConnectorID  string                  `json:"connector_id"`
	Type          constants.ConnectorType `json:"type"`
	PlugName      constants.PlugName       `json:"plug_name"`
	PricePerUnit  float64                 `json:"price_per_unit"`
	PowerOutput   int                     `json:"power_output"`
	Booking       *Booking                `json:"booking"` // 👈 ใช้ Pointer
}

type Booking struct {
	Username       string `json:"username"`
	BookingEndTime string `json:"booking_end_time"`
}
