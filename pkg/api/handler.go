package api

import (
	"github.com/chux0519/yeti/pkg/config"
	"github.com/chux0519/yeti/pkg/service/cqhttp"
)

type YetiHandler struct {
	Config *config.ServerConfig

	CQ *cqhttp.YetiCQHTTPService
}

func NewYetiHandler(config *config.ServerConfig) *YetiHandler {
	cq := cqhttp.NewYetiCQHTTPService(config.CQHTTP.Host, config.CQHTTP.AccessToken)
	return &YetiHandler{config, cq}
}
