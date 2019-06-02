package middleware

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

func WrapAppEngineContext(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	c.Request = c.Request.WithContext(ctx)

	c.Next()
}
