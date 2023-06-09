package api

import (
	"net/http"

	"github.com/chux0519/yeti/pkg/handlers"
	"github.com/labstack/echo/v4"
)

func (h *YetiHandler) LilygoRankHandler(c echo.Context) error {
	ign := c.QueryParam("ign")

	res, err := handlers.LilygoRankHandler(h.R, ign)
	if err != nil {
		return handlers.NewBadRequestError(err.Error())
	}
	return c.JSON(http.StatusOK, res)
}
