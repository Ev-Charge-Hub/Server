package repository

import (
	"Ev-Charge-Hub/Server/internal/constants"
	"Ev-Charge-Hub/Server/internal/repository/models"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EVStationRepository interface {
	FindStations(ctx context.Context, company string, stationType string, search string, plugName string, isOpen *bool) ([]models.EVStationDB, error)
	FindAllStations(ctx context.Context) ([]models.EVStationDB, error)
	FindStationByID(ctx context.Context, id string) (*models.EVStationDB, error)
	CreateStation(ctx context.Context, station models.EVStationDB) error
	EditStation(ctx context.Context, id string, station models.EVStationDB) error
	RemoveStation(ctx context.Context, id string) error
	SetBooking(ctx context.Context, id string, booking models.BookingDB) error
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

// 🟢 Create New Station
func (repo *evStationRepository) CreateStation(ctx context.Context, station models.EVStationDB) error {
	station.ID = primitive.NewObjectID() // Only Generate New ID, No Timestamps Added
	_, err := repo.collection.InsertOne(ctx, station)
	return err
}

// 🟡 Edit Station Details
func (repo *evStationRepository) EditStation(ctx context.Context, id string, station models.EVStationDB) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := repo.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": station},
	)

	if result.MatchedCount == 0 {
		return errors.New("station not found")
	}

	return err
}

// 🔴 Remove Station
func (repo *evStationRepository) RemoveStation(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := repo.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if result.DeletedCount == 0 {
		return errors.New("station not found")
	}

	return err
}
// SetBooking เพิ่มข้อมูลการจองให้กับสถานีชาร์จไฟฟ้า
func (repo *evStationRepository) SetBooking(ctx context.Context, id string, booking models.BookingDB) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid station ID")
	}

	filter := bson.M{
		"_id":          objectID,
		"connectors.0": bson.M{"$exists": true}, // Ensure connectors array exists
	}

	update := bson.M{
		"$set": bson.M{
			"connectors.0.booking": booking, // Set booking for the first connector
		},
	}

	result, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("station or connector not found")
	}

	return nil
}

// 🔍 Utility Function - Filter Connectors by Type
func filterConnectorsByType(connectors []models.ConnectorDB, stationType string) []models.ConnectorDB {
	var filtered []models.ConnectorDB

	connectorType := constants.ConnectorType(stationType)

	for _, connector := range connectors {
		if connector.Type == connectorType {
			// Auto-clear expired bookings
			if connector.Booking != nil && isBookingExpired(connector.Booking.BookingEndTime) {
				connector.Booking = nil // เคลียร์การจองที่หมดอายุโดยการตั้งค่าเป็น nil
			}
			filtered = append(filtered, connector)
		}
	}
	return filtered
}

// 🔍 Utility Function - Filter Connectors by Plug Name
func filterConnectorsByPlugName(connectors []models.ConnectorDB, plugName string) []models.ConnectorDB {
	var filtered []models.ConnectorDB

	connectorPlugName := constants.PlugName(plugName)

	for _, connector := range connectors {
		if connector.PlugName == connectorPlugName {
			// Auto-clear expired bookings
			if connector.Booking != nil && isBookingExpired(connector.Booking.BookingEndTime) {
				connector.Booking = nil // เคลียร์การจองที่หมดอายุโดยการตั้งค่าเป็น nil
			}
			filtered = append(filtered, connector)
		}
	}
	return filtered
}

// 🔍 Booking Expiry Check Function
func isBookingExpired(bookingEndTime string) bool {
	// Convert string time to `time.Time`
	parsedTime, err := time.Parse("2006-01-02T15:04:05", bookingEndTime)
	if err != nil {
		return true // ถ้าแปลงเวลาไม่ได้ ให้ถือว่าหมดอายุ
	}

	return time.Now().After(parsedTime)
}
