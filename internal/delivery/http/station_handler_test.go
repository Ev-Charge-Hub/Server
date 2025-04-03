package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	deliveryHttp "Ev-Charge-Hub/Server/internal/delivery/http"
	"Ev-Charge-Hub/Server/internal/dto/request"
	"Ev-Charge-Hub/Server/internal/dto/response"
	"Ev-Charge-Hub/Server/internal/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupRouterWithStationHandler(mockUsecase *mocks.MockEVStationUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := deliveryHttp.NewEVStationHandler(mockUsecase)

	r.GET("/stations", handler.ShowAllStations)
	r.GET("/stations/filter", handler.FilterStations)
	r.POST("/stations", handler.CreateStation)
	r.POST("/stations/booking", handler.SetBooking)
	r.GET("/stations/:id", handler.GetStationByID)
	r.PUT("/stations/:id", handler.EditStation)
	r.DELETE("/stations/:id", handler.RemoveStation)

	return r
}

func TestShowAllStations_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockEVStationUsecase(ctrl)
	router := setupRouterWithStationHandler(mockUsecase)

	mockUsecase.EXPECT().ShowAllStations(gomock.Any()).Return([]response.EVStationResponse{
		{Name: "StationA"},
	}, nil)

	req := httptest.NewRequest("GET", "/stations", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "StationA")
}

func TestFilterStations_InvalidQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockEVStationUsecase(ctrl)

	// ✅ Expect the call with invalid status, return error
	mockUsecase.
		EXPECT().
		FilterStations(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("invalid status value: unknown"))

	router := setupRouterWithStationHandler(mockUsecase)

	req := httptest.NewRequest("GET", "/stations/filter?status=unknown", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), "invalid status value")
}

func TestCreateStation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockEVStationUsecase(ctrl)
	router := setupRouterWithStationHandler(mockUsecase)

	body := request.EVStationRequest{
		Name:      "Updated Station",
		Latitude:  13.5,
		Longitude: 100.5,
		Company:   "Updated Co",
		Status: request.StationStatusRequest{
			OpenHours:  "09:00",
			CloseHours: "19:00",
			IsOpen:     true,
		},
		Connectors: []request.ConnectorRequest{
			{
				Type:         "DC",
				PlugName:     "Type 2",
				PricePerUnit: 10,
				PowerOutput:  22,
			},
		},
	}

	mockUsecase.EXPECT().CreateStation(gomock.Any(), body).Return(nil)

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/stations", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Station created successfully")
}

func TestSetBooking_Fail_InvalidFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockEVStationUsecase(ctrl)
	router := setupRouterWithStationHandler(mockUsecase)

	body := `{"connector_id": "xxx"}` // incomplete and invalid
	req := httptest.NewRequest("POST", "/stations/booking", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestSetBooking_UsecaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockEVStationUsecase(ctrl)
	router := setupRouterWithStationHandler(mockUsecase)

	body := request.SetBookingRequest{
		ConnectorId:    "abc123",
		Username:       "john",
		BookingEndTime: "2025-12-31T10:00:00",
	}

	mockUsecase.
		EXPECT().
		SetBooking(gomock.Any(), body).
		Return(errors.New("booking failed"))

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/stations/booking", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), "booking failed")
}

func TestEditStation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockEVStationUsecase(ctrl)
	router := setupRouterWithStationHandler(mockUsecase)

	id := "abc123"
	reqBody := request.EVStationRequest{
		Name:      "Updated Station",
		Latitude:  13.5,
		Longitude: 100.5,
		Company:   "Updated Co",
		Status: request.StationStatusRequest{
			OpenHours:  "09:00",
			CloseHours: "19:00",
			IsOpen:     true,
		},
		Connectors: []request.ConnectorRequest{
			{
				Type:         "DC",
				PlugName:     "Type 2",
				PricePerUnit: 10,
				PowerOutput:  22,
			},
		},
	}

	mockUsecase.
		EXPECT().
		EditStation(gomock.Any(), gomock.Any()).
		Return(&response.EVStationResponse{Name: "Updated Station"}, nil)

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/stations/"+id, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Updated Station")
}

func TestGetStationByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockEVStationUsecase(ctrl)
	router := setupRouterWithStationHandler(mockUsecase)

	mockUsecase.
		EXPECT().
		GetStationByID(gomock.Any(), request.GetStationByIDRequest{ID: "notfound"}).
		Return(nil, errors.New("station not found"))

	req := httptest.NewRequest("GET", "/stations/notfound", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Contains(t, resp.Body.String(), "Station not found")
}

func TestRemoveStation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockEVStationUsecase(ctrl)
	router := setupRouterWithStationHandler(mockUsecase)

	mockUsecase.
		EXPECT().
		RemoveStation(gomock.Any(), request.RemoveStationRequest{ID: "abc123"}).
		Return(nil)

	req := httptest.NewRequest("DELETE", "/stations/abc123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Station removed successfully")
}
func TestGetStationByUserName_EmptyParam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockEVStationUsecase(ctrl)
	router := gin.Default()
	handler := deliveryHttp.NewEVStationHandler(mockUsecase)

	router.GET("/stations/user/", handler.GetStationByUserName)

	req := httptest.NewRequest("GET", "/stations/user/", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "username is required")
}
