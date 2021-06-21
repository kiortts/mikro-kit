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

	"github.com/kiortts/mikro-kit/components"
	"github.com/pkg/errors"
)

type Application struct {
	name       string
	version    string
	components []components.Runnable
	mainParams *components.MainParams
}

func New(appName string, appVersion string) *Application {

	s := &Application{
		name:    appName,
		version: appVersion,
	}

	return s
}

func (s *Application) Add(components ...components.Runnable) *Application {
	s.components = append(s.components, components...)
	return s
}

func (s *Application) Run() error {

	log.Printf("Application %s %s run", s.name, s.version)
	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)

	s.mainParams = &components.MainParams{
		Ctx:  ctx,
		Wg:   wg,
		Kill: cancel,
	}

	for _, component := range s.components {
		if err := component.Run(s.mainParams); err != nil {
			s.mainParams.Kill()
			componentType := reflect.TypeOf(component)
			msg := fmt.Sprintf("Run %s err", componentType)
			return errors.Wrap(err, msg)
		}
	}

	return nil
}

func (s *Application) Wait() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

func (s *Application) Stop() {
	s.mainParams.Kill()
	s.mainParams.Wg.Wait()
	log.Printf("Application %s %s stop", s.name, s.version)
}
