package main

import (
	"log"
	"time"

	"github.com/kiortts/mikro-kit/application"
	"github.com/kiortts/mikro-kit/components/httpserver/chiserver"
	"github.com/kiortts/mikro-kit/examples/ghibli/components/endpoints/films"
	"github.com/kiortts/mikro-kit/examples/ghibli/components/endpoints/locations"
	"github.com/kiortts/mikro-kit/examples/ghibli/components/endpoints/people"
	"github.com/kiortts/mikro-kit/examples/ghibli/components/storage"
	"github.com/pkg/errors"
)

func main() {

	// configure your log
	log.SetFlags(log.Lshortfile)

	// make the app
	appName := "GhibliAPI"
	appVersion := "v0.1.0"
	app := application.New(appName, appVersion)
	makeApplicationComponents(app)

	// run the service
	if err := app.Run(); err != nil {
		go app.Stop()
		<-time.After(time.Second * 2)
		log.Fatal(errors.Wrap(err, "app.Run"))
	}

	// waiting for shutdown signal
	app.Wait()

	// shutdown all running components
	app.Stop()
}

// make all app components and add some of them to app.Run
func makeApplicationComponents(app *application.Application) {
	stor := storage.NewLocal()
	app.Add(stor)

	filmsEP := films.New(stor)
	peopleEP := people.New(stor)
	locationsEP := locations.New(stor)
	// server := gorillaserver.New(nil, filmsEP, peopleEP, locationsEP)
	server := chiserver.New(nil, filmsEP, peopleEP, locationsEP)
	app.Add(server)
}
