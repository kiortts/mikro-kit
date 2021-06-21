package gorillaserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kiortts/mikro-kit/components"
	"github.com/kiortts/mikro-kit/components/httpserver"
)

// GorillaServer сервер
type GorillaServer struct {
	routers []httpserver.Router // коллекция переданных серверу API
	r       *mux.Router         // роутер
}

//
var _ components.Runnable = (*GorillaServer)(nil)
var cfg *Config

func (s *GorillaServer) Router() *mux.Router {
	return s.r
}

// New конструктор
func New(config *Config, routers ...httpserver.Router) *GorillaServer {

	checkConfig(config)

	s := &GorillaServer{routers: routers}

	s.r = mux.NewRouter() // gorilla
	for _, router := range s.routers {
		for _, route := range router.Routes() {
			s.r.Methods(route.Method).
				Path(route.Pattern).
				Queries(route.QueryPairs...).
				Name(route.Name).
				Handler(route.Handler)
		}
	}

	return s
}

// пустой хэндлер
func dummyHandler(w http.ResponseWriter, r *http.Request) {}

func (s *GorillaServer) Stop() {
	// TODO: cancel local context
}

// Run запуск сервиса в работу
func (s *GorillaServer) Run(main *components.MainParams) error {

	http.Handle("/", s.r)

	server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port)}
	log.Printf("Start http server at: %s", server.Addr)

	// запуск сервера
	main.Wg.Add(1)
	go func() {
		defer main.Wg.Done()
		if err := server.ListenAndServe(); err != nil {
			log.Printf("http server DONE: %v", err)
		}
		main.Kill() // при падении http сервера завершается работа всех остальных сервисов приложения
	}()

	// прекращение работы сервера
	main.Wg.Add(1)
	go func() {
		defer main.Wg.Done()
		<-main.Ctx.Done()
		ctx, _ := context.WithTimeout(context.TODO(), time.Second*2) // мы даем 2 секунды серверу чтобы самостоятельно прекратить работу, вряд ли потребуется отменить этот контекст быстрее
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown err: %v", err)
		}
	}()

	return nil
}
