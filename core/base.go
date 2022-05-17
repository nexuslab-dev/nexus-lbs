package core

import "github.com/oschwald/geoip2-golang"

type GeoQuery struct {
	DBPath string
	geoip  *geoip2.Reader
}

func New(DBPath string) (*GeoQuery, error) {
	gi, err := geoip2.Open(DBPath)
	if err != nil {
		return nil, err
	}
	c := &GeoQuery{
		DBPath: DBPath,
		geoip:  gi,
	}
	return c, nil
}

func maxmindGetMLNameFallback(theMap map[string]string, targetLang string) string {
	if name, ok := theMap[targetLang]; ok {
		return name
	}
	return theMap["en"]
}
