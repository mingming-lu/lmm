package ui

import (
	"fmt"
	"math"
	"net/http"

	httpUtil "lmm/api/pkg/http"
	"lmm/api/pkg/transaction"
	"lmm/api/service/article/application"
	"lmm/api/service/article/application/command"
	"lmm/api/service/article/application/query"
	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/model"
	"lmm/api/service/article/domain/repository"
	"lmm/api/service/article/domain/viewer"
	"lmm/api/util/stringutil"

	"github.com/gin-gonic/gin"
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
func (ui *UI) PostNewArticle(c *gin.Context) {
	user, ok := httpUtil.AuthFromGinContext(c)
	if !ok {
		httpUtil.Unauthorized(c)
		return
	}

	article := postArticleAdapter{}
	if err := c.ShouldBindJSON(&article); err != nil {
		httpUtil.BadRequest(c)
	}

	if err := ui.validatePostArticleAdaptor(&article); err != nil {
		httpUtil.ErrorResponse(c, http.StatusBadRequest, err.Error())
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
		httpUtil.Unauthorized(c)
	default:
		httpUtil.LogCritf(c, err.Error())
	}
}

// PutV1Articles handles PUT /v1/article/:articleID
func (ui *UI) PutV1Articles(c *gin.Context) {
	user, ok := httpUtil.AuthFromGinContext(c)
	if !ok {
		httpUtil.Unauthorized(c)
		return
	}

	articleID, err := stringutil.ParseInt64(c.Param("articleID"))
	if err != nil {
		c.String(http.StatusNotFound, domain.ErrNoSuchArticle.Error())
		return
	}

	article := postArticleAdapter{}
	if err := c.ShouldBindJSON(&article); err != nil {
		httpUtil.BadRequest(c)
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
		httpUtil.Response(c, http.StatusOK, "Success")

	case
		domain.ErrArticleTitleTooLong,
		domain.ErrEmptyArticleTitle,
		domain.ErrInvalidArticleTitle,
		domain.ErrInvalidAliasArticleID:

		c.String(http.StatusBadRequest, original.Error())

	case domain.ErrNoSuchArticle, domain.ErrInvalidArticleID:
		c.String(http.StatusNotFound, domain.ErrNoSuchArticle.Error())

	case domain.ErrNoSuchUser:
		httpUtil.Unauthorized(c)

	case domain.ErrNotArticleAuthor:
		c.String(http.StatusForbidden, original.Error())

	default:
		httpUtil.LogCritf(c, err.Error())
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
func (ui *UI) ListArticles(c *gin.Context) {
	v, err := ui.appService.Query().ListArticlesByPage(
		c,
		ui.buildListArticleQueryFromContext(c),
	)
	switch errors.Cause(err) {
	case nil:
		if c.DefaultQuery("flavor", "") == "true" {
			c.JSON(http.StatusOK, ui.articleListViewToJSONV2(c, v))
		} else {
			c.JSON(http.StatusOK, ui.articleListViewToJSON(v))
		}
	case application.ErrInvalidCount, application.ErrInvalidPage:
		c.JSON(http.StatusBadRequest, err.Error())
	default:
		httpUtil.LogCritf(c, err.Error())
	}
}

func (ui *UI) buildListArticleQueryFromContext(c *gin.Context) query.ListArticleQuery {
	var tag *string
	if c.Query("tag") == "" {
		tmp := c.Query("tag")
		tag = &tmp
	}
	return query.ListArticleQuery{
		Page:  c.DefaultQuery("page", "1"),
		Count: c.DefaultQuery("perPage", "5"),
		Tag:   tag,
	}
}

func (ui *UI) articleListViewToJSON(view *model.ArticleListView) *articleListAdapter {
	items := make([]articleListItem, len(view.Items()), len(view.Items()))
	for i, item := range view.Items() {
		items[i].ID = item.ID()
		items[i].Link = item.LinkName()
		items[i].Title = item.Title()
		items[i].PostAt = item.PostAt().Unix()
	}
	return &articleListAdapter{
		Articles:    items,
		HasNextPage: view.HasNextPage(),
	}
}

func (ui *UI) articleListViewToJSONV2(c *gin.Context, view *model.ArticleListView) *articleListAdapterV2 {
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
		adapterV2.PrevPage = buildURI(c.Request.URL.Path, adapterV2.Page-1, adapterV2.PerPage)
	}

	if view.HasNextPage() {
		adapterV2.NextPage = buildURI(c.Request.URL.Path, adapterV2.Page+1, adapterV2.PerPage)
	}

	if adapterV2.Total > 0 {
		adapterV2.FirstPage = buildURI(c.Request.URL.Path, 1, adapterV2.PerPage)
		adapterV2.LastPage = buildURI(c.Request.URL.Path, last, adapterV2.PerPage)
	}

	return adapterV2
}

func buildURI(path string, page, perPage int) *string {
	uri := fmt.Sprintf("%s?page=%d&perPage=%d", path, page, perPage)
	return &uri
}

// GetArticle handles GET /v1/articles/:articleID
func (ui *UI) GetArticle(c *gin.Context) {
	view, err := ui.appService.Query().ArticleByLinkName(c,
		c.Param("articleID"),
	)
	switch errors.Cause(err) {
	case nil:
		c.JSON(http.StatusOK, ui.articleViewToJSON(view))
	case domain.ErrInvalidArticleID, domain.ErrNoSuchArticle:
		c.String(http.StatusNotFound, domain.ErrNoSuchArticle.Error())
	default:
		httpUtil.LogCritf(c, err.Error())
	}
}

func (ui *UI) articleViewToJSON(model *model.Article) *articleViewResponse {
	tags := make([]articleViewTag, len(model.Content().Tags()), len(model.Content().Tags()))
	for i, tag := range model.Content().Tags() {
		tags[i].Name = tag
	}
	return &articleViewResponse{
		ID:           model.ID().ID(),
		Link:         model.LinkName(),
		Title:        model.Content().Text().Title(),
		Body:         model.Content().Text().Body(),
		PostAt:       model.CreatedAt().Unix(),
		LastEditedAt: model.LastModified().Unix(),
		Tags:         tags,
	}
}

// GetAllArticleTags handles GET /v1/articleTags
func (ui *UI) GetAllArticleTags(c *gin.Context) {
	tags, err := ui.appService.Query().AllArticleTags(c)

	switch errors.Cause(err) {
	case nil:
		c.JSON(http.StatusOK, ui.tagListViewToJSON(tags))
	default:
		httpUtil.LogCritf(c, err.Error())
	}
}

func (ui *UI) tagListViewToJSON(views []*model.TagView) articleTagListView {
	tags := make([]articleTagListItemView, len(views), len(views))
	for i, tag := range views {
		tags[i].Name = tag.Name()
	}
	return tags
}
