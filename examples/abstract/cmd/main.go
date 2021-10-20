package main

import (
	"log"
	"time"

	"github.com/kiortts/mikro-kit/application"
	componenta "github.com/kiortts/mikro-kit/examples/abstract/components/componentA"
	"github.com/pkg/errors"
)

var AppName = "AbstractApplication" // FIXME: change to <your service> name
var AppVersion = "dev"              // FIXME: change to current version

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

	// TODO: implements me like:

	caConfig := &componenta.Config{
		Param: "Value",
	}
	ca := componenta.New(caConfig)
	app.Add(ca)
}
