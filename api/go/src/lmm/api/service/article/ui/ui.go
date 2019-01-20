package ui

import (
	"fmt"
	"math"

	"github.com/pkg/errors"

	"lmm/api/http"
	"lmm/api/service/article/application"
	"lmm/api/service/article/application/query"
	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/finder"
	"lmm/api/service/article/domain/model"
	"lmm/api/service/article/domain/repository"
	"lmm/api/service/article/domain/service"
	"lmm/api/util/stringutil"
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
	articleFinder finder.ArticleFinder,
	articleRepository repository.ArticleRepository,
	authorService service.AuthorService,
) *UI {
	appService := application.NewService(
		application.NewArticleCommandService(articleRepository, authorService),
		application.NewArticleQueryService(articleFinder),
	)
	return &UI{appService: appService}
}

// PostNewArticle handles POST /1/articles
func (ui *UI) PostNewArticle(c http.Context) {
	userName := c.Request().Header.Get("X-LMM-ID")
	if userName == "" {
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

	articleID, err := ui.appService.ArticleCommandService().PostNewArticle(c,
		userName,
		*article.Title,
		*article.Body,
		article.Tags,
	)
	switch errors.Cause(err) {
	case nil:
		c.Header("Location", "/v1/articles/"+articleID.String())
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

// EditArticle handles PUT /1/article/:articleID
func (ui *UI) EditArticle(c http.Context) {
	userName := c.Request().Header.Get("X-LMM-ID")
	if userName == "" {
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

	err := ui.appService.ArticleCommandService().EditArticle(c,
		userName,
		c.Request().PathParam("articleID"),
		*article.Title,
		*article.Body,
		article.Tags,
	)
	switch errors.Cause(err) {
	case nil:
		http.NoContent(c)
	case domain.ErrArticleTitleTooLong, domain.ErrEmptyArticleTitle:
		c.String(http.StatusBadRequest, err.Error())
	case domain.ErrInvalidArticleID:
		c.String(http.StatusNotFound, domain.ErrNoSuchArticle.Error())
	case domain.ErrInvalidArticleTitle:
		c.String(http.StatusBadRequest, err.Error())
	case domain.ErrNoSuchArticle:
		c.String(http.StatusNotFound, err.Error())
	case domain.ErrNoSuchUser:
		http.Unauthorized(c)
	case domain.ErrNotArticleAuthor:
		c.String(http.StatusForbidden, err.Error())
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
	v, err := ui.appService.ArticleQueryService().ListArticlesByPage(
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
		items[i].ID = item.ID().String()
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

	last := uint(math.Ceil(
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

func buildURI(base string, page, perPage uint) *string {
	return stringutil.Pointer(fmt.Sprintf("%s?page=%d&perPage=%d", base, page, perPage))
}

// GetArticle handles GET /v1/articles/:articleID
func (ui *UI) GetArticle(c http.Context) {
	view, err := ui.appService.ArticleQueryService().ArticleByID(c,
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

func (ui *UI) articleViewToJSON(view *model.ArticleView) *articleViewResponse {
	tags := make([]articleViewTag, len(view.Content().Tags()), len(view.Content().Tags()))
	for i, tag := range view.Content().Tags() {
		tags[i].Name = tag.Name()
	}
	return &articleViewResponse{
		ID:           view.ID().String(),
		Title:        view.Content().Text().Title(),
		Body:         view.Content().Text().Body(),
		PostAt:       view.PostAt().Unix(),
		LastEditedAt: view.LastEditedAt().Unix(),
		Tags:         tags,
	}
}

// GetAllArticleTags handles GET /v1/articleTags
func (ui *UI) GetAllArticleTags(c http.Context) {
	view, err := ui.appService.ArticleQueryService().AllArticleTags(c)

	switch errors.Cause(err) {
	case nil:
		c.JSON(http.StatusOK, ui.tagListViewToJSON(view))
	default:
		http.Log().Panic(c, err.Error())
	}
}

func (ui *UI) tagListViewToJSON(view model.TagListView) articleTagListView {
	tags := make([]articleTagListItemView, len(view), len(view))
	for i, tag := range view {
		tags[i].Name = tag.Name()
	}
	return tags
}
