package profile

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	Result     Result `json:"result"`
}

type Result struct {
	Name string `json:"name"`
}

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		StatusCode: http.StatusOK,
		Result: Result{
			Name: "卢明鸣",
		},
	}
	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), resp.StatusCode)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, string(b))
}
