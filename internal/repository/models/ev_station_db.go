package models

import (
	"Ev-Charge-Hub/Server/internal/constants"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EVStationDB represents the core domain model for EV Stations
type EVStationDB struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	Latitude   float64            `bson:"latitude"`
	Longitude  float64            `bson:"longitude"`
	Company    string             `bson:"company"`
	Status     StationStatusDB    `bson:"status"`
	Connectors []ConnectorDB      `bson:"connectors"`
}

// StationStatusDB represents the status details of an EV Station
type StationStatusDB struct {
	OpenHours  string `bson:"open_hours"`
	CloseHours string `bson:"close_hours"`
	IsOpen     bool   `bson:"is_open"`
}

// ConnectorDB represents a charging connector within an EV Station
type ConnectorDB struct {
	ConnectorID  string                  `bson:"connector_id"`
	Type         constants.ConnectorType `bson:"type"`
	PlugName     constants.PlugName      `bson:"plug_name"`
	PricePerUnit float64                 `bson:"price_per_unit"`
	PowerOutput  int                     `bson:"power_output"`
	Booking      *BookingDB              `bson:"booking,omitempty"`
}

// BookingDB represents booking details for each connector
type BookingDB struct {
	Username       string `bson:"username"`
	BookingEndTime string `bson:"booking_end_time"`
}
