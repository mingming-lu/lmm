package ui

import (
	"fmt"
	"math"

	"lmm/api/http"
	"lmm/api/pkg/auth"
	"lmm/api/pkg/transaction"
	"lmm/api/service/article/application"
	"lmm/api/service/article/application/command"
	"lmm/api/service/article/application/query"
	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/model"
	"lmm/api/service/article/domain/repository"
	"lmm/api/service/article/domain/viewer"
	"lmm/api/util/stringutil"

	"github.com/pkg/errors"
)

var (
	errTitleRequired = errors.New("title required")
	errBodyRequired  = errors.New("body requried")
	errTagsRequired  = errors.New("tags requried")
)

// UI is the user interface to contact with network
type UI struct {
	appService *application.Service
}

// NewUI returns a new ui
func NewUI(
	articleViewer viewer.ArticleViewer,
	articleRepository repository.ArticleRepository,
	transactionManager transaction.Manager,
) *UI {
	appService := application.NewService(
		application.NewArticleCommandService(articleRepository, transactionManager),
		application.NewArticleQueryService(articleViewer, transactionManager),
	)
	return &UI{appService: appService}
}

// PostNewArticle handles POST /1/articles
func (ui *UI) PostNewArticle(c http.Context) {
	user, ok := auth.FromContext(c)
	if !ok {
		http.Unauthorized(c)
		return
	}

	article := postArticleAdapter{}
	if err := c.Request().Bind(&article); err != nil {
		http.BadRequest(c)
		return
	}

	if err := ui.validatePostArticleAdaptor(&article); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	articleID, err := ui.appService.Command().PostNewArticle(c, command.PostArticle{
		AuthorID: user.ID,
		Title:    *article.Title,
		Body:     *article.Body,
		Tags:     article.Tags,
	})
	switch errors.Cause(err) {
	case nil:
		c.Header("Location", fmt.Sprintf("/v1/articles/%d", articleID.ID()))
		c.String(http.StatusCreated, "Success")
	case domain.ErrArticleTitleTooLong, domain.ErrEmptyArticleTitle:
		c.String(http.StatusBadRequest, err.Error())
	case domain.ErrInvalidArticleTitle:
		c.String(http.StatusBadRequest, err.Error())
	case domain.ErrNoSuchUser:
		http.Unauthorized(c)
	default:
		http.Log().Panic(c, err.Error())
	}
}

// PutV1Articles handles PUT /v1/article/:articleID
func (ui *UI) PutV1Articles(c http.Context) {
	user, ok := auth.FromContext(c)
	if !ok {
		http.Unauthorized(c)
		return
	}

	articleID, err := stringutil.ParseInt64(c.Request().PathParam("articleID"))
	if err != nil {
		c.String(http.StatusNotFound, domain.ErrNoSuchArticle.Error())
		return
	}

	article := postArticleAdapter{}
	if err := c.Request().Bind(&article); err != nil {
		http.BadRequest(c)
		return
	}

	if err := ui.validatePostArticleAdaptor(&article); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = ui.appService.Command().EditArticle(c, command.EditArticle{
		UserID:    user.ID,
		ArticleID: articleID,
		LinkName:  *article.LinkName,
		Title:     *article.Title,
		Body:      *article.Body,
		Tags:      article.Tags,
	})

	original := errors.Cause(err)
	switch original {
	case nil:
		http.NoContent(c)

	case
		domain.ErrArticleTitleTooLong,
		domain.ErrEmptyArticleTitle,
		domain.ErrInvalidArticleTitle,
		domain.ErrInvalidAliasArticleID:

		c.String(http.StatusBadRequest, original.Error())

	case domain.ErrNoSuchArticle, domain.ErrInvalidArticleID:
		c.String(http.StatusNotFound, domain.ErrNoSuchArticle.Error())

	case domain.ErrNoSuchUser:
		http.Unauthorized(c)

	case domain.ErrNotArticleAuthor:
		c.String(http.StatusForbidden, original.Error())

	default:
		http.Log().Panic(c, err.Error())
	}
}

func (ui *UI) validatePostArticleAdaptor(adaptor *postArticleAdapter) error {
	if adaptor.Title == nil {
		return errTitleRequired
	}
	if adaptor.Body == nil {
		return errBodyRequired
	}
	if adaptor.Tags == nil {
		return errTagsRequired
	}
	return nil
}

