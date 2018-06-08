package factory

import (
	"time"

	"github.com/sony/sonyflake"
)

var factory Factory

func init() {
	st := sonyflake.Settings{
		StartTime: time.Now(),
	}
	factory = factoryT{idGenerator: sonyflake.NewSonyflake(st)}
}

type Factory interface {
	GenerateID() (uint64, error)
}

func Default() Factory {
	return factory
}

type factoryT struct {
	idGenerator *sonyflake.Sonyflake
}

func (f factoryT) GenerateID() (uint64, error) {
	return f.idGenerator.NextID()
}
