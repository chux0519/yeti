package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/chux0519/yeti/pkg/service/cqhttp"
)

type Sender struct {
	Nickname string `json:"nickname"`
	UserId   int64  `json:"user_id"`
}

type GroupMessage struct {
	GroupId   int64   `json:"group_id"`
	MessageId int64   `json:"message_id"`
	Message   string  `json:"message"`
	Sender    *Sender `json:"sender"`
}

func RollHanlder(event map[string]interface{}, cq *cqhttp.YetiCQHTTPService) {
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

	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 100
	res := rand.Intn(max-min+1) + min

	msg := fmt.Sprintf("[CQ:reply,id=%d][CQ:at,qq=%d] %d", gmsg.MessageId, gmsg.Sender.UserId, res)

	hLog.Debug(msg)

	cq.SendGroupMessage(gmsg.GroupId, msg)
}
