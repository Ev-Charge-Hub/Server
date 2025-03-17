package configs

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Database {
	// กำหนด URI ของ MongoDB

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}
	mongoURI := os.Getenv("MONGO_URI")

	if mongoURI == "" {
		log.Fatalf("MONGO_URI is not set in .env")
	}
	// ตั้งค่า Timeout สำหรับการเชื่อมต่อ
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
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
