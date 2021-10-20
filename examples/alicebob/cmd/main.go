package main

import (
	"log"
	"time"

	"github.com/kiortts/mikro-kit/application"
	"github.com/kiortts/mikro-kit/components/messages/bigmessageservice/minioclient"
	"github.com/kiortts/mikro-kit/components/messages/messageservice/natsclient"
	"github.com/kiortts/mikro-kit/examples/alicebob/components/alice"
	"github.com/kiortts/mikro-kit/examples/alicebob/components/bob"
	"github.com/kiortts/mikro-kit/examples/alicebob/components/eva"
	"github.com/pkg/errors"
)

var AppName = "Alice-Bob"
var AppVersion = "v0.1.0"

func main() {

	// configure your logger
	log.SetFlags(log.Lshortfile)

	// make the app

	app := application.New(AppName, AppVersion)
	makeApplicationComponents(app)

	// run the service and handle run errors
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

	// big messages service
	bmsCfg := &minioclient.Config{
		Endpoint:        "192.168.31.25:9000",
		AccessKeyID:     "minio",
		SecretAccessKey: "12345678",
		BucketName:      "bigmsgs",
	}
	bms := minioclient.New(bmsCfg)
	app.Add(bms)

	// messages service
	ms := natsclient.New(nil, bms)
	app.Add(ms)

	// Alice, tell something to bob or ask him
	Alice := alice.New(ms)
	app.Add(Alice)

	// Bob, listen to Alice and answer her
	Bob := bob.New(ms)
	app.Add(Bob)

	// Eva, daemon, listen to both and make some service
	Eva := eva.New(ms)
	app.Add(Eva)
}
