package http_test

import (
	deliveryHttp "Ev-Charge-Hub/Server/internal/delivery/http"
	response "Ev-Charge-Hub/Server/internal/dto/response"
	"Ev-Charge-Hub/Server/internal/mocks"
	"Ev-Charge-Hub/Server/middleware"
	"Ev-Charge-Hub/Server/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mock route ที่ถูก protect ด้วย middleware
func setupProtectedStationRoute(mockUsecase *mocks.MockEVStationUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := deliveryHttp.NewEVStationHandler(mockUsecase)

	protected := r.Group("/stations")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("", handler.ShowAllStations)

	return r
}

func TestAuthMiddleware_NoTokenProvided(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockEVStationUsecase(ctrl)
	router := setupProtectedStationRoute(mockUsecase)

	req := httptest.NewRequest("GET", "/stations", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Authorization header missing")
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockEVStationUsecase(ctrl)
	router := setupProtectedStationRoute(mockUsecase)

	req := httptest.NewRequest("GET", "/stations", nil)
	req.Header.Set("Authorization", "InvalidFormatToken")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid authorization header format")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockEVStationUsecase(ctrl)
	router := setupProtectedStationRoute(mockUsecase)

	req := httptest.NewRequest("GET", "/stations", nil)
	req.Header.Set("Authorization", "Bearer invalid.jwt.token")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid token")
}

func TestAuthMiddleware_ValidToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockEVStationUsecase(ctrl)

	mockUsecase.EXPECT().
		ShowAllStations(gomock.Any()).
		Return([]response.EVStationResponse{{Name: "Authorized Station"}}, nil)

	router := setupProtectedStationRoute(mockUsecase)

	// ✅ Generate valid token
	token, _ := utils.CreateToken("u123", "testuser", "USER")

	req := httptest.NewRequest("GET", "/stations", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Authorized Station")
}
