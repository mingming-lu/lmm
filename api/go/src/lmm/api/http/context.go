package http

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// StrCtxKey stands for string context key
type StrCtxKey string

// Context is a abstraction of http context
type Context interface {
	context.Context

	Header(key, value string)

	JSON(statusCode int, schema interface{})

	Request() *Request

	Response() Response

	String(statusCode int, s string)

	With(ctx context.Context) Context
}

type contextImpl struct {
	req *Request
	res Response
}

func (c *contextImpl) Deadline() (time.Time, bool) {
	return c.Request().Context().Deadline()
}

func (c *contextImpl) Done() <-chan struct{} {
	return c.Request().Context().Done()
}

func (c *contextImpl) Err() error {
	return c.Request().Context().Err()
}

func (c *contextImpl) Value(key interface{}) interface{} {
	return c.Request().Context().Value(key)
}

func (c *contextImpl) Header(key, value string) {
	if c.res.Header().Get(key) != "" {
		zap.L().Warn("unexpected to set same header more than once",
			zap.String("request_id", c.Request().RequestID()),
			zap.String("header", key),
			zap.String("value", value),
		)
		return
	}
	c.res.Header().Set(key, value)
}

func (c *contextImpl) JSON(statusCode int, data interface{}) {
	c.writeContentType("application/json")
	c.res.WriteHeader(statusCode)
	if err := json.NewEncoder(c.res).Encode(data); err != nil {
		Error(c, err.Error())
	}
}

func (c *contextImpl) Request() *Request {
	return c.req
}

func (c *contextImpl) Response() Response {
	return c.res
}

func (c *contextImpl) String(statusCode int, s string) {
	c.writeContentType("text/plain")
	c.res.WriteHeader(statusCode)
	if _, err := fmt.Fprint(c.res, s); err != nil {
		Error(c, err.Error())
	}
}

func (c *contextImpl) With(ctx context.Context) Context {
	c.req.Request = c.Request().WithContext(ctx)
	return c
}

func (c *contextImpl) writeContentType(value string) {
	c.Header("Content-Type", value)
}

func extractRequestID(c context.Context) string {
	reqID, ok := c.Value(StrCtxKey("request_id")).(string)
	if !ok || reqID == "" {
		reqID = "-"
	}
	return reqID
}
