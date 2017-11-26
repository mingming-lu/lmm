package profile

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Name           string           `json:"name"`
	AvatarURL      string           `json:"avatar_url"`
	Bio            string           `json:"bio"`
	Location       string           `json:"location"`
	Profession     string           `json:"profession"`
	Email          string           `json:"email"`
	Skills         []string         `json:"skills"`
	Languages      []Language       `json:"languages"`
	Education      []Education      `json:"education"`
	WorkExperience []WorkExperience `json:"work_experience"`
	Qualifications []Qualification  `json:"qualifications"`
}

type Language struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

type Education struct {
	YearFrom    int    `json:"year_from"`
	MonthFrom   int    `json:"month_from"`
	YearTo      int    `json:"year_to"`
	MonthTo     int    `json:"month_to"`
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Current     bool   `json:"current"`
	Sort        int    `json:"sort"`
}

type WorkExperience struct {
	YearFrom  int    `json:"year_from"`
	MonthFrom int    `json:"month_from"`
	YearTo    int    `json:"year_to"`
	MonthTo   int    `json:"month_to"`
	Company   string `json:"company"`
	Position  string `json:"position"`
	Status    string `json:"status"`
	Current   bool   `json:"current"`
	Sort      int    `json:"sort"`
}

type Qualification struct {
	Year  int    `json:"year"`
	Month int    `json:"month"`
	Name  string `json:"name"`
	Sort  int    `json:"sort"`
}

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Name:       "卢明鸣",
		AvatarURL:  "http://localhost:8082/avatar",
		Bio:        "吾輩が人間である",
		Profession: "Server-side engineer",
		Location:   "鎌倉, 日本",
		Email:      "akinaru.lu@gmail.com",
		Skills: []string{
			"Golang",
			"Perl",
			"Java",
			"Android",
			"Vue.js",
		},
		Languages: []Language{
			{
				Name:  "Chinese",
				Level: "Native",
			},
			{
				Name:  "Japanese",
				Level: "Conversational",
			},
			{
				Name:  "English",
				Level: "Conversational",
			},
		},
		Education: []Education{
			{
				YearFrom:    2010,
				MonthFrom:   9,
				YearTo:      2014,
				MonthTo:     6,
				Institution: "绿头蘑菇学校",
				Degree:      "工学学士",
				Current:     false,
				Sort:        2,
			},
			{
				YearFrom:    2015,
				MonthFrom:   10,
				YearTo:      2017,
				MonthTo:     9,
				Institution: "帝国杜王町大学",
				Degree:      "工学硕士",
				Current:     false,
				Sort:        1,
			},
		},
		WorkExperience: []WorkExperience{
			{
				YearFrom:  2016,
				MonthFrom: 8,
				YearTo:    2016,
				MonthTo:   9,
				Company:   "富士通研究所",
				Position:  "Research",
				Status:    "Internship",
				Current:   false,
				Sort:      3,
			},
			{
				YearFrom:  2016,
				MonthFrom: 11,
				YearTo:    2016,
				MonthTo:   12,
				Company:   "Toshiba",
				Position:  "IoT R&D",
				Status:    "Internship",
				Current:   false,
				Sort:      2,
			},
			{
				YearFrom:  2017,
				MonthFrom: 10,
				Company:   "とある面白法人",
				Position:  "Internship",
				Current:   true,
				Sort:      1,
			},
		},
		Qualifications: []Qualification{
			{
				Year:  2014,
				Month: 1,
				Name:  "JLPT N2",
				Sort:  2,
			},
			{
				Year:  2014,
				Month: 7,
				Name:  "JLPT N1",
				Sort:  1,
			},
		},
	}
	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, string(b))
}
