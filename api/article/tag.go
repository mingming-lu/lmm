package article

import (
	"fmt"
	"net/http"

	"github.com/akinaru-lu/elesion"

	"lmm/api/db"
)

type Tag struct {
	ID   int64  `json:"id"`
	User int64  `json:"user"`
	Name string `json:"name"`
}

// GetTags is the handler of GET /users/:user/tags
// get all tags under the given user (id)
func GetTags(c *elesion.Context) {
	queryParams := c.Query()
	if len(queryParams) == 0 {
		c.Status(http.StatusBadRequest).String("empty query parameter")
		return
	}

	values := db.NewValuesFromURL(queryParams)

	tags, err := getTags(values)
	if err != nil {
		c.Status(http.StatusNotFound).Error(err.Error()).String("tags not found")
		return
	}
	c.Status(http.StatusOK).JSON(tags)
}

func getTags(values db.Values) ([]Tag, error) {
	d := db.UseDefault()
	defer d.Close()

	query := fmt.Sprintf(
		`SELECT MIN(id), user, name FROM tag %s GROUP BY name ORDER BY name`,
		values.Where(),
	)

	itr, err := d.Query(query)
	if err != nil {
		return nil, err
	}
	defer itr.Close()

	tags := make([]Tag, 0)
	for itr.Next() {
		tag := Tag{}
		err = itr.Scan(&tag.ID, &tag.User, &tag.Name)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}
	return tags, nil
}

// GetArticleTags get all tags under given article (id)
// GET /users/:user/articles/:article/tags
/*
func GetArticleTags(c *elesion.Context) {
	userIDStr := c.Params.ByName("user")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		c.Status(http.StatusBadRequest).String("invalid user id: " + userIDStr)
		return
	}

	articleIDStr := c.Params.ByName("article")
	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil || articleID <= 0 {
		c.Status(http.StatusBadRequest).String("invalid article id: " + articleIDStr)
		return
	}

	tags, err := getArticleTags(userID, articleID)
	if err != nil {
		c.Status(http.StatusNotFound).Error(err.Error()).String("tags nout found")
		return
	}
	c.Status(http.StatusOK).JSON(tags)
}

func getArticleTags(userID, articleID int64) ([]Tag, error) {
	d := db.UseDefault()
	defer d.Close()

	itr, err := d.Query(
		"SELECT tag FROM article_tag WHERE user = ? AND article = ?",
		userID, articleID,
	)
	if err != nil {
		return make([]Tag, 0), nil
	}
	defer itr.Close()

	ids := make([]int64, 0)
	for itr.Next() {
		var id int64
		err := itr.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return getTagsByIDs(ids...)
}

func getTagsByIDs(ids ...int64) ([]Tag, error) {
	d := db.UseDefault()
	defer d.Close()

	stmt, err := d.Prepare("SELECT id, user, name FROM tag where id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	tags := make([]Tag, 0)
	for _, id := range ids {
		var tag Tag
		err = stmt.QueryRow(id).Scan(&tag.ID, &tag.User, &tag.Name)
		if err != nil {
			return tags, err
		}
		tags = append(tags, tag)
	}
	sort.Sort(byTagName(tags))
	return tags, nil
}

type byTagName []Tag

func (tags byTagName) Len() int {
	return len(tags)
}

func (tags byTagName) Swap(i, j int) {
	tags[i], tags[j] = tags[j], tags[i]
}

func (tags byTagName) Less(i, j int) bool {
	return tags[i].Name < tags[j].Name
}

/*
func NewTags(c *elesion.Context) {
	tags := make([]Tag, 0)
	err := json.NewDecoder(c.Request.Body).Decode(&tags)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid body")
		return
	}
	defer c.Request.Body.Close()

	_, err = newTags(tags)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("failed to add tags")
		return
	}
	c.Status(http.StatusOK).String("success")
}

func newTags(tags []Tag) (int64, error) {
	if tags == nil || len(tags) == 0 {
		return 0, nil
	}
	d := db.New().Use("lmm")
	defer d.Close()

	query := "INSERT INTO tags (user_id, article_id, name) VALUES "
	var values []interface{}
	for _, tag := range tags {
		query += "(?, ?, ?), "
		values = append(values, tag.User, tag.Name)
	}
	query = strings.TrimSuffix(query, ", ")

	stmtIns, err := d.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(values...)
	if err != nil {
		return 0, err
	}
	if rows, err := result.RowsAffected(); err != nil {
		return 0, err
	} else if rows != int64(len(tags)) {
		return 0, errors.Newf("rows affected should be %d, but got %d", len(tags), rows)
	}
	return result.LastInsertId()
}

func DeleteTag(c *elesion.Context) {
	idStr := c.Params.ByName("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invalid id: " + idStr)
		return
	}

	err = deleteTag(id)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("failed to delete tag")
		return
	}
	c.Status(http.StatusOK).String("success")
}

func deleteTag(id int64) error {
	d := db.New().Use("lmm")
	defer d.Close()

	result, err := d.Exec("DELETE FROM tags WHERE id = ?", id)
	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); err != nil {

	} else if rows != 1 {
		return errors.Newf("rows affected should be 1 but got %d", rows)
	}
	return nil
}
*/
