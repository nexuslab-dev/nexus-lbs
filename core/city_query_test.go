package core

import (
	"encoding/json"
	"testing"
)

func TestCityFromCityQuery(t *testing.T) {
	q, err := New("../ipdb/GeoIP2-City.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	c, err := q.City("23.210.75.255", "en")
	if err != nil {
		t.Fatal(err)
	}
	out, _ := json.Marshal(c)
	t.Logf("City=%s", out)
	if c.Country.IsoCode != "JP" {
		t.Fail()
	}

	if c.City.Name != "Tokyo" {
		t.Fail()
	}
}

func TestCountryFromCityQuery(t *testing.T) {
	q, err := New("../ipdb/GeoIP2-City.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	c, err := q.Country("8.8.8.8", "en")
	if err != nil {
		t.Fatal(err)
	}
	out, _ := json.Marshal(c)
	t.Logf("Country=%s", out)
	if c.Country.IsoCode != "US" {
		t.Fail()
	}

	if c.Country.Name != "United States" {
		t.Fail()
	}
}

func TestCityFromCityQueryV6(t *testing.T) {
	q, err := New("../ipdb/dbip-city-lite-2022-06.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	c, err := q.City("2001:bc8:1200:3:208:a2ff:fe0c:5d72", "en")
	if err != nil {
		t.Fatal(err)
	}
	out, _ := json.Marshal(c)
	t.Logf("Country=%s", out)
	if c.Country.IsoCode != "FR" {
		t.Fail()
	}

	if c.City.Name != "Paris" {
		t.Fail()
	}
}
