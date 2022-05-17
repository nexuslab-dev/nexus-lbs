package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

const serviceNameHeader = "Nexus-Service-Name"

type Options struct {
	serviceName    string
	serviceVersion string
}

func NewServerNameHeader(serviceName, serviceVersion string) echo.MiddlewareFunc {
	return (&Options{
		serviceName:    serviceName,
		serviceVersion: serviceVersion,
	}).handle
}

// ServerNameHeader middleware adds a `Server` header to the response.
func (s *Options) handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(serviceNameHeader, fmt.Sprintf("%v/%v", s.serviceName, s.serviceVersion))
		return next(c)
	}
}
