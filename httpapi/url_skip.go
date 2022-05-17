package httpapi

import (
	"github.com/labstack/echo/v4"
	"strings"
)

func URLSkipper(c echo.Context) bool {
	if strings.HasPrefix(c.Path(), "/healthz") {
		return true
	}
	if strings.HasPrefix(c.Path(), "/metrics") {
		return true
	}
	if strings.HasPrefix(c.Path(), "/debug/pprof") {
		return true
	}
	if strings.HasPrefix(c.Path(), "/pprof") {
		return true
	}
	if strings.HasPrefix(c.Path(), "/config") {
		return true
	}
	return false
}
