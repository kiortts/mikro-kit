package natsclient

import (
	"fmt"
	"log"
	"time"

	"github.com/kiortts/mikro-kit/components/messages"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// обертка для *nats.Msg, реализующая интерфейса Message
type NatsMsg struct {
	*nats.Msg
}

// статическая проверка реализаии интерфеса
var _ messages.Message = (*NatsMsg)(nil)

// получение данных из сообщения, реализация интерфейса Message
func (m *NatsMsg) Data() []byte {
	return m.Msg.Data
}

func (m *NatsMsg) SetData(data []byte) {
	m.Msg.Data = data
}

// @ Overwrite
func (m *NatsMsg) Respond(data []byte) error {

	if bmService == nil || len(data) < messages.BigMessageSize {
		return m.Msg.Respond(data)
	}

	// пересылка через модуль больших сообщений
	log.Printf("Giant respond message: %d bytes", len(data))
	subj := uuid.Must(uuid.NewV4(), nil).String() // генерация ключа
	proxyData, err := bmService.Put([]byte(subj), data)
	if err != nil {
		log.Println(errors.Wrap(err, "bigMsgService.Put"))
		return m.Msg.Respond(data)
	}

	// запуск отложенного процедуры удаления сообщения из хранилища больших сообщений
	go func() {
		<-time.After(time.Second * 10)
		bmService.Remove([]byte(subj)) // FIXME: BigMessageBucket
	}()

	return m.Msg.Respond(proxyData)
}

// обертка для *StanMsg, реализующая интерфейса Message
type StanMsg struct {
	*stan.Msg
}

// статическая проверка реализаии интерфеса
var _ messages.Message = (*StanMsg)(nil)

// получение данных из сообщения, реализация интерфейса Message
func (m *StanMsg) Data() []byte {
	return m.Msg.Data
}

func (m *StanMsg) SetData(data []byte) {
	m.Msg.Data = data
}

// Отправка ответа, реализация интерфейса Message.
// Т.к. в NATS-Streaming отправка сообщений с запросом ответа не поддерживается, метод возвращает ошибку.
func (m *StanMsg) Respond(data []byte) error {
	return fmt.Errorf("NATS-Streaming don`t implement Respond method")
}
