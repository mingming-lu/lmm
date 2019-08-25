package http

import (
	"net/http"

	"go.uber.org/zap"

	"google.golang.org/appengine"

	"lmm/api/pkg/auth"

	"github.com/gin-gonic/gin"
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

var logger = zap.NewExample()

func AuthFromGinContext(c *gin.Context) (*auth.Auth, bool) {
	return auth.FromContext(c.Request.Context())
}

func reqIDFromC(c *gin.Context) string {
	ctx := appengine.NewContext(c.Request)

	return appengine.RequestID(ctx)
}

func LogDebug(c *gin.Context, msg string, err error) {
	logger.Debug(msg,
		zap.String("reqID", reqIDFromC(c)),
		zap.String("handler", c.HandlerName()),
		zap.Error(err),
	)
}

func LogInfo(c *gin.Context, msg string, err error) {
	logger.Info(msg,
		zap.String("reqID", reqIDFromC(c)),
		zap.String("handler", c.HandlerName()),
		zap.Error(err),
	)
}

func LogError(c *gin.Context, msg string, err error) {
	logger.Error(msg,
		zap.String("reqID", reqIDFromC(c)),
		zap.String("handler", c.HandlerName()),
		zap.Error(err),
	)
}

func LogWarn(c *gin.Context, msg string, err error) {
	logger.Warn(msg,
		zap.String("reqID", reqIDFromC(c)),
		zap.String("handler", c.HandlerName()),
		zap.Error(err),
	)
}

func LogPanic(c *gin.Context, msg string, err error) {
	logger.Panic(msg,
		zap.String("reqID", reqIDFromC(c)),
		zap.String("handler", c.HandlerName()),
		zap.Error(err),
	)
}
