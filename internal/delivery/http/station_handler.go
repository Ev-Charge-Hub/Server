package http

import (
	"Ev-Charge-Hub/Server/internal/dto/request"
	"Ev-Charge-Hub/Server/internal/repository/models"
	"Ev-Charge-Hub/Server/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Class
type EVStationHandler struct {
	stationUsecase usecase.EVStationUsecase
}

// Init Class
func NewEVStationHandler(usecase usecase.EVStationUsecase) *EVStationHandler {
	return &EVStationHandler{stationUsecase: usecase}
}

func (h *EVStationHandler) FilterStations(c *gin.Context) {
	var filterRequest request.StationFilterRequest

	// Query and save data from parameter
	if err := c.ShouldBindQuery(&filterRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stations, err := h.stationUsecase.FilterStations(c, filterRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stations)
}

func (h *EVStationHandler) ShowAllStations(c *gin.Context) {
	stations, err := h.stationUsecase.ShowAllStations(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stations)
}

func (h *EVStationHandler) GetStationByID(c *gin.Context) {
	id := c.Param("id")

	station, err := h.stationUsecase.GetStationByID(c.Request.Context(), request.GetStationByIDRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Station not found"})
		return
	}

	c.JSON(http.StatusOK, station)
}

func (h *EVStationHandler) SetBooking(c *gin.Context) {
	var bookingReq request.SetBookingRequest

	// Log Request Data เพื่อเช็คข้อมูลก่อนเช็ค Validation
	if err := c.ShouldBindJSON(&bookingReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	//  Call Usecase
	err := h.stationUsecase.SetBooking(c.Request.Context(), bookingReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking successfully added"})
}

func (h *EVStationHandler) CreateStation(c *gin.Context) {
	var stationRequest request.EVStationRequest
	if err := c.ShouldBindJSON(&stationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid station data"})
		return
	}

	// Validate request
	if err := validate.Struct(stationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_error": err.Error()})
		return
	}

	// ส่งต่อ request ไป Usecase เลย
	if err := h.stationUsecase.CreateStation(c.Request.Context(), stationRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Station created successfully"})
}

func (h *EVStationHandler) EditStation(c *gin.Context) {
	id := c.Param("id")

	var stationReq request.EVStationRequest
	if err := c.ShouldBindJSON(&stationReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid station data",
			"detail": err.Error(),
		})
		return
	}

	editReq := request.EditStationRequest{
		ID:         id,
		Name:       &stationReq.Name,
		Latitude:   &stationReq.Latitude,
		Longitude:  &stationReq.Longitude,
		Company:    &stationReq.Company,
		Status:     &stationReq.Status,
		Connectors: &stationReq.Connectors,
	}

	updated, err := h.stationUsecase.EditStation(c.Request.Context(), editReq)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Station updated successfully",
		"station": updated,
	})
}


func (h *EVStationHandler) RemoveStation(c *gin.Context) {
	id := c.Param("id")

	err := h.stationUsecase.RemoveStation(c.Request.Context(), request.RemoveStationRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Station not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Station removed successfully"})
}

func (h *EVStationHandler) GetBookingByUserName(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	booking, err := h.stationUsecase.GetBookingByUserName(c.Request.Context(), request.GetBookingRequest{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (h *EVStationHandler) GetBookingsByUserName(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	bookings, err := h.stationUsecase.GetBookingsByUserName(c.Request.Context(), request.GetBookingsRequest{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *EVStationHandler) GetStationByConnectorID(c *gin.Context) {
	connectorID := c.Param("connector_id")

	station, err := h.stationUsecase.GetStationByConnectorID(c.Request.Context(), request.GetStationByConnectorIDRequest{ConnectorId: connectorID})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Station not found"})
		return
	}

	c.JSON(http.StatusOK, station)
}

func (h *EVStationHandler) GetStationByUserName(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	station, err := h.stationUsecase.GetStationByUserName(c.Request.Context(), request.GetStationByUsernameRequest{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, station)
}

func mapRequestToModel(req request.EVStationRequest) models.EVStationDB {
	var connectors []models.ConnectorDB
	for _, c := range req.Connectors {
		var booking *models.BookingDB
		if c.Booking != nil {
			booking = &models.BookingDB{
				Username:       c.Booking.Username,
				BookingEndTime: c.Booking.BookingEndTime,
			}
		}

		connectors = append(connectors, models.ConnectorDB{
			Type:         c.Type,
			PlugName:     c.PlugName,
			PricePerUnit: c.PricePerUnit,
			PowerOutput:  c.PowerOutput,
			Booking:      booking,
		})
	}

	return models.EVStationDB{
		Name:      req.Name,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Company:   req.Company,
		Status: models.StationStatusDB{
			OpenHours:  req.Status.OpenHours,
			CloseHours: req.Status.CloseHours,
			IsOpen:     req.Status.IsOpen,
		},
		Connectors: connectors,
	}
}
