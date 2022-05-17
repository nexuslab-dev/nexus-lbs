package core

import (
	"fmt"
	"net"
)

func (a *GeoQuery) City(ip string, maxmindLang string) (*CityRecord, error) {
	if a.geoip == nil {
		return nil, fmt.Errorf("err nil geoip2.Reader")
	}
	if ip == "" {
		return nil, fmt.Errorf("err ip param empty")
	}
	netip := net.ParseIP(ip)
	city, err := a.geoip.City(netip)
	if err != nil {
		return nil, fmt.Errorf("geoip2 get City failed, err=%v", err)
	}
	if city == nil {
		return nil, fmt.Errorf("geoip2 get CityName success but nil city object")
	}
	if maxmindLang == "" {
		maxmindLang = "en"
	}

	ret := FromMaxmindCity(city, maxmindLang)
	return ret, nil
}
