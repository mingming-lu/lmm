package middleware

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS middleware
func CORS(whiteList ...string) gin.HandlerFunc {
	sort.Strings(whiteList)

	return cors.New(cors.Config{
		AllowMethods:  []string{http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:  []string{"Authorization", "Content-Type", "Origin"},
		ExposeHeaders: []string{"Content-Length", "Location"},
		AllowOriginFunc: func(origin string) bool {
			idx := sort.SearchStrings(whiteList, origin)

			return idx < len(whiteList) && whiteList[idx] == origin
		},
		MaxAge: 24 * time.Hour,
	})
}
