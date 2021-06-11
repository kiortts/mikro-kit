package main

import (
	"log"
	"time"

	"github.com/kiortts/mikro-kit/application"
	"github.com/kiortts/mikro-kit/examples/helloworld/services/helloworld"
	"github.com/pkg/errors"
)

func main() {

	// configure your log
	log.SetFlags(log.Lshortfile)

	// make and build the app
	appName := "Hello-World"
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
	hw := helloworld.New(nil)
	app.Add(hw)
}
