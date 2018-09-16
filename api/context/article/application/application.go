package application

// Service is like a registry for services in article bounded context
type Service struct {
	articleApplicationService *ArticleApplicationService
	articleQueryService       *ArticleQueryService
}

// NewService is a constructor of Service
func NewService(articleApplicationService *ArticleApplicationService, articleQueryService *ArticleQueryService) *Service {
	return &Service{
		articleApplicationService: articleApplicationService,
		articleQueryService:       articleQueryService,
	}
}

// ArticleApplicationService getter
func (s *Service) ArticleApplicationService() *ArticleApplicationService {
	return s.articleApplicationService
}

// ArticleQueryService getter
func (s *Service) ArticleQueryService() *ArticleQueryService {
	return s.articleQueryService
}
