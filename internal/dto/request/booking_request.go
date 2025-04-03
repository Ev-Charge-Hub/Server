package request

type SetBookingRequest struct {
	ConnectorId    string `json:"connector_id" binding:"required"`
	Username       string `json:"username" binding:"required"`
	BookingEndTime string `json:"booking_end_time" binding:"required,datetime=2006-01-02T15:04:05"`
}

type GetBookingRequest struct {
	Username string `json:"username" binding:"required"`
}

type GetBookingsRequest struct {
	Username string `json:"username" binding:"required"`
}