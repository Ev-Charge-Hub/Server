package models

import (
	"time"
)

type UserModel struct {
	ID        string    
	Username  string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// json: key name when Post to client(ถูกแปลงเป็น JSON)
// bson: key name when Get from DB(ถูกใช้เก็บหรือดึงข้อมูลจาก MongoDB)