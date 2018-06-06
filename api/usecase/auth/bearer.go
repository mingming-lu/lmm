package auth

import (
	accountService "lmm/api/context/account/appservice"
	accountRepository "lmm/api/context/account/domain/repository"
	"lmm/api/http"
	"regexp"
)

var (
	PatternBearerAuthorization = regexp.MustCompile(`^Bearer (.+)$`)
)

func BearerAuth(handler http.Handler) http.Handler {
	return func(c *http.Context) {
		auth := c.Request.Header.Get("Authorization")
		mathed := PatternBearerAuthorization.FindStringSubmatch(auth)
		if len(mathed) != 2 {
			http.Unauthorized(c)
			return
		}
		token := mathed[1]
		user, err := accountService.New(accountRepository.New()).VerifyToken(token)
		if err != nil {
			http.Unauthorized(c)
			return
		}
		c.Values().Set("user", user)
		handler(c)
	}
}
