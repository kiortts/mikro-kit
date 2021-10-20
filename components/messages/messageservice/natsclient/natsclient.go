/*
Сервис сообщений.
*/

package natsclient

import (
	"context"
	"log"

	"github.com/kiortts/mikro-kit/application"
	"github.com/kiortts/mikro-kit/components/messages"
	"github.com/nats-io/stan.go"
)

// Клиент для отправки сообщений через NATS и NATS-Streaming
type NatsClient struct {
	cancel context.CancelFunc
	nsc    stan.Conn // соединение NATS-Streaming
}

// статическая проверка реализаии интерфесов
var _ messages.Service = (*NatsClient)(nil)
var _ application.Runnable = (*NatsClient)(nil)

var cfg *Config
var bmService messages.BigMsgService

// передач сервису соединения с NATS
func (s *NatsClient) SetConn(nsc stan.Conn) *NatsClient {
	s.nsc = nsc
	return s
}

// Соединение в данный момент установлено
func (s *NatsClient) IsConnected() bool {
	if s.nsc == nil {
		return false
	}
	return s.nsc.NatsConn().IsConnected()
}

// Конструктор NatsClient.
// В случае передачи nil в качестве конфигурации, конфигурация будет получена из переменных окружения.
func New(config *Config, bigMsgService messages.BigMsgService) *NatsClient {
	checkConfig(config)
	bmService = bigMsgService
	s := &NatsClient{}
	return s
}

// Конструктор NatsClient c произвольными параметрами.
func UnsafeNew(params ...interface{}) *NatsClient {

	var config *Config
	var bigMsgService messages.BigMsgService

	for _, param := range params {
		switch p := param.(type) {
		case *Config:
			config = p

		case messages.BigMsgService:
			bigMsgService = p
		}
	}

	// вызов конструктора с типизированными парметрами
	return New(config, bigMsgService)
}

// Завершение работы модуля.
// Реализация интерфейса service.Runnable.
func (s *NatsClient) Stop() {
	s.cancel()
	s.nsc.Close()
	s.nsc.NatsConn().Close()
}

// Запуск клиента в работу.
// Реализация интерфейса application.Runnable.
func (s *NatsClient) Run(mainParams *application.MainParams) error {

	// неблокирующие методы
	s.getConnection(mainParams.Ctx, mainParams.Wg, cfg)
	s.runCloser(mainParams.Ctx, mainParams.Wg)

	log.Printf("NatsClient Run")
	return nil
}
