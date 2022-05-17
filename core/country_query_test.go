package core

import (
	"encoding/json"
	"testing"
)

func TestCountryQuery(t *testing.T) {
	q, err := New("../ipdb/dbip-country.mmdb")
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
