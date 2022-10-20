package database

import (
	"errors"

	"github.com/nicholasanthonys/go-recipe/adapter/repository"
)

var (
	errInvalidKeyValStoreDatabaseInstance = errors.New("invalid keyvalstore db instance")
)

const (
	InstanceRedis int = iota
)

func NewDatabaseKeyValStoreFactory(instance int) (repository.KeyValStore, error) {
	switch instance {
	case InstanceRedis:
		return NewRedistHanlder(newConfigRedis())
	default:
		return nil, errInvalidKeyValStoreDatabaseInstance
	}
}
