package v1

import "github.com/nexuslab-dev/nexus-lbs/core"

type Request struct {
	IP   string `param:"ip" query:"ip" json:"ip"`
	Lang string `param:"lang" query:"lang" json:"lang"`
}

type Response struct {
	Code    Code        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// CountryResponse for Declarative Comments
type CountryResponse struct {
	Code    Code                `json:"code"`
	Message string              `json:"message"`
	Data    *core.CountryRecord `json:"data,omitempty"`
}

// CityResponse for Declarative Comments
type CityResponse struct {
	Code    Code             `json:"code"`
	Message string           `json:"message"`
	Data    *core.CityRecord `json:"data,omitempty"`
}

// batch api

type RequestBatch struct {
	IP   []string `param:"ip" query:"ip" json:"ip"`
	Lang string   `param:"lang" query:"lang" json:"lang"`
}

// CountryResponseBatch for Declarative Comments
type CountryResponseBatch struct {
	Code    Code                           `json:"code"`
	Message string                         `json:"message"`
	Data    map[string]*core.CountryRecord `json:"data,omitempty"`
}

// CityResponseBatch for Declarative Comments
type CityResponseBatch struct {
	Code    Code                        `json:"code"`
	Message string                      `json:"message"`
	Data    map[string]*core.CityRecord `json:"data,omitempty"`
}
