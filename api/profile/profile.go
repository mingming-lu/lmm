package profile

import (
	"net/http"
	"strconv"

	"github.com/akinaru-lu/elesion"

	"lmm/api/db"
)

type Profile struct {
	Name           string           `json:"name"`
	AvatarURL      string           `json:"avatar_url"`
	Description    string           `json:"description"`
	Profession     string           `json:"profession"`
	Location       string           `json:"location"`
	Email          string           `json:"email"`
	Skills         []Skill          `json:"skills"`
	Languages      []Language       `json:"languages"`
	Education      []Education      `json:"education"`
	WorkExperience []WorkExperience `json:"work_experience"`
	Qualifications []Qualification  `json:"qualifications"`
}

// GetProfile get profile by given user id
// GET /users/:user/profile
func GetProfile(c *elesion.Context) {
	userIDStr := c.Params.ByName("user")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid user id: " + userIDStr)
		return
	}

	profile, err := getProfile(userID)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error()).String(err.Error())
		return
	}
	c.Status(http.StatusOK).JSON(profile)
}

func getProfile(userID int64) (*Profile, error) {
	d := db.UseDefault()
	defer d.Close()

	profile := Profile{}
	err := d.QueryRow("SELECT name, avatar_url, description, profession, location, email from user where id = 1").Scan(
		&profile.Name, &profile.AvatarURL, &profile.Description, &profile.Profession, &profile.Location, &profile.Email)
	if err != nil {
		return nil, err
	}

	skills, err := skillsByUserID(1)
	if err != nil {
		return nil, err
	}
	profile.Skills = skills

	languages, err := languagesByUserID(1)
	if err != nil {
		return nil, err
	}
	profile.Languages = languages

	workExperience, err := workExperienceByUserID(1)
	if err != nil {
		return nil, err
	}
	profile.WorkExperience = workExperience

	education, err := educationByUserID(1)
	if err != nil {
		return nil, err
	}
	profile.Education = education

	qualifications, err := qualificationByUserID(1)
	if err != nil {
		return nil, err
	}
	profile.Qualifications = qualifications

	return &profile, nil
}

type Skill struct {
	ID   int64  `json:"id"`
	User int64  `json:"user"`
	Name string `json:"name"`
	Sort int64  `json:"sort"`
}

func skillsByUserID(userID int) ([]Skill, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	itr, err := d.Query("SELECT id, user, name, sort FROM skill WHERE user = ? ORDER BY sort", userID)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	skills := make([]Skill, 0)

	for itr.Next() {
		skill := Skill{}
		if e := itr.Scan(&skill.ID, &skill.User, &skill.Name, &skill.Sort); e != nil {
			return skills, err
		}
		skills = append(skills, skill)
	}
	return skills, nil
}

type Language struct {
	ID   int64  `json:"id"`
	User int64  `json:"user"`
	Name string `json:"name"`
	Sort int64  `json:"sort"`
}

func languagesByUserID(id int) ([]Language, error) {
	d := db.New().Use("lmm")
	defer d.Close()

	itr, err := d.Query("SELECT id, user, name, sort FROM language WHERE user = ? ORDER BY sort", id)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	languages := make([]Language, 0)

	for itr.Next() {
		language := Language{}
		if e := itr.Scan(&language.ID, &language.User, &language.Name, &language.Sort); e != nil {
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

	itr, err := d.Query("SELECT date_from, date_to, institution, department, major, degree, current+0 FROM education WHERE user = ? ORDER BY date_from DESC", id)
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

	itr, err := d.Query("SELECT date_from, date_to, company, position, status, current+0 FROM work_experience WHERE user = ? ORDER BY date_from DESC", id)
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

	itr, err := d.Query("SELECT name, date FROM qualification WHERE user = ? ORDER BY date DESC", id)
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
