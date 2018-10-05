package http

import (
	"time"

	"lmm/api/util/bashutil"

	"go.uber.org/zap"
)

var reloader = NewReloader("lmm/api")

type Reloader struct {
	mainPath   string
	lastHash   string
	lastModify time.Time
}

func NewReloader(mainDir string) *Reloader {
	r := &Reloader{mainPath: mainDir}
	r.Recompile()
	r.CompareMD5AndSwap()
	return r
}

func (r *Reloader) Recompile() {
	if out, err := bashutil.SyncRun("go install " + r.mainPath); err != nil {
		zap.L().Error(err.Error(), zap.String("result", out))
	}
}

func (r *Reloader) CompareMD5AndSwap() {
	out, err := bashutil.SyncRun("cat $(which api) | md5sum")
	if err != nil {
		zap.L().Error(err.Error(), zap.String("result", out))
		return
	}
	hash := out[:32]
	if r.lastHash != hash {
		r.lastHash = hash
		r.lastModify = time.Now()
	}
}
