package blog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/akinaru-lu/elesion"
	"github.com/akinaru-lu/errors"

	"lmm/api/db"
	"lmm/api/user"
)

type Blog struct {
	ID        int64  `json:"id"`
	User      int64  `json:"user"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetBlog gets all blog according to query parameters
// GET /blog
func GetBlog(c *elesion.Context) {
	queryParams := c.Query()
	if len(queryParams) == 0 {
		c.Status(http.StatusBadRequest).String("empty query parameter")
		return
	}

	values := db.NewValuesFromURL(c.Query())
	blog, err := getBlog(values)
	if err != nil {
		c.Status(http.StatusNotFound).Error(err.Error()).String("blog not found")
		return
	}
	c.Status(http.StatusOK).JSON(blog)
}

func getBlog(values db.Values) ([]Blog, error) {
	d := db.UseDefault()
	defer d.Close()

	query := fmt.Sprintf(
		`SELECT id, user, title, text, created_at, updated_at FROM blog %s ORDER BY created_at DESC`,
		values.Where(),
	)

	blogList := make([]Blog, 0)
	cursor, err := d.Query(query)
	if err != nil {
		return blogList, err
	}
	defer cursor.Close()

	for cursor.Next() {
		blog := Blog{}
		err = cursor.Scan(&blog.ID, &blog.User, &blog.Title, &blog.Text, &blog.CreatedAt, &blog.UpdatedAt)
		if err != nil {
			return blogList, err // return all blogList found with error
		}
		blogList = append(blogList, blog)
	}
	return blogList, nil
}

// NewBlog post new blog to the user given by url path
// POST /blog
func NewBlog(c *elesion.Context) {
	usr, err := user.CheckAuth(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String(err.Error())
		return
	}

	body := Blog{}
	err = json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("invalid body")
		return
	}
	defer c.Request.Body.Close()

	body.User = usr.ID

	id, err := newBlog(body)
	if err != nil {
		c.Status(http.StatusInternalServerError).Error(err.Error()).String("failed to post blog")
		return
	}

	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("success to post blog but failed when post tags")
		return
	}
	c.Header("location", fmt.Sprintf("blog/%d", id)).Status(http.StatusCreated).String("success")
}

func newBlog(blog Blog) (int64, error) {
	d := db.UseDefault()
	defer d.Close()

	result, err := d.Exec(
		"INSERT INTO blog (user, title, text) VALUES (?, ?, ?)",
		blog.User, strings.TrimSpace(blog.Title), strings.TrimSpace(blog.Text),
	)
	if err != nil {
		return 0, err
	}
	if rows, err := result.RowsAffected(); err != nil {
		return 0, err
	} else if rows != 1 {
		return 0, errors.WithCaller("rows affected should be 1", 2)
	}
	return result.LastInsertId()
}

// TODO
// NewTestBlog creates a new user, and creates a new blog by the created user

// UpdateBlog update the blog where user name and blog id are matched
// PUT /blog
func UpdateBlog(c *elesion.Context) {
	usr, err := user.CheckAuth(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String(err.Error())
		return
	}

	body := Blog{}
	body.User = usr.ID
	err = json.NewDecoder(c.Request.Body).Decode(&body)
	if err != nil {
		c.Status(http.StatusBadRequest).Error(err.Error()).String("invalid body")
		return
	}
	defer c.Request.Body.Close()

	err = updateBlog(body)

	switch err {
	case nil:
		c.Status(http.StatusOK).String("success")
	case db.ErrNoChange:
		c.Status(http.StatusAccepted).String("no change")
	case db.ErrNoRows:
		c.Status(http.StatusNotFound).String(fmt.Sprintf("no such blog: %d", body.ID))
	default:
		c.Status(http.StatusInternalServerError).Error(err.Error()).String(err.Error())
	}
}

func updateBlog(blog Blog) error {
	d := db.UseDefault()
	defer d.Close()

	ok, err := d.Exists("SELECT 1 FROM blog WHERE id = ? AND user = ?", blog.ID, blog.User)
	if err != nil {
		return err
	}
	if !ok {
		return db.ErrNoRows
	}

	res, err := d.Exec("UPDATE blog SET title = ?, text = ? WHERE id = ? AND user = ?",
		blog.Title, blog.Text, blog.ID, blog.User,
	)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); err != nil {
		return err
	} else if rows == 0 {
		return db.ErrNoChange
	} else if rows > 1 {
		return errors.Newf("rows affected should be larger than 1 but got ", rows)
	}
	return nil
}

// DeleteBlog delete the blog where user name and blog id are matched
// DELETE /blog
func DeleteBlog(c *elesion.Context) {
	usr, err := user.CheckAuth(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.Status(http.StatusUnauthorized).String(err.Error())
		return
	}

	blog := Blog{}
	err = json.NewDecoder(c.Request.Body).Decode(&blog)
	if err != nil {
		c.Status(http.StatusBadRequest).String("invald body")
		return
	}

	err = deleteBlog(usr.ID, blog.ID)
	switch err {
	case nil:
		c.Status(http.StatusOK).String("success")
	case db.ErrNoRows:
		c.Status(http.StatusNotFound).Stringf("blog not found")
	default:
		c.Status(http.StatusInternalServerError).Error(err.Error()).String(err.Error())
	}
}

func deleteBlog(user, id int64) error {
	d := db.UseDefault()
	defer d.Close()

	ok, err := d.Exists("SELECT 1 FROM blog WHERE id = ? AND user = ?", id, user)
	if err != nil {
		return err
	}
	if !ok {
		return db.ErrNoRows
	}

	result, err := d.Exec("DELETE FROM blog WHERE id = ? AND user = ?", id, user)
	if err != nil {
		return err
	}
	if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows == 0 {
		return db.ErrNoChange
	} else if rows > 1 {
		return errors.Newf("rows affected shoulb be 1 but got %d", rows)
	}
	return nil
}
