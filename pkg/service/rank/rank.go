package rank

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chux0519/yeti/pkg/config"
	"github.com/chux0519/yeti/pkg/service/sqlite"

	logging "github.com/ipfs/go-log"
)

var rLog = logging.Logger("rank")

type YetiRankService struct {
	S      *sqlite.YetiSQLiteService
	Config *config.ServerConfig
}

func NewYetiRankService(s *sqlite.YetiSQLiteService, config *config.ServerConfig) *YetiRankService {
	return &YetiRankService{s, config}
}

func (r *YetiRankService) FetchUserRank(ign string) (*RankData, error) {
	record, err := r.S.GetMapleGGByIGN(ign)

	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return nil, err
		}
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
		profileImgBytes, err := rank.GetProfileImage()
		if err != nil {
			rLog.Error(err)
			return nil, err
		}

		if err := r.saveProfileImageToCache(r.GetProfileImageName(&rank), profileImgBytes); err != nil {
			rLog.Error(err)
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
		// upsert
		profileImgBytes, err := rank.GetProfileImage()
		if err != nil {
			rLog.Error(err)
			return nil, err
		}

		if err := r.saveProfileImageToCache(r.GetProfileImageName(rank), profileImgBytes); err != nil {
			rLog.Error(err)
			return nil, err
		}

		_, err = r.S.UpsertMapleGG(ign, data, profileImgBytes)
		if err != nil {
			return nil, err
		}

		return rank, nil
	}
}

func (r *YetiRankService) GetProfileImageName(rank *RankData) string {
	return filepath.Join(r.Config.CQHTTP.CacheDir, fmt.Sprintf("user_profile_%s.png", rank.Name))
}

func (r *YetiRankService) saveProfileImageToCache(fileName string, bytes []byte) error {
	out, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	out.Write(bytes)
	return nil
}
