package models

import (
	"Ev-Charge-Hub/Server/internal/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EVStation struct {
	ID         primitive.ObjectID
	Name       string
	Latitude   float64
	Longitude  float64
	Company    string
	Status     StationStatus
	Connectors []Connector
}

type StationStatus struct {
	OpenHours  string
	CloseHours string
	IsOpen     bool
}

type Connector struct {
	ConnectorID  string
	Type         constants.ConnectorType
	PlugName     constants.PlugName
	PricePerUnit float64
	PowerOutput  int
	Booking      *Booking
}

type Booking struct {
	Username       string
	BookingEndTime time.Time
}
