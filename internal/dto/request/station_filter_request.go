package request

type StationFilterRequest struct {
	Company  string `form:"company"`
	Type     string `form:"type"`
	PlugName string `form:"plug_name"`
	Search   string `form:"search"`
	Status	 string `form:"status"`
}
