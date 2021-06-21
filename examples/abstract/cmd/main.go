package main

import (
	"log"
	"time"

	"github.com/kiortts/mikro-kit/application"
	componenta "github.com/kiortts/mikro-kit/examples/abstract/components/componentA"
	"github.com/pkg/errors"
)

func main() {

	// configure your log
	log.SetFlags(log.Lshortfile)

	// make the app
	appName := "AbstractApplication" // FIXME: change to <your service> name
	appVersion := "v0.1.0"           // FIXME: change to current version
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

	// TODO: implements me like:

	caConfig := &componenta.Config{
		Param: "Value",
	}
	ca := componenta.New(caConfig)
	app.Add(ca)
}
