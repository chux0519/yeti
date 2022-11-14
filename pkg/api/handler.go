package api

import (
	"github.com/chux0519/yeti/pkg/config"
	"github.com/chux0519/yeti/pkg/service/cqhttp"
	"github.com/chux0519/yeti/pkg/service/rank"
	"github.com/chux0519/yeti/pkg/service/sqlite"
)

type YetiHandler struct {
	Config *config.ServerConfig

	CQ *cqhttp.YetiCQHTTPService
	R  *rank.YetiRankService
}

func NewYetiHandler(config *config.ServerConfig) *YetiHandler {
	cq := cqhttp.NewYetiCQHTTPService(config.CQHTTP.Host, config.CQHTTP.AccessToken)
	s, err := sqlite.NewYetiSQLiteService(config.DB.File)
	if err != nil {
		serverLog.Fatal(err)
	}
	r := rank.NewYetiRankService(s, config)

	return &YetiHandler{config, cq, r}
}
