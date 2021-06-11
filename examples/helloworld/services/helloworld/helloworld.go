package helloworld

import (
	"context"
	"fmt"
	"sync"

	"github.com/kiortts/mikro-kit/api"
)

// HelloWordService struct
type HelloWordService struct {
	cancel context.CancelFunc
}

// static interface implementation check
var _ api.Runnable = (*HelloWordService)(nil)

var cfg *Config

// Make the service
func New(config *Config) *HelloWordService {

	checkConfig(config)

	s := &HelloWordService{}
	return s
}

// Run the service
func (s *HelloWordService) Run(mainParams *api.MainParams) error {

	var localCtx context.Context
	localCtx, s.cancel = context.WithCancel(mainParams.Ctx)

	mainParams.Wg.Add(1)
	go printHello(localCtx, mainParams.Wg)

	return nil
}

func printHello(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Hello, %s!!!\n", cfg.Name)
}
