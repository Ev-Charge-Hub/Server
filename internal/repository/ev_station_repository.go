package repository

import (
	"Ev-Charge-Hub/Server/internal/constants"
	"Ev-Charge-Hub/Server/internal/repository/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EVStationRepository interface {
    FindStations(ctx context.Context, company string, stationType string, search string, plugName string, isOpen *bool) ([]models.EVStationDB, error)
	FindAllStations(ctx context.Context) ([]models.EVStationDB, error)
	FindStationByID(ctx context.Context, id string) (*models.EVStationDB, error)
	// CreateStation(ctx context.Context, station models.EVStationDB)
	// EditStation(ctx context.Context, id string , station models.EVStationDB)
	// RemoveStation(ctx context.Context,id string)
}

type evStationRepository struct {
	collection *mongo.Collection
}

func NewEVStationRepository(db *mongo.Database) EVStationRepository {
	return &evStationRepository{collection: db.Collection("ev_station")}
}

func (repo *evStationRepository) FindStations(
    ctx context.Context, 
    company string, 
    stationType string, 
    search string, 
    plugName string,
    isOpen *bool,
) ([]models.EVStationDB, error) {	
	filter := bson.M{}
	// กรอง Company และ Search ตามปกติ
	if company != "" {
		filter["company"] = company
	}
	if search != "" {
		filter["name"] = bson.M{"$regex": search, "$options": "i"}
	}

	if isOpen != nil {
		filter["status.is_open"] = *isOpen
	}

	// ดึงข้อมูล Ens ทั้งหมดที่ตรงกับ filterV Station
	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var stations []models.EVStationDB
	if err := cursor.All(ctx, &stations); err != nil {
		return nil, err
	}

	// map Connectors each Station (type)
	if stationType != "" {
		for i := range stations {
			stations[i].Connectors = filterConnectorsByType(stations[i].Connectors, stationType)
		}
	}

	// map Connectors each Station (plugName)
	if plugName != "" {
		for i := range stations {
			stations[i].Connectors = filterConnectorsByPlugName(stations[i].Connectors, plugName)
		}
	}

	return stations, nil
}

func (repo *evStationRepository) FindAllStations(ctx context.Context) ([]models.EVStationDB, error) {
	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var stations []models.EVStationDB
	if err := cursor.All(ctx, &stations); err != nil {
		return nil, err
	}
	return stations, nil
}

func (repo *evStationRepository) FindStationByID(ctx context.Context, id string) (*models.EVStationDB, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var station models.EVStationDB
	err = repo.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&station)
	if err != nil {
		return nil, err
	}

	return &station, nil
}

func filterConnectorsByType(connectors []models.ConnectorDB, stationType string) []models.ConnectorDB {
	var filtered []models.ConnectorDB

	connectorType := constants.ConnectorType(stationType)

	for _, connector := range connectors {
		if connector.Type == connectorType {
			filtered = append(filtered, connector)
		}
	}
	return filtered
}

func filterConnectorsByPlugName(connectors []models.ConnectorDB, plugName string) []models.ConnectorDB {
	var filtered []models.ConnectorDB

	connectorPlugName := constants.PlugName(plugName)

	for _, connector := range connectors {
		if connector.PlugName == connectorPlugName {
			filtered = append(filtered, connector)
		}
	}
	return filtered
}