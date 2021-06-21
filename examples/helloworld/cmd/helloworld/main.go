package main

import (
	"log"
	"time"

	"github.com/kiortts/mikro-kit/application"
	"github.com/kiortts/mikro-kit/examples/helloworld/components/helloworld"
	"github.com/pkg/errors"
)

func main() {

	// configure your log
	log.SetFlags(log.Lshortfile)

	// make the app
	appName := "Hello-World"
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
	hw := helloworld.New(nil)
	app.Add(hw)
}
