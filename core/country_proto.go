package core

import "github.com/oschwald/geoip2-golang"

type CountryRecord struct {
	Continent *Continent `json:"continent"`
	Country   *Country   `json:"country"`
}

type Continent struct {
	Code      string `json:"code"`
	GeoNameID uint   `json:"geoname_id"`
	Name      string `json:"name"`
}

type Country struct {
	GeoNameID         uint   `json:"geoname_id"`
	IsInEuropeanUnion bool   `json:"is_in_european_union"`
	IsoCode           string `json:"iso_code"`
	Name              string `json:"name"`
}

type MmdbContinent struct {
	Names     map[string]string `maxminddb:"names"`
	Code      string            `maxminddb:"code"`
	GeoNameID uint              `maxminddb:"geoname_id"`
}

type MmdbCountry struct {
	Names             map[string]string `maxminddb:"names"`
	IsoCode           string            `maxminddb:"iso_code"`
	GeoNameID         uint              `maxminddb:"geoname_id"`
	IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
}

func FromMaxmindContinentField(c MmdbContinent, maxmindLang string) *Continent {
	return &Continent{
		Code:      c.Code,
		GeoNameID: c.GeoNameID,
		Name:      maxmindGetMLNameFallback(c.Names, maxmindLang),
	}
}

func FromMaxmindCountryField(c MmdbCountry, maxmindLang string) *Country {
	return &Country{
		GeoNameID:         c.GeoNameID,
		IsInEuropeanUnion: c.IsInEuropeanUnion,
		IsoCode:           c.IsoCode,
		Name:              maxmindGetMLNameFallback(c.Names, maxmindLang),
	}
}

func FromMaxmindCountry(country *geoip2.Country, maxmindLang string) *CountryRecord {
	ret := CountryRecord{
		Continent: FromMaxmindContinentField(country.Continent, maxmindLang),
		Country:   FromMaxmindCountryField(country.Country, maxmindLang),
	}
	return &ret
}
