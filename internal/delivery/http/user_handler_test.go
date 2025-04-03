package http_test

import (
	deliveryHttp "Ev-Charge-Hub/Server/internal/delivery/http" // ✅ ตั้งชื่อให้ไม่ชน
	"Ev-Charge-Hub/Server/internal/dto/request"
	"Ev-Charge-Hub/Server/internal/dto/response"
	"Ev-Charge-Hub/Server/internal/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupRouterWithUserHandler(mockUsecase *mocks.MockUserUsecaseInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := deliveryHttp.NewUserHandler(mockUsecase)
	router.POST("/register", handler.RegisterUser)
	router.POST("/login", handler.LoginUser)
	return router
}

func TestRegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecaseInterface(ctrl)
	handler := deliveryHttp.NewUserHandler(mockUsecase)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", handler.RegisterUser)

	reqBody := request.RegisterUserRequest{
		Username: "test",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "USER",
	}
	mockUsecase.EXPECT().RegisterUser(gomock.Any(), reqBody).Return(nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
}

func TestLoginUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecaseInterface(ctrl)
	handler := deliveryHttp.NewUserHandler(mockUsecase)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/login", handler.LoginUser)

	reqBody := request.LoginRequest{
		UsernameOrEmail: "test@example.com",
		Password:        "password123",
	}
	expectedResp := &response.LoginResponse{Token: "mocked-jwt-token"}
	mockUsecase.EXPECT().LoginUser(gomock.Any(), reqBody).Return(expectedResp, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	var result response.LoginResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.Equal(t, "mocked-jwt-token", result.Token)
}

func TestRegisterUser_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecaseInterface(ctrl)
	router := setupRouterWithUserHandler(mockUsecase)

	// Missing required fields (invalid JSON)
	body := `{"username": "x"}`

	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "validation_error")
}

func TestLoginUser_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecaseInterface(ctrl)
	router := setupRouterWithUserHandler(mockUsecase)

	loginReq := request.LoginRequest{
		UsernameOrEmail: "user@example.com",
		Password:        "wrongpass",
	}

	mockUsecase.
		EXPECT().
		LoginUser(gomock.Any(), loginReq).
		Return(nil, errors.New("invalid email or password"))

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "invalid email or password")
}

func TestRegisterUser_UsecaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockUserUsecaseInterface(ctrl)
	router := setupRouterWithUserHandler(mockUsecase)

	registerReq := request.RegisterUserRequest{
		Username: "test",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "USER",
	}

	mockUsecase.
		EXPECT().
		RegisterUser(gomock.Any(), registerReq).
		Return(errors.New("email already exists"))

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), "email already exists")
}
