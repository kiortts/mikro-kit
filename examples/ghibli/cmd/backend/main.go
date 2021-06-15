package main

import (
	"log"
	"time"

	"github.com/kiortts/mikro-kit/application"
	"github.com/kiortts/mikro-kit/examples/ghibli/services/endpoints/films"
	"github.com/kiortts/mikro-kit/examples/ghibli/services/endpoints/locations"
	"github.com/kiortts/mikro-kit/examples/ghibli/services/endpoints/people"
	"github.com/kiortts/mikro-kit/examples/ghibli/services/storage"
	"github.com/kiortts/mikro-kit/services/httpserver/chiserver"
	"github.com/pkg/errors"
)

func main() {

	// configure your log
	log.SetFlags(log.Lshortfile)

	// make and build the app
	appName := "GhibliAPI"
	appVersion := "v0.1.0"
	app := application.New(appName, appVersion)
	buildApp(app)

	// run the app
	if err := app.Run(); err != nil {
		go app.Stop()
		<-time.After(time.Second * 2)
		log.Fatal(errors.Wrap(err, "app.Run"))
	}

	// waiting for shutdown signal
	app.Wait()

	// shutdown all run modules
	app.Stop()
}

// make all app modules and add some of them to app Run
func buildApp(app *application.Application) {
	stor := storage.NewLocal()
	app.Add(stor)

	filmsEP := films.New(stor)
	peopleEP := people.New(stor)
	locationsEP := locations.New(stor)
	// server := gorillaserver.New(nil, filmsEP, peopleEP, locationsEP)
	server := chiserver.New(nil, filmsEP, peopleEP, locationsEP)
	app.Add(server)
}
