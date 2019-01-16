package ui

import (
	"context"
	"encoding/json"
	"time"

	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/util/testutil"

	"github.com/google/uuid"
)

func TestListArticleV1(tt *testing.T) {
	lock.Lock()
	defer lock.Unlock()

	t := testing.NewTester(tt)
	c := context.Background()

	user := testutil.NewUser(mysql)

	if _, err := mysql.Exec(context.Background(), "truncate table article"); err != nil {
		t.FailNow()
	}

	var articles [9]articleListItem

	for i := range articles {
		articles[i].ID = uuid.New().String()[:8]
		articles[i].Title = uuid.New().String()

		postAtUnix := time.Now().Add(10 - time.Duration(i)*time.Minute).Unix()
		articles[i].PostAt = postAtUnix

		postAtFromUnix := time.Unix(postAtUnix, 0)

		if _, err := mysql.Exec(c,
			"insert into article (uid, user, title, body, created_at, updated_at) values(?, ?, ?, ?, ?, ?)",
			articles[i].ID, user.ID(), articles[i].Title, uuid.New().String(), postAtFromUnix, postAtFromUnix,
		); err != nil {
			t.Log(err.Error())
		}
	}
	cases := map[string]articleListAdapter{
		"": {
			Articles:    articles[:5],
			HasNextPage: true,
		},
		"?page=2": {
			Articles:    articles[5:],
			HasNextPage: false,
		},
		"?page=3": {
			Articles:    []articleListItem{},
			HasNextPage: false,
		},
		"?page=100": {
			Articles:    []articleListItem{},
			HasNextPage: false,
		},
		"?perPage=10": {
			Articles:    articles[:],
			HasNextPage: false,
		},
		"?page=1&perPage=10": {
			Articles:    articles[:],
			HasNextPage: false,
		},
		"?page=2&perPage=10": {
			Articles:    []articleListItem{},
			HasNextPage: false,
		},
		"?page=100&perPage=10": {
			Articles:    []articleListItem{},
			HasNextPage: false,
		},
		"?perPage=1": {
			Articles:    articles[:1],
			HasNextPage: true,
		},
		"?perPage=1&page=5": {
			Articles:    articles[4:5],
			HasNextPage: true,
		},
		"?perPage=1&page=10": {
			Articles:    []articleListItem{},
			HasNextPage: false,
		},
		"?perPage=2": {
			Articles:    articles[:2],
			HasNextPage: true,
		},
		"?perPage=2&page=5": {
			Articles:    articles[8:],
			HasNextPage: false,
		},
	}

	for testName, testCase := range cases {
		t.Run(testName, func(_ *testing.T) {
			req := testing.GET("/v1/articles"+testName, nil)
			res := testing.DoRequest(req, router)

			t.Is(http.StatusOK, res.StatusCode())

			body := articleListAdapter{}
			t.NoError(json.NewDecoder(res.RawBody()).Decode(&body))

			t.Is(testCase, body)
		})
	}
}
