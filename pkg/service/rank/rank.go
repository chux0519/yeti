package rank

import (
	"encoding/json"
	"fmt"
)

type RankData struct {
	CharacterImageURL  string
	Class              string
	ClassRank          int
	EXP                int64
	EXPPercent         float64
	GlobalRanking      int
	Guild              string
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
	url := "https://api.maplestory.gg/v1/public/character/gms/" + ign
	content, err := httpGet(url)
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
