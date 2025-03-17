package request

type SetBookingRequest struct {
	ID             string `json:"id" binding:"required"`
	Username       string `json:"username" binding:"required"`
	BookingEndTime string `json:"booking_end_time" binding:"required,datetime=2006-01-02T15:04:05"`
}


