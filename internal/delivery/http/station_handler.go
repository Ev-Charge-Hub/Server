package http

import (
	"Ev-Charge-Hub/Server/internal/dto/request"
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

	station, err := h.stationUsecase.GetStationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Station not found"})
		return
	}

	c.JSON(http.StatusOK, station)
}