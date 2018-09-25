package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// Context is a abstraction of http context
type Context interface {
	context.Context

	Header(key, value string)

	JSON(statusCode int, schema interface{})

	Request() Request

	String(statusCode int, s string)
}

type contextImpl struct {
	context.Context
	r Request
	w Response
}

func (c *contextImpl) Header(key, value string) {
	if c.w.Header().Get(key) != "" {
		log.Printf("unexpected to set same header more than once, header = %s, value = %s\n", key, value)
		return
	}
	c.w.Header().Set(key, value)
}

func (c *contextImpl) JSON(statusCode int, data interface{}) {
	c.writeContentType("application/json")
	c.w.WriteHeader(statusCode)
	if err := json.NewEncoder(c.w).Encode(data); err != nil {
		panic(err)
	}
}

func (c *contextImpl) Request() Request {
	return c.r
}

func (c *contextImpl) String(statusCode int, s string) {
	c.writeContentType("text/plain")
	c.w.WriteHeader(statusCode)
	if _, err := fmt.Fprint(c.w, s); err != nil {
		panic(err)
	}
}

func (c *contextImpl) writeContentType(value string) {
	c.Header("Content-Type", value)
}
