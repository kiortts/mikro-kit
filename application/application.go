package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

type Application struct {
	name       string
	version    string
	components []Runnable
	mainParams *MainParams
}

func New(appName string, appVersion string) *Application {

	s := &Application{
		name:    appName,
		version: appVersion,
	}

	return s
}

func (s *Application) Add(components ...Runnable) *Application {
	s.components = append(s.components, components...)
	return s
}

func (s *Application) Run() error {

	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)

	s.mainParams = &MainParams{
		Ctx:     ctx,
		Wg:      wg,
		AppStop: cancel,
	}

	for _, component := range s.components {
		if err := component.Run(s.mainParams); err != nil {
			s.mainParams.AppStop()
			componentType := reflect.TypeOf(component)
			msg := fmt.Sprintf("Run %s err", componentType)
			return errors.Wrap(err, msg)
		}
	}

	log.Printf("Application %s %s RUN", s.name, s.version)
	return nil
}

func (s *Application) Wait() {

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-quit:
	case <-s.mainParams.Ctx.Done():
	}
}

func (s *Application) Stop(timeout time.Duration) {

	s.mainParams.AppStop()
	for _, component := range s.components {
		component.Stop()
	}

	if waitTimeout(s.mainParams.Wg, timeout) {
		log.Println("Timed out waiting for main wait group")
	} else {
		log.Println("Main wait group finished")
	}

	log.Printf("Application %s %s DONE", s.name, s.version)
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
