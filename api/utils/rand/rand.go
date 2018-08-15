package rand

import (
	"math/rand"
	"time"
)

var (
	base64      = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_")
	defaultRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func Seed(seed int64) {
	defaultRand.Seed(seed)
}

func Base64(n int) string {
	return string(base(defaultRand, n, 64))
}

func Base62(n int) string {
	return string(base(defaultRand, n, 62))
}

func Base36(n int) string {
	return string(base(defaultRand, n, 36))
}

func base(r *rand.Rand, n, base int) []byte {
	b := make([]byte, n, n)
	for i := 0; i < n; i++ {
		b[i] = base64[r.Intn(base)]
	}
	return b
}

type Rand struct {
	*rand.Rand
}

func New(seed int64) *Rand {
	internal := rand.New(rand.NewSource(seed))
	r := Rand{Rand: internal}
	return &r
}

func (r *Rand) Base64(n int) string {
	return string(base(r.Rand, n, 64))
}

func (r *Rand) Base62(n int) string {
	return string(base(r.Rand, n, 62))
}

func (r *Rand) Base36(n int) string {
	return string(base(r.Rand, n, 36))
}
