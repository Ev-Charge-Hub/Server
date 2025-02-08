package constants

type ConnectorType string
type PlugName string

const (
	AC ConnectorType = "AC"
	DC ConnectorType = "DC"
)

const (
	J1772    PlugName = "J1772 TYPE 1"
	Type2    PlugName = "TYPE 2"
	CCSType1 PlugName = "CCS TYPE 1"
	CCSType2 PlugName = "CCS TYPE 2"
	CHAdeMO  PlugName = "CHAdeMO"
	GBTAC    PlugName = "GB/T AC"
	GBTDC    PlugName = "GB/T DC"
)
