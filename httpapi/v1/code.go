package v1

type Code string

const (
	OK       Code = "OK"
	ErrParam Code = "ERR.PARAM"
	ErrGeoip Code = "ERR.GEOIP"
)
