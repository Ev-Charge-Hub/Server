package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"Ev-Charge-Hub/Server/internal/constants"
)

type EVStationDB struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	StationID  string             `bson:"station_id"`
	Name       string             `bson:"name"`
	Latitude   float64            `bson:"latitude"`
	Longitude  float64            `bson:"longitude"`
	Company    string             `bson:"company"`
	Status     StationStatusDB    `bson:"status"`
	Connectors []ConnectorDB      `bson:"connectors"`
}

type StationStatusDB struct {
	OpenHours  string `bson:"open_hours"`
	CloseHours string `bson:"close_hours"`
	IsOpen     bool   `bson:"is_open"`
}

type ConnectorDB struct {
	ConnectorID   string                  `bson:"connector_id"`
	Type          constants.ConnectorType `bson:"type"`
	PlugName	  constants.PlugName  	  `bson:"plug_name"`
	PricePerUnit  float64                 `bson:"price_per_unit"`
	PowerOutput   int                     `bson:"power_output"`
	IsAvailable   bool                    `bson:"is_available"`
}
