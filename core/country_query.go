package core

import (
	"fmt"
	"net"
)

func (a *GeoQuery) Country(ip string, maxmindLang string) (*CountryRecord, error) {
	if a.geoip == nil {
		return nil, fmt.Errorf("err nil geoip2.Reader")
	}
	if ip == "" {
		return nil, fmt.Errorf("err ip param empty")
	}
	netip := net.ParseIP(ip)
	country, err := a.geoip.Country(netip)
	if err != nil {
		return nil, fmt.Errorf("geoip2 get country failed, err=%v", err)
	}
	if country == nil {
		return nil, fmt.Errorf("geoip2 get country success but nil country object")
	}
	if maxmindLang == "" {
		maxmindLang = "en"
	}

	ret := FromMaxmindCountry(country, maxmindLang)
	return ret, nil
}
