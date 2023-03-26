package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
)

var KVNotFound = redis.Nil

type KeyValStoreIn interface {
	Set(context.Context, string, string, time.Duration) error
	Get(context.Context, string) (string, error)
	Del(ctx context.Context, key []string) (int64, error)
}

type KeyValStoreImpl struct {
	db KeyValStoreIn
}

func NewKeyValStore(db KeyValStoreIn) KeyValStoreImpl {
	return KeyValStoreImpl{
		db: db,
	}
}

func (r KeyValStoreImpl) Set(ctx context.Context, key string, val string, time time.Duration) error {
	err := r.db.Set(ctx, key, val, time)
	return err
}

func (r KeyValStoreImpl) Get(ctx context.Context, key string) (string, error) {
	res, err := r.db.Get(ctx, key)
	return res, err
}
func (r KeyValStoreImpl) Del(ctx context.Context, key []string) (int64, error) {
	res, err := r.db.Del(ctx, key)
	return res, err
}
