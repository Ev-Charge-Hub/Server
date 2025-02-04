package request

// type RegisterUserRequest struct {
// 	Username string `json:"username" validate:"required,min=4,max=20"`
// 	Email    string `json:"email" validate:"required,email"`
// 	Password string `json:"password" validate:"required,min=6"`
// 	Role     string `json:"role" validate:"required,eq=ADMIN|eq=USER"`
// }

// type LoginUserRequest struct {
// 	UsernameOrEmail string `json:"username_or_email" validate:"required,min=4,max=20"`
// 	Password string `json:"password" validate:"required,min=6"`
// }

// json: key name when Post to client(ถูกแปลงเป็น JSON)
// bson: key name when Get from DB(ถูกใช้เก็บหรือดึงข้อมูลจาก MongoDB)