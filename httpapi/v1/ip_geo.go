package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/nexuslab-dev/nexus-lbs/core"
	"net/http"
)

type IPGeo struct {
	country *core.GeoQuery
	city    *core.GeoQuery
}

func New(countryQuery *core.GeoQuery, cityQuery *core.GeoQuery) *IPGeo {
	return &IPGeo{
		country: countryQuery,
		city:    cityQuery,
	}
}

// CountryHandler godoc
// @Summary      query country by IP
// @Description  query country by IP, with optional lang param
// @Tags         lbs
// @Produce      json
// @Param        ip   path      string  true  "IP address"
// @Param        lang query      string  false  "response language"
// @Success      200  {object}  CountryResponse
// @Router       /country/{ip} [get]
func (a *IPGeo) CountryHandler(c echo.Context) error {
	req := &Request{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusOK, &Response{
			Code:    ErrParam,
			Message: "invalid param",
		})
	}

	country, err := a.country.Country(req.IP, req.Lang)
	if err != nil {
		return c.JSON(http.StatusOK, &Response{
			Code:    ErrGeoip,
			Message: "query country from db failed",
		})
	}

	return c.JSON(http.StatusOK, &Response{
		Code:    OK,
		Message: "success",
		Data:    country,
	})
}

// CountryBatchHandler godoc
// @Summary      batch query country by IPs
// @Description  batch query country by IPs, with optional lang param
// @Tags         lbs
// @accept       json
// @Produce      json
// @param        payload   body RequestBatch true "json payload"
// @Success      200  {object}  CountryResponseBatch
// @Router       /country/batch [post]
func (a *IPGeo) CountryBatchHandler(c echo.Context) error {
	req := &RequestBatch{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusOK, &Response{
			Code:    ErrParam,
			Message: "invalid param",
		})
	}

	rs := make(map[string]*core.CountryRecord)
	for _, ip := range req.IP {
		country, err := a.country.Country(ip, req.Lang)
		if err != nil {
			return c.JSON(http.StatusOK, &Response{
				Code:    ErrGeoip,
				Message: "query country from db failed, failed IP=" + ip,
			})
		}
		rs[ip] = country
	}

	return c.JSON(http.StatusOK, &Response{
		Code:    OK,
		Message: "success",
		Data:    rs,
	})
}

// CityHandler godoc
// @Summary      query city by IP
// @Description  query city by IP, with optional lang param
// @Tags         lbs
// @Produce      json
// @Param        ip   path      string  true  "IP address"
// @Param        lang query     string  false  "response language"
// @Success      200  {object}  CityResponse
// @Router       /city/{ip} [get]
func (a *IPGeo) CityHandler(c echo.Context) error {
	req := &Request{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusOK, &Response{
			Code:    ErrParam,
			Message: "invalid param",
		})
	}

	city, err := a.city.City(req.IP, req.Lang)
	if err != nil {
		return c.JSON(http.StatusOK, &Response{
			Code:    ErrGeoip,
			Message: "query city from db failed",
		})
	}

	return c.JSON(http.StatusOK, &Response{
		Code:    OK,
		Message: "success",
		Data:    city,
	})
}

// CityBatchHandler godoc
// @Summary      batch query city by IPs
// @Description  batch query city by IPs, with optional lang param
// @Tags         lbs
// @accept       json
// @Produce      json
// @param        payload   body RequestBatch true "json payload"
// @Success      200  {object}  CityResponseBatch
// @Router       /city/batch [post]
func (a *IPGeo) CityBatchHandler(c echo.Context) error {
	req := &RequestBatch{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusOK, &Response{
			Code:    ErrParam,
			Message: "invalid param",
		})
	}

	rs := make(map[string]*core.CityRecord)
	for _, ip := range req.IP {
		city, err := a.city.City(ip, req.Lang)
		if err != nil {
			return c.JSON(http.StatusOK, &Response{
				Code:    ErrGeoip,
				Message: "query city from db failed, failed IP=" + ip,
			})
		}
		rs[ip] = city
	}

	return c.JSON(http.StatusOK, &Response{
		Code:    OK,
		Message: "success",
		Data:    rs,
	})
}
