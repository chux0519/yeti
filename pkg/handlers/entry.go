package handlers

import (
	"strings"

	"github.com/chux0519/yeti/pkg/service/cqhttp"
	"github.com/chux0519/yeti/pkg/service/rank"
)

const (
	PREFIX_ROLL = "/roll"
	PREFIX_MG   = "/mg"
	PREFIX_HELP = "/help"
)

func EntryHandler(event map[string]interface{}, cq *cqhttp.YetiCQHTTPService, r *rank.YetiRankService) {

	imsg := event["message"]
	msg, ok := imsg.(string)

	if ok {
		hLog.Debug(event)
		lowerMsg := strings.ToLower(msg)
		if strings.HasPrefix(lowerMsg, PREFIX_HELP) {
			HelpHanlder(event, cq)
		}
		if strings.HasPrefix(lowerMsg, PREFIX_ROLL) {
			RollHanlder(event, cq)
		}
		if strings.HasPrefix(lowerMsg, PREFIX_MG) {
			MgHanlder(event, cq, r)
		}
	}
}
