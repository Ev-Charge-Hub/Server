package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserDB struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Role      string             `bson:"role"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// json: key name when Post to client(ถูกแปลงเป็น JSON)
// bson: key name when Get from DB(ถูกใช้เก็บหรือดึงข้อมูลจาก MongoDB)