package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chux0519/yeti/pkg/config"
	logging "github.com/ipfs/go-log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var serverLog = logging.Logger("handlers")

type YetiServer struct {
	e *echo.Echo
	c *config.ServerConfig
	h *YetiHandler
}

func NewYetiServer(config *config.ServerConfig) *YetiServer {
	e := echo.New()

	e.HideBanner = true

	e.Debug = config.Debug

	// ref: https://echo.labstack.com/guide/customization/#logging
	if e.Debug {
		e.Logger.SetLevel(log.DEBUG)
	} else {
		e.Logger.SetLevel(log.INFO)
	}

	// timeout
	e.Server.ReadTimeout = 30 * time.Second
	e.Server.WriteTimeout = 30 * time.Second

	// logs the information about each HTTP request
	// ref: https://echo.labstack.com/middleware/logger/
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper: func(c echo.Context) bool {
			if e.Debug {
				return false
			}
			switch c.Path() {
			case "/ping":
				return true
			default:
				return false
			}
		},
		LogStatus:  true,
		LogURI:     true,
		LogLatency: true,
		LogMethod:  true,
		LogError:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			msg := fmt.Sprintf("method=%s uri=%v status=%v, latency=%v, error=%v", v.Method, v.URI, v.Status, v.Latency.String(), v.Error)
			serverLog.Info(msg)
			return nil
		},
	}))

	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch,
			http.MethodPost, http.MethodDelete, http.MethodOptions,
		},
	}))

	server := YetiServer{}
	server.e = e
	server.c = config
	server.h = NewYetiHandler(config)

	return &server
}

func (k *YetiServer) InitHandlers() {

	k.e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	k.e.POST("/", k.h.EntryHandler)

	k.e.GET("/rank", k.h.LilygoRankHandler)
}

func (k *YetiServer) Serve() {
	k.e.Logger.Fatal(k.e.Start(fmt.Sprintf(":%d", k.c.Port)))
}

func StartYetiServer(config *config.ServerConfig) {
	server := NewYetiServer(config)
	server.InitHandlers()
	server.Serve()
}
