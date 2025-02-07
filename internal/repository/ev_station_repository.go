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
	FindStations(ctx context.Context, company string, stationType string, search string) ([]models.EVStationDB, error)
	FindAllStations(ctx context.Context) ([]models.EVStationDB, error)
	FindStationByID(ctx context.Context, id string) (*models.EVStationDB, error)
}

type evStationRepository struct {
	collection *mongo.Collection
}

func NewEVStationRepository(db *mongo.Database) EVStationRepository {
	return &evStationRepository{collection: db.Collection("ev_station")}
}

// func (repo *evStationRepository) FindStations(ctx context.Context, company string, stationType string, search string) ([]models.EVStationDB, error) {
// 	filter := bson.M{}

// 	if company != "" {
// 		filter["company"] = company
// 	}
// 	if stationType != "" {
// 		filter["connectors.type"] = stationType
// 	}
// 	if search != "" {
// 		filter["name"] = bson.M{"$regex": search, "$options": "i"} // ค้นหาด้วยชื่อที่คล้ายกัน
// 	}

// 	cursor, err := repo.collection.Find(ctx, filter)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var stations []models.EVStationDB
// 	if err = cursor.All(ctx, &stations); err != nil {
// 		return nil, err
// 	}

// 	return stations, nil
// }

func (repo *evStationRepository) FindStations(ctx context.Context, company string, stationType string, search string) ([]models.EVStationDB, error) {
	filter := bson.M{}

	// กรอง Company และ Search ตามปกติ
	if company != "" {
		filter["company"] = company
	}
	if search != "" {
		filter["name"] = bson.M{"$regex": search, "$options": "i"}
	}

	// ดึงข้อมูล EV Stations ทั้งหมดที่ตรงกับ filter
	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var stations []models.EVStationDB
	if err := cursor.All(ctx, &stations); err != nil {
		return nil, err
	}

	// กรอง Connectors ภายในแต่ละ Station
	if stationType != "" {
		for i := range stations {
			stations[i].Connectors = filterConnectorsByType(stations[i].Connectors, stationType)
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

	// แปลง stationType (string) เป็น ConnectorType
	connectorType := constants.ConnectorType(stationType)

	for _, connector := range connectors {
		if connector.Type == connectorType {
			filtered = append(filtered, connector)
		}
	}
	return filtered
}
