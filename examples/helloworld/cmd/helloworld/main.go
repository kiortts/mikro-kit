package main

import (
	"log"
	"time"

	"github.com/kiortts/mikro-kit/examples/helloworld/components/helloworld"
	"github.com/kiortts/mikro-kit/service"
	"github.com/pkg/errors"
)

func main() {

	// configure your log
	log.SetFlags(log.Lshortfile)

	// make the service
	serviceName := "Hello-World"
	serviceVersion := "v0.1.0"
	s := service.New(serviceName, serviceVersion)
	makeComponents(s)

	// run the service
	if err := s.Run(); err != nil {
		go s.Stop()
		<-time.After(time.Second * 2)
		log.Fatal(errors.Wrap(err, "service.Run"))
	}

	// waiting for shutdown signal
	s.Wait()

	// shutdown all running components
	s.Stop()
}

// make all service modules and add some of them to service Run
func makeComponents(s *service.Service) {
	hw := helloworld.New(nil)
	s.Add(hw)
}
