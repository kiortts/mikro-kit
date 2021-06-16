package main

import (
	"log"
	"time"

	"github.com/kiortts/mikro-kit/components/httpserver/chiserver"
	"github.com/kiortts/mikro-kit/examples/ghibli/components/endpoints/films"
	"github.com/kiortts/mikro-kit/examples/ghibli/components/endpoints/locations"
	"github.com/kiortts/mikro-kit/examples/ghibli/components/endpoints/people"
	"github.com/kiortts/mikro-kit/examples/ghibli/components/storage"
	"github.com/kiortts/mikro-kit/service"
	"github.com/pkg/errors"
)

func main() {

	// configure your log
	log.SetFlags(log.Lshortfile)

	// make the service
	serviceName := "GhibliAPI"
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
	stor := storage.NewLocal()
	s.Add(stor)

	filmsEP := films.New(stor)
	peopleEP := people.New(stor)
	locationsEP := locations.New(stor)
	// server := gorillaserver.New(nil, filmsEP, peopleEP, locationsEP)
	server := chiserver.New(nil, filmsEP, peopleEP, locationsEP)
	s.Add(server)
}
