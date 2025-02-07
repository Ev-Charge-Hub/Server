package models

import "Ev-Charge-Hub/Server/internal/constants"

type EVStation struct {
	ID         string
	StationID  string
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
	ConnectorID   string
	Type          constants.ConnectorType
	PricePerUnit  float64
	PowerOutput   int
	IsAvailable   bool
}