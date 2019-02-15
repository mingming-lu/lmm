package cache

import (
	"context"

	"github.com/gomodule/redigo/redis"

	"lmm/api/http"
	"lmm/api/service/asset/domain/model"
	"lmm/api/storage/cache"
)

const photosListRedisKey = "assets:photos"

type RedisCache struct {
	redisClient *cache.RedisClient
}

func NewRedisCache() *RedisCache {
	client, err := cache.NewRedisClient()
	if err != nil {
		panic(err)
	}
	return &RedisCache{
		redisClient: client,
	}
}

func (cache *RedisCache) FetchPhotos(c context.Context, page, perPage uint) (*model.PhotoCollection, bool) {
	conn, err := cache.redisClient.GetContext(c)
	if err != nil {
		http.Log().Warn(c, err.Error())
		return nil, false
	}
	defer conn.Close()

	photos := make([]*model.PhotoDescriptor, 0)

	begin := (page - 1) * perPage
	end := begin + perPage
	values, err := redis.Values(conn.Do("ZRANGE", photosListRedisKey, begin, end+1))
	if err != nil {
		http.Log().Warn(c, err.Error())
		return nil, false
	}

	hasNextPage := false
	if len(values) > int(perPage) {
		values = values[:int(perPage)]
		hasNextPage = true
	}

	if err := redis.ScanSlice(values, &photos); err != nil {
		http.Log().Warn(c, err.Error())
		return nil, false
	}

	return model.NewPhotoCollection(photos, hasNextPage), true
}

func (cache *RedisCache) StorePhotos(c context.Context, page, perPage uint, photos []*model.PhotoDescriptor) error {
	conn, err := cache.redisClient.GetContext(c)
	if err != nil {
		return err
	}

	begin := (page - 1) * perPage

	args := make([]interface{}, len(photos)*2)

	for i, photo := range photos {
		baseIdx := 2 * i
		args[baseIdx] = begin + uint(i)
		args[baseIdx+1] = photo
	}

	if _, err := conn.Do("ZADD", photosListRedisKey, args); err != nil {
		return err
	}
	return nil
}
