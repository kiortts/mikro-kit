package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kiortts/mikro-kit/application"
)

// HttpServer сервер
type HttpServer struct {
	routers []Router // коллекция переданных серверу API
	r       *mux.Router
}

var _ application.Runnable = (*HttpServer)(nil)
var cfg *Config

func (s *HttpServer) Router() *mux.Router {
	return s.r
}

// New конструктор
func New(config *Config, routers ...Router) *HttpServer {

	checkConfig(config)

	s := &HttpServer{routers: routers}

	// создание и конфигурация роутера
	// s.r = chi.NewRouter() // chi
	// for _, api := range s.apis {
	// 	for _, route := range api.Routes() {
	// 		log.Println(route.Name)
	// 		s.r.Method(route.Method, route.Pattern, route.Handler)
	// 	}
	// }

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

// Run запуск сервиса в работу
func (s *HttpServer) Run(main *application.MainParams) error {

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
