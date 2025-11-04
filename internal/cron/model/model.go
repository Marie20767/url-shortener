package model

import "context"

type CronLike interface {
	Add(ctx context.Context) error
	Start()
	Stop() context.Context
}