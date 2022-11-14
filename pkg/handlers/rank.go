package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chux0519/yeti/pkg/service/cqhttp"
	"github.com/chux0519/yeti/pkg/service/rank"
)

func MgHanlder(event map[string]interface{}, cq *cqhttp.YetiCQHTTPService, r *rank.YetiRankService) {
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

	msg := gmsg.Message
	words := strings.Fields(msg)
	if len(words) < 2 {
		return
	}

	ign := words[1]
	rank, err := r.FetchUserRank(ign)
	if err != nil {
		hLog.Error(err)
	}

	if rank == nil {
		hLog.Error("failed to get rank, with no error")
		return
	}

	profileFile := r.GetProfileImageName(rank)

	hLog.Debug(profileFile)

	reply := fmt.Sprintf("[CQ:image,file=file://%s]", profileFile)

	hLog.Debug(reply)

	cq.SendGroupMessage(gmsg.GroupId, reply)
}
