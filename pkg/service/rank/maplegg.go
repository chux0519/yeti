package rank

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type RankData struct {
	CharacterImageURL  string
	Class              string
	ClassRank          int
	EXP                int64
	EXPPercent         float64
	GlobalRanking      int
	Guild              string
	LegionCoinsPerDay  *int
	LegionLevel        int
	LegionPower        int
	LegionRank         int
	Level              int
	Name               string
	Server             string
	ServerClassRanking int
	ServerRank         int
	ServerSlug         string
	AchievementPoints  int
	AchievementRank    int
	GraphData          []GraphDataItem
}

type GraphDataItem struct {
	AvatarURL        string
	ClassID          int
	ClassRankGroupID int
	CurrentEXP       int64
	DateLabel        string
	EXPDifference    int64
	EXPToNextLevel   int64
	ImportTime       int64
	Level            int
	Name             string
	ServerID         int
	ServerMergeID    int
	TotalOverallEXP  int64
}

/*
GetGMSRank get data from maplestory.gg
*/
func GetGMSRank(ign string) (*RankData, error) {
	content, err := GetGMSRankRaw(ign)
	if err != nil {
		fmt.Printf("Failed to get rank: %s", err.Error())
		return nil, err
	}

	rank := RankData{}
	err = json.Unmarshal(content, &rank)
	if err != nil {
		fmt.Printf("Failed to decode rank: %s", err.Error())
		fmt.Printf("data: %s", string(content))
		return nil, err
	}
	return &rank, nil
}

func GetGMSRankRaw(ign string) ([]byte, error) {
	url := "https://api.maplestory.gg/v1/public/character/gms/" + ign
	content, err := httpGet(url)
	if err != nil {
		fmt.Printf("Failed to get rank: %s", err.Error())
		return nil, err
	}

	return content, nil
}

// get msg with cq code
func (rank *RankData) GetRankReplyString() string {
	var reply string

	reply = fmt.Sprintf("角色：%s \r\n", rank.Name) +
		fmt.Sprintf("服务器：%s \r\n", rank.Server) +
		fmt.Sprintf("等级：%d - %.2f%%  (排名 %d )\r\n", rank.Level, rank.EXPPercent, rank.ServerRank) +
		fmt.Sprintf("职业：%s  (排名 %d )\r\n", rank.Class, rank.ServerClassRanking)

	if rank.LegionCoinsPerDay == nil {
		reply = reply + "非联盟最高角色，无法查询联盟信息"
	} else {
		reply = reply +
			fmt.Sprintf("联盟等级：%d  (排名 %d )\r\n", rank.LegionLevel, rank.LegionRank) +
			fmt.Sprintf("联盟战斗力：%s  (每天 %d 币)", prettyNumber(rank.LegionPower), *rank.LegionCoinsPerDay)
	}

	return reply
}

func prettyNumber(i int) string {
	s := strconv.Itoa(i)
	r1 := ""
	idx := 0

	// Reverse and interleave the separator.
	for i = len(s) - 1; i >= 0; i-- {
		idx++
		if idx == 4 {
			idx = 1
			r1 = r1 + ","
		}
		r1 = r1 + string(s[i])
	}

	// Reverse back and return.
	r2 := ""
	for i = len(r1) - 1; i >= 0; i-- {
		r2 = r2 + string(r1[i])
	}
	return r2
}
