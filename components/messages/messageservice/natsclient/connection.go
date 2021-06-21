package natsclient

import (
	"context"
	"sync"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	uuid "github.com/satori/go.uuid"
)

var bigInt int = 1<<63 - 1 // очень большое число

func (s *NatsClient) getConnection(ctx context.Context, wg *sync.WaitGroup, cfg *Config) {

	killed := ctx.Done()
	wg.Add(1)

	go func() {
		defer wg.Done()

		var nc *nats.Conn
		var err error

	Loop1:
		for {

			// соединение с таймаутом, чтобы не висеть в нем вечно, а отписываться хэлзчекеру
			nc, err = nats.Connect(cfg.Url, nats.Timeout(time.Second), nats.MaxReconnects(bigInt))

			if err != nil {
				// log.Printf(dict.LOG_CONNECTED_TO_NATS_ERR, err)
				select {

				// отмена контекста
				case <-killed:
					return
				// еще попытка
				case <-time.After(time.Second):
					// imHealthyNow()
					continue Loop1
				}

			}

			// log.Printf(dict.LOG_CONNECTED_TO_NATS, nc.ConnectedUrl())

			closedHandler := func(nc *nats.Conn) {
				// log.Println(di/ct.LOG_NATS_CONN_CLOSED)
			}
			disconnectHandlae := func(nc *nats.Conn, err error) {
				if err != nil {
					// log.Printf(dict.LOG_NATS_DESCONNECTED+" err: %s", err.Error())
					return
				}
				// log.Printf(dict.LOG_NATS_DESCONNECTED)
			}
			reconnectHandler := func(nc *nats.Conn) {
				// log.Println(dict.LOG_NATS_RECONNECTED)
			}

			nc.SetDisconnectErrHandler(disconnectHandlae)
			nc.SetReconnectHandler(reconnectHandler)
			nc.SetClosedHandler(closedHandler)

			break Loop1

		}

		// если имя клиента не задано

		clientID := uuid.Must(uuid.NewV4(), nil).String()

	Loop2:
		for {

			// соединение с таймаутом, чтобы не висеть в нем вечно, а отписываться хэлзчекеру
			nsc, err := stan.Connect(cfg.StreamingCluster, clientID, stan.NatsConn(nc), stan.ConnectWait(time.Second))

			if err != nil {
				select {

				// отмена контекста
				case <-killed: // отмена контекста
					return

				// еще попытка
				case <-time.After(time.Second):
					// imHealthyNow()
					continue Loop2
				}

			}

			// log.Printf(dict.LOG_CONNECTED_TO_NATS_ST, nc.ConnectedUrl())
			s.nsc = nsc
			break Loop2
		}
	}()

}

func (s *NatsClient) runCloser(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		s.nsc.Close()
		s.nsc.NatsConn().Close()
	}()
}
