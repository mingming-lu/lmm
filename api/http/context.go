package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

// ContextKey is used to indentify context value key
type ContextKey int32

// Context is a abstraction of http context
type Context interface {
	context.Context

	Header(key, value string)

	JSON(statusCode int, schema interface{})

	KeyRegistry(name string) ContextKey

	Request() *Request

	String(statusCode int, s string)

	With(ctx context.Context) Context
}

type contextImpl struct {
	keyCount int32
	keyMap   map[string]ContextKey
	req      *Request
	res      Response
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
		log.Printf("unexpected to set same header more than once, header = %s, value = %s\n", key, value)
		return
	}
	c.res.Header().Set(key, value)
}

func (c *contextImpl) JSON(statusCode int, data interface{}) {
	c.writeContentType("application/json")
	c.res.WriteHeader(statusCode)
	if err := json.NewEncoder(c.res).Encode(data); err != nil {
		panic(err)
	}
}

func (c *contextImpl) KeyRegistry(name string) ContextKey {
	if key, ok := c.keyMap[name]; ok {
		return ContextKey(key)
	}

	newKey := ContextKey(atomic.AddInt32(&c.keyCount, 1))
	c.keyMap[name] = newKey
	return newKey
}

func (c *contextImpl) Request() *Request {
	return c.req
}

func (c *contextImpl) String(statusCode int, s string) {
	c.writeContentType("text/plain")
	c.res.WriteHeader(statusCode)
	if _, err := fmt.Fprint(c.res, s); err != nil {
		panic(err)
	}
}

func (c *contextImpl) With(ctx context.Context) Context {
	c.req.Request = c.Request().WithContext(ctx)
	return c
}

func (c *contextImpl) writeContentType(value string) {
	c.Header("Content-Type", value)
}
