package configs

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB เชื่อมต่อกับ MongoDB และคืนค่า Database Object
func ConnectDB() *mongo.Database {
	// กำหนด URI ของ MongoDB
	mongoURI := "mongodb+srv://setthanan50:Admin1234@ev-charge-hub-db.lgqxg.mongodb.net/"

	// ตั้งค่า Timeout สำหรับการเชื่อมต่อ
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// สร้างและเชื่อมต่อ MongoDB Client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// ตรวจสอบการเชื่อมต่อ
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB successfully.")
	return client.Database("ev_charge_hub")
}
