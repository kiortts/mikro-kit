package helloworld

import (
	"context"
	"fmt"
	"sync"

	"github.com/kiortts/mikro-kit/application"
)

// HelloWordModule struct
type HelloWordModule struct {
	cancel context.CancelFunc
}

// static interface implementation check
var _ application.Runnable = (*HelloWordModule)(nil)

// static variables
var cfg *Config

// Make the instance.
// This method never returns any error.
func New(config *Config) *HelloWordModule {

	checkConfig(config)

	s := &HelloWordModule{}
	return s
}

func (s *HelloWordModule) Stop() {}

// Run the module.
// This is non blocking method returns only module starting errors.
func (s *HelloWordModule) Run(mainParams *application.MainParams) error {

	// make the local context for this module instance
	var localCtx context.Context
	localCtx, s.cancel = context.WithCancel(mainParams.Ctx)

	// run module workers
	mainParams.Wg.Add(1)                             // every time increment the WaitGroup before start goroutine
	go printHello(localCtx, mainParams.Wg, cfg.Name) // run goroutine

	return nil
}

// Some runtime function.
func printHello(ctx context.Context, wg *sync.WaitGroup, name string) {
	defer wg.Done()                    // every time Done the WaitGroup before leave goroutine
	fmt.Printf("Hello, %s!!!\n", name) // do something
}
