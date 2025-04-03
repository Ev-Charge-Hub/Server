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
	FindStationByConnectorID(ctx context.Context, connectorID string) (*models.EVStationDB, error)
	FindBookingByUserName(ctx context.Context, userName string) (*models.BookingDB, error)
	FindBookingsByUserName(ctx context.Context, username string) ([]models.BookingDB, error)
	FindStationByUserName(ctx context.Context, userName string) (*models.EVStationDB, error)
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
	isOpen *bool) ([]models.EVStationDB, error) {
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

func (repo *evStationRepository) CreateStation(ctx context.Context, station models.EVStationDB) error {
	station.ID = primitive.NewObjectID() // Only Generate New ID, No Timestamps Added
	_, err := repo.collection.InsertOne(ctx, station)
	return err
}

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

func (repo *evStationRepository) SetBooking(ctx context.Context, connector_id string, booking models.BookingDB) error {
	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to find stations: %v", err)
	}

	var stations []models.EVStationDB
	if err := cursor.All(ctx, &stations); err != nil {
		return fmt.Errorf("failed to iterate over stations: %v", err)
	}

	connectorFound := false

	for _, station := range stations {
		for i, connector := range station.Connectors {
			if connector.ConnectorID == connector_id {
				// Found the connector, update booking
				station.Connectors[i].Booking = &booking

				filter := bson.M{
					"_id":                     station.ID,
					"connectors.connector_id": connector_id,
				}

				update := bson.M{
					"$set": bson.M{
						"connectors.$.booking": booking,
					},
				}

				// Update the station in the database
				result, err := repo.collection.UpdateOne(ctx, filter, update)
				if err != nil {
					return fmt.Errorf("failed to update booking: %v", err)
				}

				// Check if the connector was matched and updated
				if result.MatchedCount == 0 {
					return fmt.Errorf("no matching connector found to update booking")
				}

				// Mark that we have found the connector
				connectorFound = true
				break
			}
		}
		if connectorFound {
			break
		}
	}

	// If connector id was not found in any of the stations
	if !connectorFound {
		return fmt.Errorf("connector id %s not found", connector_id)
	}

	return nil
}

// func (repo *evStationRepository) SetBooking(ctx context.Context, connector_id string, booking models.BookingDB) error {
// 	// สร้าง filter เพื่อหาสถานีที่มี connector_id ตรงกัน
// 	filter := bson.M{
// 		"connectors.connector_id": connector_id,
// 	}
//
// 	// สร้างข้อมูลที่จะอัปเดต (เพียงแค่ booking)
// 	update := bson.M{
// 		"$set": bson.M{
// 			"connectors.$.booking": booking, // อัปเดตเฉพาะ booking ใน connector ที่ตรง
// 		},
// 	}
//
// 	// ทำการอัปเดตใน MongoDB
// 	result, err := repo.collection.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		return fmt.Errorf("ไม่สามารถอัปเดต booking ได้: %v", err)
// 	}
//
// 	// ตรวจสอบว่ามีการอัปเดตจริงหรือไม่
// 	if result.MatchedCount == 0 {
// 		return fmt.Errorf("ไม่พบ connector ที่มี ID %s", connector_id)
// 	}
//
// 	return nil
// }

func (repo *evStationRepository) FindBookingByUserName(ctx context.Context, userName string) (*models.BookingDB, error) {
	filter := bson.M{"connectors.booking.username": userName}

	var station models.EVStationDB
	err := repo.collection.FindOne(ctx, filter).Decode(&station)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			return nil, fmt.Errorf("no booking found for user name %s", userName)
		}
		return nil, fmt.Errorf("error finding station: %v", err)
	}

	for _, connector := range station.Connectors {
		if connector.Booking != nil && connector.Booking.Username == userName {
			return connector.Booking, nil
		}
	}
	return nil, fmt.Errorf("no booking found for user name %s", userName)
}

func (repo *evStationRepository) FindBookingsByUserName(ctx context.Context, username string) ([]models.BookingDB, error) {
	// Filter: Find all stations with any connector having a booking of this username
	filter := bson.M{"connectors.booking.username": username}

	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error querying stations: %v", err)
	}
	defer cursor.Close(ctx)

	var bookings []models.BookingDB

	for cursor.Next(ctx) {
		var station models.EVStationDB
		if err := cursor.Decode(&station); err != nil {
			return nil, fmt.Errorf("error decoding station: %v", err)
		}

		for _, connector := range station.Connectors {
			if connector.Booking != nil && connector.Booking.Username == username {
				bookings = append(bookings, *connector.Booking)
			}
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	if len(bookings) == 0 {
		return nil, fmt.Errorf("no bookings found for username %s", username)
	}

	return bookings, nil
}

func (repo *evStationRepository) FindStationByConnectorID(ctx context.Context, connectorID string) (*models.EVStationDB, error) {
	filter := bson.M{
		"connectors.connector_id": connectorID,
	}

	var station models.EVStationDB
	err := repo.collection.FindOne(ctx, filter).Decode(&station)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no station found with connector id %s", connectorID)
		}
		return nil, fmt.Errorf("error finding station: %v", err)
	}

	return &station, nil
}

// func (repo *evStationRepository) FindStationByConnectorID(ctx context.Context, connector_id string) (*models.EVStationDB, error) {
// 	// ใช้ elemMatch เพื่อให้แม่นยำในการค้นหา
// 	filter := bson.M{
// 		"connectors.connector_id": connector_id,
// 	}
//
// 	var station models.EVStationDB
// 	err := repo.collection.FindOne(ctx, filter).Decode(&station)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, fmt.Errorf("ไม่พบสถานีที่มี connector id %s", connector_id)
// 		}
// 		return nil, fmt.Errorf("เกิดข้อผิดพลาดขณะค้นหาสถานี: %v", err)
// 	}
//
// 	return &station, nil
// }

func (repo *evStationRepository) FindStationByUserName(ctx context.Context, userName string) (*models.EVStationDB, error) {
	filter := bson.M{"connectors.booking.username": userName}

	var station models.EVStationDB
	err := repo.collection.FindOne(ctx, filter).Decode(&station)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no station found for user name %s", userName)
		}
		return nil, fmt.Errorf("error finding station: %v", err)
	}

	// กรองเฉพาะ connector ที่ user คนนี้จองไว้
	var filteredConnectors []models.ConnectorDB
	for _, c := range station.Connectors {
		if c.Booking != nil && c.Booking.Username == userName {
			filteredConnectors = append(filteredConnectors, c)
		}
	}
	station.Connectors = filteredConnectors

	return &station, nil
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
