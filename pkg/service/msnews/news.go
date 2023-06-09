package msnews

import (
	"encoding/json"
	"time"

	"github.com/chux0519/yeti/pkg/service/rank"
)

type NewsResp struct {
	Title       string      `json:"title"`
	HomePageUrl string      `json:"home_page_url"`
	Items       *[]NewsItem `json:"items"`
}

// for example
//
//	{
//		"id": "https://maplestory.nexon.net/news/80155/completed-unscheduled-channel-maintenance-january-27-2023",
//		"url": "https://maplestory.nexon.net/news/80155/completed-unscheduled-channel-maintenance-january-27-2023",
//		"title": "[Completed] Unscheduled Channel Maintenance - January 27, 2023",
//		"content_html": "The maintenance has been completed and all channels are now available. Thank you for your patience.",
//		"date_published": "2023-01-28T03:13:44.000Z"
//	}

type NewsItem struct {
	Url           string `json:"url"`
	Title         string `json:"title"`
	Content       string `json:"content_html"`
	DatePublushed string `json:"date_published"`
}

// TODO: check in 10 mins(rss has 5 mins cache by default)
// use a redis to cache if we push the news or not
// push to configured groups
func CheckMapleNews(url string) (*[]NewsItem, error) {
	data, err := rank.HttpGet(url)
	if err != nil {
		return nil, err
	}

	var resp NewsResp

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	var ret []NewsItem
	for _, item := range *resp.Items {
		pubTime, err := time.Parse(time.RFC3339, item.DatePublushed)
		if err == nil {
			now := time.Now()
			diffHours := now.Sub(pubTime).Hours()
			if diffHours < 8 {
				ret = append(ret, item)
			}
		}
	}

	return &ret, nil
}
