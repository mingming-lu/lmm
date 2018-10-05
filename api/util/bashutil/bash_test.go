package bashutil

import (
	"lmm/api/testing"
	"runtime"
)

func TestAsynRun(tt *testing.T) {
	t := testing.NewTester(tt)
	res, err := SyncRun("go version")
	t.NoError(err)
	t.Regexp(runtime.Version(), res)
}
