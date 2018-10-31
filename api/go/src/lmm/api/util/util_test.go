package util

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"lmm/api/testing"
)

func TestRetry(t *testing.T) {
	t.Run("NO_ERROR", func(tt *testing.T) {
		t := testing.NewTester(tt)
		var count int
		err := Retry(3, func() error {
			count++
			return nil
		})
		t.Nil(err)
		t.Is(1, count)
	})

	t.Run("NOT_RETRY", func(tt *testing.T) {
		t := testing.NewTester(tt)
		var count int
		err := Retry(0, func() error {
			count++
			return errors.New("error")
		})
		t.Is(1, count)
		t.NotNil(err)
		t.Is("error", err.Error())
	})

	t.Run("FINITE_THREE_TIMES", func(tt *testing.T) {
		t := testing.NewTester(tt)
		var count int
		err := Retry(3, func() error {
			count++
			return errors.New("error")
		})
		t.Is(4, count)
		t.NotNil(err)
		t.Is("error", err.Error())
	})

	t.Run("INFINITE", func(tt *testing.T) {
		t := testing.NewTester(tt)

		c, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
		defer cancel()

		go Retry(-1, func() error {
			return errors.New("error")
		})

		<-c.Done()
		t.Is(context.DeadlineExceeded, c.Err())
	})
}
