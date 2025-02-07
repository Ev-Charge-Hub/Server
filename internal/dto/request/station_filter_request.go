package request

type StationFilterRequest struct {
	Company string `form:"company"`
	Type    string `form:"type"`
	Search  string `form:"search"`
}
