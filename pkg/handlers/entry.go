package handlers

import (
	"strings"

	"github.com/chux0519/yeti/pkg/service/cqhttp"
)

const (
	PREFIX_ROLL = "/roll"
)

func EntryHandler(event map[string]interface{}, cq *cqhttp.YetiCQHTTPService) {

	imsg := event["message"]
	msg, ok := imsg.(string)

	if ok {
		hLog.Debug(event)
		if strings.HasPrefix(msg, PREFIX_ROLL) {
			RollHanlder(event, cq)
		}
	}
}
