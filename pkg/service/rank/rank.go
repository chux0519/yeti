package rank

import (
	"encoding/json"
	"strings"

	"github.com/chux0519/yeti/pkg/service/sqlite"

	logging "github.com/ipfs/go-log"
)

var rLog = logging.Logger("rank")

type YetiRankService struct {
	S *sqlite.YetiSQLiteService
}

func NewYetiRankService(s *sqlite.YetiSQLiteService) *YetiRankService {
	return &YetiRankService{s}
}

func (r *YetiRankService) FetchUserRank(ign string) (*RankData, error) {
	record, err := r.S.GetMapleGGByIGN(ign)
	if err != nil {
		return nil, err
	}

	if record != nil && strings.EqualFold(record.IGN, ign) {
		data := record.Data
		dBytes, err := json.Marshal(&data)
		if err != nil {
			return nil, err
		}
		var rank RankData
		if err := json.Unmarshal(dBytes, &rank); err != nil {
			return nil, err
		}
		return &rank, nil
	} else {
		rank, err := GetGMSRank(ign)
		if err != nil {
			return nil, err
		}

		dBytes, err := json.Marshal(rank)
		if err != nil {
			return nil, err
		}

		var data map[string]interface{}
		if err := json.Unmarshal(dBytes, &data); err != nil {
			return nil, err
		}

		// TODO: fetch the avatar bytes
		// upsert
		_, err = r.S.UpsertMapleGG(ign, data, []byte("test"))
		if err != nil {
			return nil, err
		}

		return rank, nil
	}
}
