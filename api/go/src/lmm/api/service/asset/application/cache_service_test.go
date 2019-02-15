package application

import (
	"context"
	"lmm/api/testing"
)

func TestNopCacheService(tt *testing.T) {
	c := context.Background()

	tt.Run("Fetch", func(tt *testing.T) {
		t := testing.NewTester(tt)

		photos, hit := NopCacheService().FetchPhotos(c, 0, 0)
		t.False(hit)
		t.Nil(photos)
	})

	tt.Run("Store", func(tt *testing.T) {
		t := testing.NewTester(tt)

		t.NoError(NopCacheService().StorePhotos(c, 0, 0, nil))
	})

	tt.Run("Clear", func(tt *testing.T) {
		t := testing.NewTester(tt)

		t.NoError(NopCacheService().ClearPhotos(c))
	})
}
