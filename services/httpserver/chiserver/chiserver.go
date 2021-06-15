package chiserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/kiortts/mikro-kit/application"
	"github.com/kiortts/mikro-kit/services/httpserver"
)

// ChiServer сервер
type ChiServer struct {
	routers []httpserver.Router // коллекция переданных серверу API
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

	s := &ChiServer{routers: routers}

	// создание и конфигурация роутера
	s.r = chi.NewRouter() // chi
	for _, api := range s.routers {
		for _, route := range api.Routes() {
			log.Println(route.Name)
			s.r.Method(route.Method, route.Pattern, route.Handler)
		}
	}

	return s
}

// Run запуск сервиса в работу
func (s *ChiServer) Run(main *application.MainParams) error {

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
