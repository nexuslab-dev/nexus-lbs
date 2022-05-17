package core

import "github.com/oschwald/geoip2-golang"

type CityRecord struct {
	CountryRecord
	Subdivisions []*Subdivision `json:"subdivisions"`
	City         *City          `json:"city"`
	Location     *Location      `json:"location"`
	PostalCode   string         `json:"postal_code"`
}

type Subdivision struct {
	GeoNameID uint   `json:"geoname_id"`
	IsoCode   string `json:"iso_code"`
	Name      string `json:"name"`
}

type City struct {
	GeoNameID uint   `json:"geoname_id"`
	Name      string `json:"name"`
}

type Location struct {
	AccuracyRadius uint16  `json:"accuracy_radius"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	MetroCode      uint    `json:"metro_code"`
	TimeZone       string  `json:"time_zone"`
}

func FromMaxmindCity(city *geoip2.City, maxmindLang string) *CityRecord {
	ret := CityRecord{
		CountryRecord: CountryRecord{
			Continent: FromMaxmindContinentField(city.Continent, maxmindLang),
			Country:   FromMaxmindCountryField(city.Country, maxmindLang),
		},
		Subdivisions: nil,
		City: &City{
			GeoNameID: city.City.GeoNameID,
			Name:      maxmindGetMLNameFallback(city.City.Names, maxmindLang),
		},
		Location:   (*Location)(&city.Location),
		PostalCode: city.Postal.Code,
	}
	if len(city.Subdivisions) > 0 {
		for _, s := range city.Subdivisions {
			ret.Subdivisions = append(ret.Subdivisions, &Subdivision{
				GeoNameID: s.GeoNameID,
				IsoCode:   s.IsoCode,
				Name:      maxmindGetMLNameFallback(s.Names, maxmindLang),
			})
		}
	}
	return &ret
}
