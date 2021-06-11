package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"github.com/kiortts/mikro-kit/api"
	"github.com/pkg/errors"
)

type Application struct {
	name       string
	version    string
	services   []api.Runnable
	mainParams *api.MainParams
}

func New(appName string, appVersion string) *Application {

	app := &Application{
		name:    appName,
		version: appVersion,
	}

	return app
}

func (app *Application) Add(services ...api.Runnable) *Application {
	app.services = append(app.services, services...)
	return app
}

func (app *Application) Run() error {

	log.Printf("Application %s %s run", app.name, app.version)
	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)

	app.mainParams = &api.MainParams{
		Ctx:  ctx,
		Wg:   wg,
		Kill: cancel,
	}

	for _, service := range app.services {
		if err := service.Run(app.mainParams); err != nil {
			app.mainParams.Kill()
			serviceType := reflect.TypeOf(service)
			msg := fmt.Sprintf("Run %s err", serviceType)
			return errors.Wrap(err, msg)
		}
	}

	return nil
}

func (app *Application) Wait() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

func (app *Application) Stop() {
	app.mainParams.Kill()
	app.mainParams.Wg.Wait()
	log.Printf("Application %s %s stop", app.name, app.version)
}
