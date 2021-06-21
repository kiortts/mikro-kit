package messages

import (
	"context"
	"sync"
	"time"
)

const (
	BigMessageBucket = "bigmsg"
	BigMessageSize   = 1024 * 128
)

// Интерфейс сообщения
type Message interface {
	Data() []byte              // получение данных из сообщения
	SetData(data []byte)       // метод лоя подмены нагрузки сообщения
	Respond(data []byte) error // отправка ответа н сообщение с запросом ответа
}

// Интерфейс модуля сообщений.
// Сервис сообщений должен обеспечивать функции отправки и получения сообщений удобным образом.
type Service interface {
	SendMessage(subj string, data []byte) error                                                                                                          // отправка сообщения обычным способом
	SendStreamingMessage(subj string, data []byte) error                                                                                                 // отправка сообщения с очередью
	SendRequest(subj string, data []byte, timeout time.Duration) (Message, error)                                                                        // отправка сообщения с запросом ответа
	HandleMessages(ctx context.Context, wg *sync.WaitGroup, subj string, handleFunc func(ctx context.Context, wg *sync.WaitGroup, msg Message))          // обработка обычных сообщений
	HandleStreamingMessages(ctx context.Context, wg *sync.WaitGroup, subj string, handleFunc func(ctx context.Context, wg *sync.WaitGroup, msg Message)) // обработка сообщений из очереди
	Sub(subj string) (out <-chan Message, unsub func() error)                                                                                            // подписка на обычные сообщения
	SubStreaming(subj string) (out <-chan Message, unsub func() error)                                                                                   // подписка на сообщения с очередью
	IsConnected() bool                                                                                                                                   // есть соединение
}

// Интерфейс сервиса сохранения и получения больших сообщений.
type BigMsgService interface {
	Put(key []byte, data []byte) (proxydata []byte, err error) // блокирующий метод сохранения данных
	Get(key []byte) ([]byte, error)                            // блокирующий метод получения данных
	Remove(key []byte) error                                   // удаление данных
	DataIsProxy(data []byte) bool                              // проверка, что принятое сообщение является прокси сообщением
	GetKey(proxyData []byte) []byte                            // проверка, что принятое сообщение является прокси сообщением
}
