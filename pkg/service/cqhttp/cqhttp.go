package cqhttp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chux0519/yeti/pkg/service/utils"
	logging "github.com/ipfs/go-log"
)

var cqLog = logging.Logger("handlers")

type YetiCQHTTPService struct {
	BaseUrl     string
	AccessToken string
}

func NewYetiCQHTTPService(url string, token string) *YetiCQHTTPService {
	return &YetiCQHTTPService{url, token}
}

type GroupMessageBody struct {
	GroupId   int64  `json:"group_id"`
	Message   string `json:"message"`
	AsRawText bool   `json:"auto_escape"`
}

func (cq *YetiCQHTTPService) SendGroupMessage(groupId int64, msg string) error {
	url := fmt.Sprintf("%s/send_group_msg?access_token=%s", cq.BaseUrl, cq.AccessToken)

	body := GroupMessageBody{
		GroupId:   groupId,
		Message:   msg,
		AsRawText: false,
	}

	data, err := json.Marshal(&body)
	if err != nil {
		cqLog.Error(err)
		return err
	}

	res, err := utils.HttpPost(url, data, 10*time.Second)

	if err != nil {
		cqLog.Error(err)
		return err
	}
	cqLog.Debug(string(res))
	return nil
}
