package handlers

import (
	"strings"

	"github.com/chux0519/yeti/pkg/service/cqhttp"
	"github.com/chux0519/yeti/pkg/service/rank"
)

const (
	PREFIX_ROLL  = "/roll"
	PREFIX_MGETA = "/mgeta"
)

func EntryHandler(event map[string]interface{}, cq *cqhttp.YetiCQHTTPService, r *rank.YetiRankService) {

	imsg := event["message"]
	msg, ok := imsg.(string)

	if ok {
		hLog.Debug(event)
		if strings.HasPrefix(msg, PREFIX_ROLL) {
			RollHanlder(event, cq)
		}
		if strings.HasPrefix(msg, PREFIX_MGETA) {
			MgetaHanlder(event, cq, r)
		}
	}
}
