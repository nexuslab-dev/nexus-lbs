package httpapi

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ttys3/echo-pprof/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"net/http"
)

func NewServer(serviceName string, pprof, tracing bool) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	// e.Logger =
	// e.Validator = NewValidator()

	// middleware
	e.Use(middleware.Recover())

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})

	e.GET("/metrics", func(c echo.Context) error {
		promhttp.Handler().ServeHTTP(c.Response(), c.Request())
		return nil
	})

	// automatically add routers for net/http/pprof
	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	if pprof {
		echopprof.Wrap(e)
	}

	if tracing {
		e.Use(otelecho.Middleware(serviceName, otelecho.WithSkipper(URLSkipper)))
	}
	return e
}
