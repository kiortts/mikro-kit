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

var AppName = "GhibliAPI"
var AppVersion = "v0.1.0"

func main() {

	// configure your logger
	log.SetFlags(log.Lshortfile)

	// make the app
	app := application.New(AppName, AppVersion)
	makeApplicationComponents(app)

	// run the service
	if err := app.Run(); err != nil {
		app.Stop(time.Second * 2)
		log.Fatal(errors.Wrap(err, "app.Run"))
	}

	// waiting for shutdown signal
	app.Wait()

	// shutdown all running components
	app.Stop(time.Second * 5)
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
