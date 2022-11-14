package handlers

import (
	"encoding/json"

	"github.com/chux0519/yeti/pkg/service/cqhttp"
)

func HelpHanlder(event map[string]interface{}, cq *cqhttp.YetiCQHTTPService) {
	bytes, err := json.Marshal(event)
	if err != nil {
		hLog.Error(err)
		return
	}

	var gmsg GroupMessage
	err = json.Unmarshal(bytes, &gmsg)
	if err != nil {
		hLog.Error(err)
		return
	}

	reply := "支持的指令: " +
		"/help: 获取帮助\r\n" +
		"/roll: roll 点\r\n" +
		"/mgeta $IGN: 查询角色进度( maple gg eta )\r\n"

	cq.SendGroupMessage(gmsg.GroupId, reply)
}
