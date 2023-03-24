package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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

	msg := gmsg.Message
	words := strings.Fields(msg)

	min := 0
	max := 100

	if len(words) == 2 {
		n, err := strconv.Atoi(words[1])
		if err == nil && n >= 0 {
			max = n
		}
	}

	if len(words) == 3 {
		if left, err := strconv.Atoi(words[1]); err == nil {
			if right, err := strconv.Atoi(words[2]); err == nil {
				if left >= 0 {
					min = left
				}
				if right >= 0 {
					max = right
				}
			}
		}
	}

	if min > max {
		min, max = max, min
	}

	hLog.Debug(min, max)
	res := rand.Intn(max-min+1) + min

	resp := fmt.Sprintf("[CQ:reply,id=%d][CQ:at,qq=%d] %d", gmsg.MessageId, gmsg.Sender.UserId, res)

	hLog.Debug(resp)

	cq.SendGroupMessage(gmsg.GroupId, resp)
}
