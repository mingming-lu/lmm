package cache

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

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

type redisPhotoModel struct {
	ID         uint
	Name       string
	Alternates []string
}

func serialize(c context.Context, photos ...*model.PhotoDescriptor) ([][]byte, error) {
	data := make([][]byte, len(photos))

	for i := range photos {
		buffer := &bytes.Buffer{}
		model := redisPhotoModel{
			ID:         photos[i].ID(),
			Name:       photos[i].Name(),
			Alternates: photos[i].AlternateTexts(),
		}
		if err := gob.NewEncoder(buffer).Encode(model); err != nil {
			return nil, errors.Wrapf(err, "%#v", model)
		}
		data[i] = buffer.Bytes()
	}

	return data, nil
}

func deserialize(c context.Context, data [][]byte) ([]*model.PhotoDescriptor, error) {
	photos := make([]*model.PhotoDescriptor, len(data))

	for i := range data {
		photo := redisPhotoModel{}
		if err := gob.NewDecoder(bytes.NewReader(data[i])).Decode(&photo); err != nil {
			return nil, errors.Wrapf(err, "%s", string(data[i]))
		}
		photos[i] = model.NewPhotoDescriptor(photo.ID, photo.Name)
		for _, alt := range photo.Alternates {
			photos[i].AddAlternateText(alt)
		}
	}

	return photos, nil
}

func (cache *RedisCache) FetchPhotos(c context.Context, page, perPage uint) (*model.PhotoCollection, bool) {
	conn, err := cache.redisClient.GetContext(c)
	if err != nil {
		http.Log().Warn(c, err.Error())
		return nil, false
	}
	defer conn.Close()

	begin := (page - 1) * perPage
	end := begin + perPage
	values, err := redis.Values(conn.Do("ZRANGE", photosListRedisKey, begin, end+1))
	if err != nil {
		http.Log().Warn(c, err.Error())
		return nil, false
	}
	if len(values) == 0 {
		return nil, false
	}

	hasNextPage := false
	if len(values) > int(perPage) {
		values = values[:int(perPage)]
		hasNextPage = true
	}

	models := make([][]byte, len(values))
	if err := redis.ScanSlice(values, &models); err != nil {
		http.Log().Warn(c, err.Error())
		return nil, false
	}

	photos, err := deserialize(c, models)
	if err != nil {
		http.Log().Warn(c, err.Error())
		return nil, false
	}

	return model.NewPhotoCollection(photos, hasNextPage), true
}

func (cache *RedisCache) StorePhotos(c context.Context, page, perPage uint, photos []*model.PhotoDescriptor) error {
	data, err := serialize(c, photos...)
	if err != nil {
		return err
	}
	conn, err := cache.redisClient.GetContext(c)
	if err != nil {
		return err
	}

	begin := (page - 1) * perPage

	args := make([]interface{}, len(photos)*2)

	for i, photo := range data {
		baseIdx := 2 * i
		args[baseIdx] = begin + uint(i)
		args[baseIdx+1] = photo
	}
	args = append([]interface{}{photosListRedisKey}, args...)

	if _, err := conn.Do("ZADD", args...); err != nil {
		return err
	}
	return nil
}

func (cache *RedisCache) ClearPhotos(c context.Context) error {
	conn, err := cache.redisClient.GetContext(c)
	if err != nil {
		return err
	}

	if _, err := conn.Do("DEL", photosListRedisKey); err != nil {
		return err
	}
	return nil
}
