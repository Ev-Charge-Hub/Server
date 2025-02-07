package user

import (
	"Ev-Charge-Hub/Server/dto/request"
	"Ev-Charge-Hub/Server/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserHandlerInterface สำหรับการ Mock หรือเปลี่ยน Implementation
type UserHandlerInterface interface {
	RegisterUser(c *gin.Context)
	LoginUser(c *gin.Context)
}

// class userHandler
type userHandler struct {
	userUsecase usecase.UserUsecaseInterface
}

// init attribute of class userHandler
func NewUserHandler(userUsecase usecase.UserUsecaseInterface) UserHandlerInterface {
	return &userHandler{userUsecase: userUsecase} // return Pointer(*) -> for pass by reference
}

// LoginUser implements UserHandlerInterface.
func (h *userHandler) RegisterUser(c *gin.Context) {
	var req request.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.userUsecase.RegisterUser(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// RegisterUser implements UserHandlerInterface.
func (h *userHandler) LoginUser(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.userUsecase.LoginUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token.Token})
}