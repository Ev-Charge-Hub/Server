package usecase_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"Ev-Charge-Hub/Server/internal/dto/request"
	"Ev-Charge-Hub/Server/internal/mocks"
	repoModels "Ev-Charge-Hub/Server/internal/repository/models"
	"Ev-Charge-Hub/Server/internal/usecase"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestShowAllStations_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEVStationRepository(ctrl)
	uc := usecase.NewEVStationUsecase(mockRepo)

	mockRepo.EXPECT().FindAllStations(gomock.Any()).Return([]repoModels.EVStationDB{
		{
			ID:       [12]byte{},
			Name:     "Station A",
			Latitude: 13.75,
			Status: repoModels.StationStatusDB{
				IsOpen: true,
			},
		},
	}, nil)

	resp, err := uc.ShowAllStations(context.TODO())

	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "Station A", resp[0].Name)
}

func TestFilterStations_StatusClosed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEVStationRepository(ctrl)
	uc := usecase.NewEVStationUsecase(mockRepo)

	req := request.StationFilterRequest{
		Status: "closed",
	}

	mockRepo.EXPECT().
		FindStations(gomock.Any(), "", "", "", "", gomock.Not(nil)).
		Return([]repoModels.EVStationDB{}, nil)

	_, err := uc.FilterStations(context.TODO(), req)
	assert.NoError(t, err)
}

func TestSetBooking_PastTime_ShouldFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEVStationRepository(ctrl)
	uc := usecase.NewEVStationUsecase(mockRepo)

	req := request.SetBookingRequest{
		ConnectorId:    "CT01",
		Username:       "user1",
		BookingEndTime: "2020-01-01T10:00:00",
	}

	err := uc.SetBooking(context.TODO(), req)
	assert.EqualError(t, err, "booking_end_time must be in the future")
}

func TestSetBooking_UserAlreadyBooked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEVStationRepository(ctrl)
	uc := usecase.NewEVStationUsecase(mockRepo)

	endTime := time.Now().Add(1 * time.Hour).Format("2006-01-02T15:04:05")

	req := request.SetBookingRequest{
		ConnectorId:    "CT02",
		Username:       "user1",
		BookingEndTime: endTime,
	}

	mockRepo.EXPECT().FindBookingsByUserName(gomock.Any(), "user1").Return([]repoModels.BookingDB{
		{Username: "user1", BookingEndTime: time.Now().Add(30 * time.Minute).Format("2006-01-02T15:04:05")},
	}, nil)

	err := uc.SetBooking(context.TODO(), req)
	assert.Contains(t, err.Error(), "user already has an active booking")
}

func TestSetBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEVStationRepository(ctrl)
	uc := usecase.NewEVStationUsecase(mockRepo)

	endTime := time.Now().Add(2 * time.Hour).Format("2006-01-02T15:04:05")

	req := request.SetBookingRequest{
		ConnectorId:    "CT03",
		Username:       "newuser",
		BookingEndTime: endTime,
	}

	// No previous booking
	mockRepo.EXPECT().
		FindBookingsByUserName(gomock.Any(), "newuser").
		Return([]repoModels.BookingDB{}, nil)

	// Connector not booked
	mockRepo.EXPECT().
		FindStationByConnectorID(gomock.Any(), "CT03").
		Return(&repoModels.EVStationDB{
			Connectors: []repoModels.ConnectorDB{
				{ConnectorID: "CT03"},
			},
		}, nil)

	mockRepo.EXPECT().
		SetBooking(gomock.Any(), "CT03", gomock.Any()).
		Return(nil)

	err := uc.SetBooking(context.TODO(), req)
	assert.NoError(t, err)
}

func TestGetStationByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEVStationRepository(ctrl)
	uc := usecase.NewEVStationUsecase(mockRepo)

	mockRepo.EXPECT().
		FindStationByID(gomock.Any(), "station123").
		Return(&repoModels.EVStationDB{
			ID:   primitive.NewObjectID(),
			Name: "Test Station",
		}, nil)

	resp, err := uc.GetStationByID(context.TODO(), request.GetStationByIDRequest{ID: "station123"})
	assert.NoError(t, err)
	assert.Equal(t, "Test Station", resp.Name)
}

func TestGetStationByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEVStationRepository(ctrl)
	uc := usecase.NewEVStationUsecase(mockRepo)

	mockRepo.EXPECT().
		FindStationByID(gomock.Any(), "badID").
		Return(nil, fmt.Errorf("not found"))

	resp, err := uc.GetStationByID(context.TODO(), request.GetStationByIDRequest{ID: "badID"})
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestCreateStation_CallsRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEVStationRepository(ctrl)
	uc := usecase.NewEVStationUsecase(mockRepo)

	req := request.EVStationRequest{
		Name:      "New Station",
		Latitude:  13.7,
		Longitude: 100.5,
		Company:   "EV CO",
		Status:    request.StationStatusRequest{IsOpen: true},
	}

	mockRepo.EXPECT().
		CreateStation(gomock.Any(), gomock.Any()).
		Return(nil)

	err := uc.CreateStation(context.TODO(), req)
	assert.NoError(t, err)
}

func TestEditStation_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockEVStationRepository(ctrl)
	uc := usecase.NewEVStationUsecase(mockRepo)

	req := request.EditStationRequest{
		ID: "invalid_hex_id", // not a valid ObjectID
	}

	resp, err := uc.EditStation(context.TODO(), req)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid ID")
}

func TestRemoveStation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockEVStationRepository(ctrl)
	uc := usecase.NewEVStationUsecase(mockRepo)

	mockRepo.EXPECT().
		RemoveStation(gomock.Any(), "stationXYZ").
		Return(nil)

	err := uc.RemoveStation(context.TODO(), request.RemoveStationRequest{ID: "stationXYZ"})
	assert.NoError(t, err)
}

func TestFilterStations_InvalidStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockEVStationRepository(ctrl)
	uc := usecase.NewEVStationUsecase(mockRepo)

	_, err := uc.FilterStations(context.TODO(), request.StationFilterRequest{Status: "unknown-status"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status value")
}
