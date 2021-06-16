package components

import (
	"context"
	"sync"
)

type Runnable interface {
	Run(*MainParams) error
}

type MainParams struct {
	Ctx  context.Context
	Wg   *sync.WaitGroup
	Kill func()
}
