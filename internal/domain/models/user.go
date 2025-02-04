package models

import (
	"time"
)

type User struct {
	ID        string    
	Username  string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// type User struct {
// 	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
// 	Username  string             `json:"username" bson:"username"`
// 	Email     string             `json:"email" bson:"email"`
// 	Password  string             `json:"password" bson:"password"`
// 	Role      string             `json:"role" bson:"role"`
// 	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
// 	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
// }

// json: key name when Post to client(ถูกแปลงเป็น JSON)
// bson: key name when Get from DB(ถูกใช้เก็บหรือดึงข้อมูลจาก MongoDB)