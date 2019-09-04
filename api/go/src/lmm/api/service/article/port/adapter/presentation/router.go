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

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	errTitleRequired = errors.New("title required")
	errBodyRequired  = errors.New("body requried")
	errTagsRequired  = errors.New("tags requried")
)

type GinRouterProvider struct {
	appService *application.Service
}

func NewGinRouterProvider(
	articleViewer model.ArticleViewer,
	articleRepository model.ArticleRepository,
	transactionManager transaction.Manager,
) *GinRouterProvider {
	appService := application.NewService(
		application.NewArticleCommandService(articleRepository, transactionManager),
		application.NewArticleQueryService(articleViewer, transactionManager),
	)
	return &GinRouterProvider{appService: appService}
}

func (p *GinRouterProvider) Provide(router *gin.Engine) {
	router.POST("/v1/articles", p.PostNewArticle)
	router.PUT("/v1/articles/:articleID", p.PutV1Articles)
	router.GET("/v1/articles", p.ListArticles)
	router.GET("/v1/articles/:articleID", p.GetArticle)
	router.GET("/v1/articleTags", p.GetAllArticleTags)
}

// PostNewArticle handles POST /1/articles
func (p *GinRouterProvider) PostNewArticle(c *gin.Context) {
	user, ok := httpUtil.AuthFromGinContext(c)
	if !ok {
		httpUtil.Unauthorized(c)
		return
	}

	article := postArticleAdapter{}
	if err := c.ShouldBindJSON(&article); err != nil {
		httpUtil.BadRequest(c)
	}

	if err := p.validatePostArticleAdaptor(&article); err != nil {
		httpUtil.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	articleID, err := p.appService.Command().PostNewArticle(c, command.PostArticle{
		AuthorID: user.ID,
		Title:    *article.Title,
		Body:     *article.Body,
		Tags:     article.Tags,
	})
	originalErr := errors.Cause(err)
	switch originalErr {
	case nil:
		c.Header("Location", fmt.Sprintf("/v1/articles/%s", articleID.String()))
		httpUtil.Response(c, http.StatusCreated, "Success")
	case domain.ErrArticleTitleTooLong, domain.ErrEmptyArticleTitle, domain.ErrInvalidArticleTitle:
		httpUtil.ErrorResponse(c, http.StatusBadRequest, originalErr.Error())
	case domain.ErrNoSuchUser:
		httpUtil.Unauthorized(c)
	default:
		httpUtil.LogPanic(c, "unexpected error", err)
	}
}

// PutV1Articles handles PUT /v1/article/:articleID
func (p *GinRouterProvider) PutV1Articles(c *gin.Context) {
	user, ok := httpUtil.AuthFromGinContext(c)
	if !ok {
		httpUtil.Unauthorized(c)
		return
	}

	article := postArticleAdapter{}
	if err := c.ShouldBindJSON(&article); err != nil {
		httpUtil.BadRequest(c)
		return
	}

	if err := p.validatePostArticleAdaptor(&article); err != nil {
		httpUtil.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := p.appService.Command().EditArticle(c, command.EditArticle{
		UserID:    user.ID,
		ArticleID: c.Param("articleID"),
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

		httpUtil.ErrorResponse(c, http.StatusBadRequest, original.Error())

	case domain.ErrNoSuchArticle, domain.ErrInvalidArticleID:
		httpUtil.ErrorResponse(c, http.StatusNotFound, domain.ErrNoSuchArticle.Error())

	case domain.ErrNoSuchUser:
		httpUtil.Unauthorized(c)

	case domain.ErrNotArticleAuthor:
		httpUtil.ErrorResponse(c, http.StatusForbidden, original.Error())

	default:
		httpUtil.LogPanic(c, "unexpected error", err)
	}
}

func (p *GinRouterProvider) validatePostArticleAdaptor(adaptor *postArticleAdapter) error {
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
func (p *GinRouterProvider) ListArticles(c *gin.Context) {
	q := query.ListArticleQuery{}
	if err := c.BindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": q.ValidateErrors(err),
		})
		return
	}

	v, err := p.appService.Query().ListArticlesByPage(c, q)
	switch errors.Cause(err) {
	case nil:
		if c.DefaultQuery("flavor", "") == "true" {
			c.JSON(http.StatusOK, p.articleListViewToJSONV2(c, v))
		} else {
			c.JSON(http.StatusOK, p.articleListViewToJSON(v))
		}
	default:
		httpUtil.LogPanic(c, "unexpected error", err)
	}
}

func (p *GinRouterProvider) articleListViewToJSON(view *model.ArticleListView) *articleListAdapter {
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

func (p *GinRouterProvider) articleListViewToJSONV2(c *gin.Context, view *model.ArticleListView) *articleListAdapterV2 {
	adapter := p.articleListViewToJSON(view)
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
func (p *GinRouterProvider) GetArticle(c *gin.Context) {
	view, err := p.appService.Query().ArticleByID(c,
		c.Param("articleID"),
	)
	switch errors.Cause(err) {
	case nil:
		c.JSON(http.StatusOK, p.articleViewToJSON(view))
	case domain.ErrInvalidArticleID, domain.ErrNoSuchArticle:
		httpUtil.ErrorResponse(c, http.StatusNotFound, domain.ErrNoSuchArticle.Error())
	default:
		httpUtil.LogPanic(c, "unexpected error", err)
	}
}

func (p *GinRouterProvider) articleViewToJSON(model *model.Article) *articleViewResponse {
	tags := make([]articleViewTag, len(model.Content().Tags()), len(model.Content().Tags()))
	for i, tag := range model.Content().Tags() {
		tags[i].Name = tag.Name()
	}
	return &articleViewResponse{
		ID:           model.ID().String(),
		Title:        model.Content().Text().Title(),
		Body:         model.Content().Text().Body(),
		PostAt:       model.CreatedAt().Unix(),
		LastEditedAt: model.LastModified().Unix(),
		Tags:         tags,
	}
}

// GetAllArticleTags handles GET /v1/articleTags
func (p *GinRouterProvider) GetAllArticleTags(c *gin.Context) {
	tags, err := p.appService.Query().AllArticleTags(c)

	switch errors.Cause(err) {
	case nil:
		c.JSON(http.StatusOK, p.tagListViewToJSON(tags))
	default:
		httpUtil.LogPanic(c, "unexpected error", err)
	}
}

func (p *GinRouterProvider) tagListViewToJSON(views []*model.TagView) articleTagListView {
	tags := make([]articleTagListItemView, len(views), len(views))
	for i, tag := range views {
		tags[i].Name = tag.Name()
	}
	return tags
}
