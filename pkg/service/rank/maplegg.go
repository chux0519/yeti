package rank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/golang/freetype"
	"github.com/wcharczuk/go-chart"
)

var (
	// data from https://docs.google.com/spreadsheets/d/1UjnYrOjibhI4ntPGfrPouvopbEGlKQGGMr_WvZOk4tI/edit#gid=0
	EXP = map[int]int64{
		200: 11462335230,
		205: 25483237508,
		210: 50192858013,
		215: 111178062380,
		220: 226834057694,
		225: 477834581588,
		230: 888805728115,
		235: 1656751310935,
		240: 2780379685705,
		245: 4803825501641,
		250: 7764453421743,
		255: 13095993913257,
		260: 19276710581130,
		265: 34082006515114,
		270: 49642521336419,
		275: 82351036260243,
		280: 164638698169785,
		285: 408002977089330,
		290: 1127748451436850,
		300: 10100775367634700,
	}
)

type RankDataV2 struct {
	CharacterData RankData
}

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

	rank := RankDataV2{}
	err = json.Unmarshal(content, &rank)
	if err != nil {
		fmt.Printf("Failed to decode rank: %s", err.Error())
		fmt.Printf("data: %s", string(content))
		return nil, err
	}
	return &rank.CharacterData, nil
}

func GetGMSRankRaw(ign string) ([]byte, error) {
	url := "https://api.maplestory.gg/v2/public/character/gms/" + ign
	content, err := HttpGet(url)
	if err != nil {
		fmt.Printf("Failed to get rank: %s", err.Error())
		return nil, err
	}

	return content, nil
}

func (rank *RankData) getLastNDaysGraphData(n int) []GraphDataItem {
	now := time.Now()
	now = now.Truncate(24 * time.Hour)

	var ret []GraphDataItem
	for _, data := range rank.GraphData {
		d := data
		ts := time.Unix(data.ImportTime, 0)
		diffHours := now.Sub(ts).Hours()
		diffTimeDays := int(diffHours / 24)
		if diffTimeDays < n && diffHours >= 0 {
			ret = append(ret, d)
		}
	}

	return ret
}

func (rank *RankData) GetExpChart() (*chart.BarChart, error) {
	bars := []chart.Value{}

	for _, data := range rank.getLastNDaysGraphData(7) {
		barStyle := chart.Style{
			FillColor:   chart.ColorBlue,
			StrokeColor: chart.ColorBlue,
			StrokeWidth: 0,
		}
		exp := float64(data.EXPDifference) / 1000000000.0
		dateTime := time.Unix(data.ImportTime, 0)

		v := chart.Value{
			Value: exp,
			Label: fmt.Sprintf("%s\n(%s)", dateTime.Format("2006-01-02"), prettyNumber(float64(data.EXPDifference), "")),
			Style: barStyle,
		}
		bars = append(bars, v)
	}

	graph := chart.BarChart{
		Background: chart.Style{
			Padding: chart.Box{
				Top:    40,
				Bottom: 40,
				Right:  40,
			},
		},
		Height:   640,
		BarWidth: 40,
		Bars:     bars,
	}

	return &graph, nil
}

func (rank *RankData) GetProfileImage() ([]byte, error) {
	// exp chart
	expChart, err := rank.GetExpChart()
	if err != nil {
		rLog.Error(err)
		return nil, err
	}

	infoWidth := 400
	w := expChart.GetWidth() + infoWidth
	h := expChart.GetHeight()

	fontCtx := freetype.NewContext()
	fontCtx.SetDPI(92)
	font, err := chart.GetDefaultFont()
	if err != nil {
		rLog.Error(err)
		return nil, err
	}
	fontCtx.SetFont(font)

	// bg
	profileImg := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(profileImg, profileImg.Bounds(), &image.Uniform{chart.ColorLightGray}, image.Point{}, draw.Src)

	// add avatar
	avatarTopPadding := 50
	if rank.CharacterImageURL == "" {
		return nil, fmt.Errorf("invalid avatar")
	}
	avatarBytes, _ := HttpGet(rank.CharacterImageURL)

	avartarImg, err := png.Decode(bytes.NewBuffer(avatarBytes))
	if err != nil {
		rLog.Error(err)
		return nil, err
	}
	avatarOffset := image.Pt((infoWidth-avartarImg.Bounds().Dx())/2, avatarTopPadding)
	draw.Draw(profileImg, avartarImg.Bounds().Add(avatarOffset), avartarImg, image.Point{}, draw.Over)

	// avatar info
	textPadding := 50

	addLabel := func(label string) error {
		fontCtx.SetClip(image.Rect(0, 0, infoWidth, h))
		fontCtx.SetDst(profileImg)
		fontSize := 12.0
		fontCtx.SetFontSize(fontSize)
		fontCtx.SetSrc(&image.Uniform{chart.ColorBlack})
		pt := freetype.Pt(70, avatarOffset.Y+avartarImg.Bounds().Dy()+textPadding+int(fontCtx.PointToFixed(fontSize)>>6))
		if _, err := fontCtx.DrawString(label, pt); err != nil {
			return err
		}
		textPadding += 35
		return nil
	}

	replyStr := rank.GetRankReplyString() + "\r\n" + rank.GetEtaString()
	labels := strings.Split(replyStr, "\r\n")
	for _, label := range labels {
		if err := addLabel(label); err != nil {
			rLog.Error(err)
			return nil, err
		}
	}

	// exp chart
	expChartBuffer := bytes.NewBuffer([]byte{})
	if err = expChart.Render(chart.PNG, expChartBuffer); err != nil {
		rLog.Error(err)
		return nil, err
	}

	expImg, err := png.Decode(expChartBuffer)
	if err != nil {
		rLog.Error(err)
		return nil, err
	}
	draw.Draw(profileImg, expImg.Bounds().Add(image.Pt(infoWidth, 0)), expImg, image.Point{}, draw.Over)

	buffer := bytes.NewBuffer([]byte{})

	err = png.Encode(buffer, profileImg)
	if err != nil {
		rLog.Error(err)
		return nil, err
	}

	return ioutil.ReadAll(buffer)
}

