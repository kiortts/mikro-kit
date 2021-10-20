package chiserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/kiortts/mikro-kit/application"
	"github.com/kiortts/mikro-kit/components/httpserver"
)

// ChiServer сервер
type ChiServer struct {
	*application.AbstractComponent
	routers []httpserver.Router // коллекция переданных серверу роутеров
	r       *chi.Mux
}

var _ application.Runnable = (*ChiServer)(nil)
var cfg *Config

func (s *ChiServer) Router() *chi.Mux {
	return s.r
}

// New конструктор
func New(config *Config, routers ...httpserver.Router) *ChiServer {

	checkConfig(config)

	s := &ChiServer{
		AbstractComponent: &application.AbstractComponent{},
		routers:           routers,
	}

	// создание и конфигурация роутера
	s.r = chi.NewRouter() // chi
	for _, api := range s.routers {
		for _, route := range api.Routes() {
			log.Println("Chi http server get route:", route.Name)
			s.r.Method(route.Method, route.Pattern, route.Handler)
		}
	}

	return s
}

// Run запуск сервиса в работу
func (s *ChiServer) Run(main *application.MainParams) error {

	s.MakeLocalCtxAndWg(main)
	localWg := s.Wg

	http.Handle("/", s.r)

	server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port)}
	log.Printf("Chi http server RUN at: %s", server.Addr)

	// запуск сервера
	localWg.Add(1)
	go func() {
		defer localWg.Done()
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Chi http server err: %s", err)
		}
		main.AppStop() // при падении http сервера завершается работа всех остальных сервисов приложения
	}()

	// прекращение работы сервера
	localWg.Add(1)
	go func() {
		defer localWg.Done()
		<-main.Ctx.Done()
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2) // мы даем 2 секунды серверу чтобы самостоятельно прекратить работу, вряд ли потребуется отменить этот контекст быстрее
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Chi http server shutdown err: %v", err)
		}
	}()

	s.WaitAndDo(doExitLog)

	return nil
}

func doExitLog() {
	log.Println("Chi http server DONE")
}
