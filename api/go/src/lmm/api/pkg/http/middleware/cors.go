package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/proproto/goenv"
)

// CORS middleware
func CORS(customDomain string) gin.HandlerFunc {
	var re *regexp.Regexp
	if customDomain == "" {
		config := struct {
			ProjectID string `env:"GCP_PROJECT_ID,required"`
		}{}
		goenv.MustBind(&config)
		pattern := fmt.Sprintf(`^https://(.+-dot-)*%s.appspot\.com$`, config.ProjectID)
		re = regexp.MustCompile(pattern)
	} else {
		pattern := fmt.Sprintf(`^https://(.+\.)*%s$`, regexp.QuoteMeta(customDomain))
		re = regexp.MustCompile(pattern)
	}

	return cors.New(cors.Config{
		AllowMethods:  []string{http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:  []string{"Authorization", "Content-Type", "Origin"},
		ExposeHeaders: []string{"Content-Length", "Location"},
		AllowOriginFunc: func(origin string) bool {
			return re.MatchString(origin)
		},
		MaxAge: 24 * time.Hour,
	})
}
