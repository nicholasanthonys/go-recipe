package router

import (
	"errors"
	"time"

	"github.com/nicholasanthonys/go-recipe/adapter/repository"

	"github.com/nicholasanthonys/go-recipe/adapter/logger"
	"github.com/nicholasanthonys/go-recipe/adapter/validator"
)

type Server interface {
	Listen()
}

type Port int64

var (
	errInvalidWebServerInstance = errors.New("invalid router server instance")
)

const (
	InstanceGorillaMux int = iota
	InstanceGin
)

func NewWebServerFactory(
	instance int,
	log logger.Logger,
	dbSQL repository.SQL,
	dbNoSQL repository.NoSQL,
	kv repository.KeyValStoreIn,
	validator validator.Validator,
	port Port,
	ctxTimeout time.Duration,
) (Server, error) {
	switch instance {
	case InstanceGorillaMux:
		return newGorillaMux(log, dbSQL, validator, port, ctxTimeout), nil
	case InstanceGin:
		return newGinServer(log, dbNoSQL, kv, validator, port, ctxTimeout), nil
	default:
		return nil, errInvalidWebServerInstance
	}
}
