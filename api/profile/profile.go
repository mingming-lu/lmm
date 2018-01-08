package profile

import (
	"net/http"

	"lmm/api/db"

	"github.com/akinaru-lu/elesion"
)

type Profile struct {
	Name           string           `json:"name"`
	AvatarURL      string           `json:"avatar_url"`
	Bio            string           `json:"bio"`
	Location       string           `json:"location"`
	Profession     string           `json:"profession"`
	Email          string           `json:"email"`
	Skills         []Skill          `json:"skills"`
	Languages      []Language       `json:"languages"`
	Education      []Education      `json:"education"`
	WorkExperience []WorkExperience `json:"work_experience"`
	Qualifications []Qualification  `json:"qualifications"`
}

type Skill struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Sort   int64  `json:"sort"`
}

func skillsByUserID(userID int) ([]Skill, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	itr, err := d.Query("SELECT id, user_id, name, sort FROM skill WHERE user_id = ? ORDER BY sort", userID)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	skills := make([]Skill, 0)

	for itr.Next() {
		skill := Skill{}
		if e := itr.Scan(&skill.ID, &skill.UserID, &skill.Name, &skill.Sort); e != nil {
			return skills, err
		}
		skills = append(skills, skill)
	}
	return skills, nil
}

type Language struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Sort   int64  `json:"sort"`
}

func languagesByUserID(id int) ([]Language, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	itr, err := d.Query("SELECT id, user_id, name, sort FROM language WHERE user_id = ? ORDER BY sort", id)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	languages := make([]Language, 0)

	for itr.Next() {
		language := Language{}
		if e := itr.Scan(&language.ID, &language.UserID, &language.Name, &language.Sort); e != nil {
			return languages, err
		}
		languages = append(languages, language)
	}
	return languages, nil
}

type Education struct {
	DateFrom    string `json:"date_from"`
	DateTo      string `json:"date_to"`
	Institution string `json:"institution"`
	Department  string `json:"department"`
	Major       string `json:"major"`
	Degree      string `json:"degree"`
	Current     bool   `json:"current"`
}

func educationByUserID(id int) ([]Education, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	itr, err := d.Query("SELECT date_from, date_to, institution, department, major, degree, current+0 FROM education WHERE user_id = ? ORDER BY date_from DESC", id)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	education := make([]Education, 0)
	var current int
	for itr.Next() {
		edu := Education{}
		err = itr.Scan(&edu.DateFrom, &edu.DateTo, &edu.Institution, &edu.Department, &edu.Major, &edu.Degree, &current)
		if err != nil {
			return education, err
		}
		if current == 1 {
			edu.Current = true
		}
		education = append(education, edu)
	}
	return education, nil
}

type WorkExperience struct {
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
	Company  string `json:"company"`
	Position string `json:"position"`
	Status   string `json:"status"`
	Current  bool   `json:"current"`
}

func workExperienceByUserID(id int) ([]WorkExperience, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	itr, err := d.Query("SELECT date_from, date_to, company, position, status, current+0 FROM work_experience WHERE user_id = ? ORDER BY date_from DESC", id)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	workExperience := make([]WorkExperience, 0)
	var current int
	for itr.Next() {
		we := WorkExperience{}
		if e := itr.Scan(&we.DateFrom, &we.DateTo, &we.Company, &we.Position, &we.Status, &current); e != nil {
			return workExperience, err
		}
		if current == 1 {
			we.Current = true
		}
		workExperience = append(workExperience, we)
	}
	return workExperience, nil
}

type Qualification struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

func qualificationByUserID(id int) ([]Qualification, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	itr, err := d.Query("SELECT name, date FROM qualification WHERE user_id = ? ORDER BY date DESC", id)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	qualification := make([]Qualification, 0)
	for itr.Next() {
		q := Qualification{}
		err := itr.Scan(&q.Name, &q.Date)
		if err != nil {
			return qualification, err
		}
		qualification = append(qualification, q)
	}
	return qualification, nil
}

func getProfile(c *elesion.Context) {
	d := db.New().Use("lmm")
	defer d.Close()

	profile := Profile{}
	err := d.QueryRow("SELECT name, avatar_url, bio, location, profession, email from profile where id = 1").Scan(
		&profile.Name, &profile.AvatarURL, &profile.Bio, &profile.Location, &profile.Profession, &profile.Email)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}

	skills, err := skillsByUserID(1)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}
	profile.Skills = skills

	languages, err := languagesByUserID(1)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}
	profile.Languages = languages

	workExperience, err := workExperienceByUserID(1)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}
	profile.WorkExperience = workExperience

	education, err := educationByUserID(1)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}
	profile.Education = education

	qualifications, err := qualificationByUserID(1)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error())
		return
	}
	profile.Qualifications = qualifications

	c.Status(200).JSON(profile)
}

func Handler(c *elesion.Context) {
	if c.Request.Method == http.MethodGet {
		getProfile(c)
	} else {
		c.Status(http.StatusMethodNotAllowed)
	}
}