// ListArticles handles GET /v1/articles
func (ui *UI) ListArticles(c http.Context) {
	v, err := ui.appService.Query().ListArticlesByPage(
		c,
		ui.buildListArticleQueryFromContext(c),
	)
	switch errors.Cause(err) {
	case nil:
		if c.Request().QueryParamOrDefault("flavor", "") == "true" {
			c.JSON(http.StatusOK, ui.articleListViewToJSONV2(c, v))
		} else {
			c.JSON(http.StatusOK, ui.articleListViewToJSON(v))
		}
	case application.ErrInvalidCount, application.ErrInvalidPage:
		c.JSON(http.StatusBadRequest, err.Error())
	default:
		http.Log().Panic(c, err.Error())
	}
}

func (ui *UI) buildListArticleQueryFromContext(c http.Context) query.ListArticleQuery {
	return query.ListArticleQuery{
		Page:  c.Request().QueryParamOrDefault("page", "1"),
		Count: c.Request().QueryParamOrDefault("perPage", "5"),
		Tag:   c.Request().QueryParam("tag"),
	}
}

func (ui *UI) articleListViewToJSON(view *model.ArticleListView) *articleListAdapter {
	items := make([]articleListItem, len(view.Items()), len(view.Items()))
	for i, item := range view.Items() {
		items[i].ID = fmt.Sprintf("%d", item.ID())
		items[i].Title = item.Title()
		items[i].PostAt = item.PostAt().Unix()
	}
	return &articleListAdapter{
		Articles:    items,
		HasNextPage: view.HasNextPage(),
	}
}

func (ui *UI) articleListViewToJSONV2(c http.Context, view *model.ArticleListView) *articleListAdapterV2 {
	adapter := ui.articleListViewToJSON(view)
	adapterV2 := &articleListAdapterV2{
		Articles: adapter.Articles,
		Page:     view.Page(),
		PerPage:  view.PerPage(),
		Total:    view.Total(),
	}

	last := int(math.Ceil(
		float64(adapterV2.Total) / float64(adapterV2.PerPage),
	))

	if adapterV2.Page > 1 && adapterV2.Page <= last+1 {
		adapterV2.PrevPage = buildURI(c.Request().Path(), adapterV2.Page-1, adapterV2.PerPage)
	}

	if view.HasNextPage() {
		adapterV2.NextPage = buildURI(c.Request().Path(), adapterV2.Page+1, adapterV2.PerPage)
	}

	if adapterV2.Total > 0 {
		adapterV2.FirstPage = buildURI(c.Request().Path(), 1, adapterV2.PerPage)
		adapterV2.LastPage = buildURI(c.Request().Path(), last, adapterV2.PerPage)
	}

	return adapterV2
}

func buildURI(path string, page, perPage int) *string {
	uri := fmt.Sprintf("%s?page=%d&perPage=%d", path, page, perPage)
	return &uri
}

// GetArticle handles GET /v1/articles/:articleID
func (ui *UI) GetArticle(c http.Context) {
	view, err := ui.appService.Query().ArticleByLinkName(c,
		c.Request().PathParam("articleID"),
	)
	switch errors.Cause(err) {
	case nil:
		c.JSON(http.StatusOK, ui.articleViewToJSON(view))
	case domain.ErrInvalidArticleID, domain.ErrNoSuchArticle:
		c.String(http.StatusNotFound, domain.ErrNoSuchArticle.Error())
	default:
		http.Log().Panic(c, err.Error())
	}
}

func (ui *UI) articleViewToJSON(model *model.Article) *articleViewResponse {
	tags := make([]articleViewTag, len(model.Content().Tags()), len(model.Content().Tags()))
	for i, tag := range model.Content().Tags() {
		tags[i].Name = tag
	}
	return &articleViewResponse{
		ID:           fmt.Sprintf("%d", model.ID().ID()),
		Title:        model.Content().Text().Title(),
		Body:         model.Content().Text().Body(),
		PostAt:       model.CreatedAt().Unix(),
		LastEditedAt: model.LastModified().Unix(),
		Tags:         tags,
	}
}

// GetAllArticleTags handles GET /v1/articleTags
func (ui *UI) GetAllArticleTags(c http.Context) {
	tags, err := ui.appService.Query().AllArticleTags(c)

	switch errors.Cause(err) {
	case nil:
		c.JSON(http.StatusOK, ui.tagListViewToJSON(tags))
	default:
		http.Log().Panic(c, err.Error())
	}
}

func (ui *UI) tagListViewToJSON(views []*model.TagView) articleTagListView {
	tags := make([]articleTagListItemView, len(views), len(views))
	for i, tag := range views {
		tags[i].Name = tag.Name()
	}
	return tags
}
