package repository

import (
	"context"
	"time"
)

type KeyValStore interface {
	Set(context.Context, string, string, time.Duration) error
	Get(context.Context, string) (string, error)
}
