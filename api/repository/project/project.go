package project

import (
	"lmm/api/db"

	model "lmm/api/domain/model/project"

	"github.com/akinaru-lu/errors"
)

func Add(project *model.Project) (int64, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("INSERT INTO project (user, name, icon, url, description, from_date, to_date")
	defer stmt.Close()

	res, err := stmt.Exec(project.User, project.Name, project.Icon, project.URL, project.Description, *project.FromDate, *project.ToDate)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func FetchByID(id int64) (*model.Project, error) {
	projects, err := fetchBySingle("id", id)
	if projects != nil && len(projects) == 1 {
		return &projects[0], err
	}
	return nil, errors.Wrap(err, "Failed to fetch project")
}

func FetchByUser(userID int64) ([]model.Project, error) {
	return fetchBySingle("user", userID)
}

func fetchBySingle(field string, value interface{}) ([]model.Project, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Mustf("SELECT id, user, name, icon, url, description, from_date, to_date FROM project WHERE %s = ? ORDER BY from_date DESC", field)
	defer stmt.Close()

	projects := make([]model.Project, 0)
	itr, err := stmt.Query(value)
	if err != nil {
		return projects, err
	}

	for itr.Next() {
		project := model.Project{}
		err = itr.Scan(&project.ID, &project.User, &project.Name, &project.Icon, &project.URL, &project.Description, &project.FromDate, &project.ToDate)
		if err != nil {
			return projects, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func Update() error {
	return errors.New("Not implemented")
}

func Delete() error {
	return errors.New("Not implemented")
}
