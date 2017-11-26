package articles

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Articles   []Article `json:"articles"`
	NextCursor string    `json:"next_cursor"`
}

type Result struct {
}

type Article struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	PostedTime int    `json:"posted_time"`
	EditedTime int    `json:"edited_time"`
}

func HandleArticles(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Articles: []Article{
			{Title: "怒斥香港记者", Text: "Too young too simple, sometimes naive.", PostedTime: 972749166},
			{Title: "与华莱士谈笑风生", Text: "不知道比你们搞到哪里去了", PostedTime: 1142348400},
			{Title: "视察国机二院", Text: "苟利国家生死以，岂因祸福避趋之", PostedTime: 1240488000},
		},
	}
	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, string(b))
}
