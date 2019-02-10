package contextutil

import "context"

type contextKey int

const (
	requestIDKey contextKey = iota
)

func WithRequestID(c context.Context, reqID string) context.Context {
	return context.WithValue(c, requestIDKey, reqID)
}

func RequestID(c context.Context) string {
	reqID, ok := c.Value(requestIDKey).(string)
	if reqID == "" || !ok {
		reqID = "-"
	}
	return reqID
}
