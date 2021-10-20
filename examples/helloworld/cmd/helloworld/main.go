package main

import (
	"log"
	"time"

	"github.com/kiortts/mikro-kit/application"
	"github.com/kiortts/mikro-kit/examples/helloworld/components/helloworld"
	"github.com/pkg/errors"
)

var AppName = "Hello-World"
var AppVersion = "dev"

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
	hw := helloworld.New(nil)
	app.Add(hw)
}
