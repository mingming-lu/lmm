package http

import (
	"fmt"
	"net/http"

	"lmm/api/pkg/auth"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine/log"
)

// BadRequest default response
func BadRequest(c *gin.Context) {
	ErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
}

// Unauthorized default response
func Unauthorized(c *gin.Context) {
	ErrorResponse(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
}

// Forbidden default response
func Forbidden(c *gin.Context) {
	ErrorResponse(c, http.StatusForbidden, http.StatusText(http.StatusForbidden))
}

// NotFound default response
func NotFound(c *gin.Context) {
	ErrorResponse(c, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

// Response writes message into c in default format
func Response(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"message": msg,
	})
}

// ErrorResponse writes default error response into c by given code
func ErrorResponse(c *gin.Context, code int, errMsg string) {
	c.JSON(code, gin.H{
		"error": errMsg,
	})
}

func AuthFromGinContext(c *gin.Context) (*auth.Auth, bool) {
	return auth.FromContext(c.Request.Context())
}

func LogDebugf(c *gin.Context, format string, args ...interface{}) {
	log.Debugf(c.Request.Context(), format, args...)
}

func LogInfo(c *gin.Context, format string, args ...interface{}) {
	log.Infof(c.Request.Context(), format, args...)
}

func LogErrorf(c *gin.Context, format string, args ...interface{}) {
	log.Errorf(c.Request.Context(), format, args...)
}

func LogWarnf(c *gin.Context, format string, args ...interface{}) {
	log.Warningf(c.Request.Context(), format, args...)
}

func LogCritf(c *gin.Context, format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}
