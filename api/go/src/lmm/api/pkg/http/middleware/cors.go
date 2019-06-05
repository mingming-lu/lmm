package middleware

import (
	"os"
	"sort"

	"github.com/gin-gonic/gin"
)

var whiteList []string

func init() {
	if appOrigin := os.Getenv("APP_ORIGIN"); appOrigin != "" {
		whiteList = append(whiteList, appOrigin)
	}
	if managerOrigin := os.Getenv("MANAGER_ORIGIN"); managerOrigin != "" {
		whiteList = append(whiteList, managerOrigin)
	}

	sort.Strings(whiteList)
}

// CORS middleware
func CORS(c *gin.Context) {
	c.Next()

	origin := c.GetHeader("Origin")

	if origin == "" {
		return
	}

	if idx := sort.SearchStrings(whiteList, origin); idx < len(whiteList) && whiteList[idx] == origin {
		c.Header("Access-Control-Allow-Origin", origin)
	}
}
