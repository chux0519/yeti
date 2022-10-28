package api

import (
	"net/http"

	"github.com/chux0519/yeti/pkg/handlers"
	"github.com/labstack/echo/v4"
)

func (h *YetiHandler) EntryHandler(c echo.Context) error {
	event := make(map[string]interface{})
	if err := c.Bind(&event); err != nil {
		return handlers.NewBadRequestError(err.Error())
	}

	go handlers.EntryHandler(event, h.CQ)

	return c.JSON(http.StatusOK, "")
}
