package application

import (
	"context"
	"sync"
)

type Runnable interface {
	Run(*MainParams) error
	Stop()
}

type MainParams struct {
	Ctx     context.Context
	Wg      *sync.WaitGroup
	AppStop func()
}
