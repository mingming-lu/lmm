package cli

import (
	"context"
	"lmm/api/testing"
	"sync"
	"sync/atomic"
)

type counter struct {
	count uint64
}

func (c *counter) Execute(_ context.Context) error {
	atomic.AddUint64(&c.count, 1)
	return nil
}

func TestCommandImpl(tt *testing.T) {
	t := testing.NewTester(tt)
	c := context.Background()

	cmd := &counter{}

	t.Is(uint64(0), cmd.count)

	t.NoError(cmd.Execute(c))
	t.Is(uint64(1), cmd.count)

	const times uint64 = 10
	wg := sync.WaitGroup{}

	for i := uint64(0); i < times; i++ {
		wg.Add(1)
		go func(command Command) {
			defer wg.Done()
			command.Execute(c)
		}(cmd)
	}

	wg.Wait()

	t.Is(uint64(11), cmd.count)
}

func TestRegisterCommand(tt *testing.T) {
	t := testing.NewTester(tt)

	cmd := &counter{}

	Register("count1", cmd)

	c, ok := commands["count1"]
	t.True(ok)
	t.Is(cmd, c)
}

func TestExecuteCommand(tt *testing.T) {
	t := testing.NewTester(tt)

	cmd := &counter{}
	Register("count2", cmd)
	Execute(context.Background(), "count2")

	t.Is(uint64(1), cmd.count)
}
