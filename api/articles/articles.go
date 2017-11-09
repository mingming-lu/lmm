package articles

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	Result     Result `json:"result"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor"`
}

type Result struct {
	Articles []Article `json:"articles"`
}

type Article struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func HandleArticles(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		StatusCode: http.StatusOK,
		Result: Result{
			Articles: []Article{
				{Title: "怒斥香港记者", Text: "吼啊"},
				{Title: "与华莱士谈笑风生", Text: "喝茶"},
				{Title: "视察国机二院", Text: "Excited"},
			},
		},
		HasMore:    false,
		NextCursor: "",
	}
	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), resp.StatusCode)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, string(b))
}
