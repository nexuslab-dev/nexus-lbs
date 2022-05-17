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