// get msg with cq code
func (rank *RankData) GetRankReplyString() string {
	var reply string

	reply = fmt.Sprintf("Name: %s \r\n", rank.Name) +
		fmt.Sprintf("Server: %s \r\n", rank.Server) +
		fmt.Sprintf("Level: %d - %.2f%%  (Rank %d )\r\n", rank.Level, rank.EXPPercent, rank.ServerRank) +
		fmt.Sprintf("Job: %s  (Rank %d )\r\n", rank.Class, rank.ServerClassRanking)

	if rank.LegionCoinsPerDay == nil {
		reply = reply + "Legion\r\n" + "-"
	} else {
		reply = reply +
			"Legion\r\n" +
			fmt.Sprintf("Level: %d  (Rank %d )\r\n", rank.LegionLevel, rank.LegionRank) +
			fmt.Sprintf("Power: %s  (%d Coins/Day)", prettyNumberInt(rank.LegionPower), *rank.LegionCoinsPerDay)
	}

	return reply
}

func (rank *RankData) getCurrentEXP() int64 {
	var currentExp int64 = 0
	for _, data := range rank.GraphData {
		if data.TotalOverallEXP > currentExp {
			currentExp = data.TotalOverallEXP
		}
	}

	return currentExp
}

func (rank *RankData) GetEtaString() string {
	var reply string

	// 1, 3, 7
	days := []int{3}
	currentExp := rank.getCurrentEXP()
	for _, day := range days {
		diff := rank.getAvgDiff(day)

		nextGapLvl := getNextGapLevel(rank.Level)
		nextGapLvlExp := getNextGapLevelTotalExp(rank.Level)

		eta := float64(nextGapLvlExp-currentExp) / diff

		reply += "ETA to Level\r\n"
		reply += fmt.Sprintf("Lv. %d ( %.2f days)\r\n", nextGapLvl, eta)
		// reply += fmt.Sprintf("按最近 %d 天的肝度估算，需要 %.2f 天到 %d 级", day, eta, nextGapLvl)
		if nextGapLvl < 250 {
			eta := float64(EXP[250]-currentExp) / diff
			reply += fmt.Sprintf("Lv. %d ( %.2f days)\r\n", 250, eta)
		} else if nextGapLvl < 275 {
			eta := float64(EXP[275]-currentExp) / diff
			reply += fmt.Sprintf("Lv. %d ( %.2f days)\r\n", 275, eta)
		} else if nextGapLvl < 280 {
			eta := float64(EXP[280]-currentExp) / diff
			reply += fmt.Sprintf("Lv. %d ( %.2f days)\r\n", 280, eta)
		} else if nextGapLvl < 300 {
			eta := float64(EXP[300]-currentExp) / diff
			reply += fmt.Sprintf("Lv. %d ( %.2f days)\r\n", 300, eta)
		}
		dailyExp := prettyNumber(diff, "")
		reply += fmt.Sprintf("Avg Daily Exp On %d Day(s): %s\r\n", day, dailyExp)
	}

	return reply
}

func getNextGapLevel(lvl int) int {
	ret := (int(lvl/5) + 1) * 5
	if ret < 200 {
		return 200
	}
	return ret
}

func getNextGapLevelTotalExp(lvl int) int64 {
	i := getNextGapLevel(lvl)
	return EXP[i]
}

func (rank *RankData) getAvgDiff(day int) float64 {
	if day < 1 || day > 15 {
		return 0
	}
	var diff int64 = 0
	for _, data := range rank.getLastNDaysGraphData(day) {
		diff += data.EXPDifference
	}
	return float64(diff) / float64(day)
}

func prettyNumberInt(i int) string {
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

func prettyNumber(diff float64, suffix string) string {
	var remain float64
	var base float64
	if diff < 1000 {
		remain = diff - float64(int(diff))
	} else if diff < 1000000 {
		// k
		if suffix == "" {
			suffix = "k"
		}
		base = diff / 1000
		remain = base - float64(int(base))
	} else if diff < 1000000000 {
		// m
		if suffix == "" {
			suffix = "m"
		}
		base = diff / 1000000
		remain = base - float64(int(base))
	} else if diff < 1000000000000 {
		// b
		if suffix == "" {
			suffix = "b"
		}
		base = diff / 1000000000
		remain = base - float64(int(base))
	} else {
		if suffix == "" {
			suffix = "t"
		}
		base = diff / 1000000000000
		remain = base - float64(int(base))
	}
	exp := prettyNumberInt(int(base))

	remain = math.Floor(remain*100) / 100
	remainStr := fmt.Sprintf("%.2f", remain)
	remainStrArr := strings.Split(remainStr, ".")
	point := remainStrArr[len(remainStrArr)-1]

	return fmt.Sprintf("%s.%s %s", exp, point, suffix)
}
