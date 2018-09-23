package application

import "errors"

var (
	ErrInvalidCount = errors.New("invalid count")
	ErrInvalidPage  = errors.New("invalid page")
)

// Service is like a registry for services in article bounded context
type Service struct {
	articleCommandService *ArticleCommandService
	articleQueryService   *ArticleQueryService
}

// NewService is a constructor of Service
func NewService(articleCommandService *ArticleCommandService, articleQueryService *ArticleQueryService) *Service {
	return &Service{
		articleCommandService: articleCommandService,
		articleQueryService:   articleQueryService,
	}
}

// ArticleCommandService getter
func (s *Service) ArticleCommandService() *ArticleCommandService {
	return s.articleCommandService
}

// ArticleQueryService getter
func (s *Service) ArticleQueryService() *ArticleQueryService {
	return s.articleQueryService
}
