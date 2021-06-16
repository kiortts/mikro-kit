package service

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

type Service struct {
	name       string
	version    string
	components []components.Runnable
	mainParams *components.MainParams
}

func New(serviceName string, serviceVersion string) *Service {

	s := &Service{
		name:    serviceName,
		version: serviceVersion,
	}

	return s
}

func (s *Service) Add(services ...components.Runnable) *Service {
	s.components = append(s.components, services...)
	return s
}

func (s *Service) Run() error {

	log.Printf("Service %s %s run", s.name, s.version)
	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)

	s.mainParams = &components.MainParams{
		Ctx:  ctx,
		Wg:   wg,
		Kill: cancel,
	}

	for _, service := range s.components {
		if err := service.Run(s.mainParams); err != nil {
			s.mainParams.Kill()
			serviceType := reflect.TypeOf(service)
			msg := fmt.Sprintf("Run %s err", serviceType)
			return errors.Wrap(err, msg)
		}
	}

	return nil
}

func (s *Service) Wait() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

func (s *Service) Stop() {
	s.mainParams.Kill()
	s.mainParams.Wg.Wait()
	log.Printf("Service %s %s stop", s.name, s.version)
}
