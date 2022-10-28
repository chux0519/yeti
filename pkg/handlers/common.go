package handlers

import (
	"net/http"

	logging "github.com/ipfs/go-log"
	"github.com/labstack/echo/v4"
)

var hLog = logging.Logger("handlers")

type OkResult struct {
	Ok bool `json:"ok"`
}

func NewBadRequestError(message ...interface{}) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, message...)
}

func NewUnauthorizedError(message ...interface{}) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusUnauthorized, message...)
}

func NewForbiddenError(message ...interface{}) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusForbidden, message...)
}

func NewNotFoundError(message ...interface{}) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusNotFound, message...)
}

func NewInternalServerError(message ...interface{}) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusInternalServerError, message...)
}
