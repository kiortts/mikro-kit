package natsclient

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/kiortts/mikro-kit/components/messages"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/pkg/errors"
)

// SendMessage публикация сообщения для NATS
func (s *NatsClient) SendMessage(subj string, data []byte) error {

	if s.nsc == nil {
		return fmt.Errorf("Empty nats connection 1")
	}

	// пересылка сообщения через Nats обычным способом
	if bmService == nil || len(data) < messages.BigMessageSize {
		return s.nsc.NatsConn().Publish(subj, data)
	}

	// пересылка через сервис больших сообщений
	log.Printf("Giant message: %d bytes", len(data))
	proxyData, err := bmService.Put([]byte(subj), data)
	if err != nil {
		log.Println(errors.Wrap(err, "bigMsgService Put"))
		return s.nsc.NatsConn().Publish(subj, data)
	}
	return s.nsc.NatsConn().Publish(subj, proxyData)
}

// SendREquest публикация сообщения для NATS и NATS Streaming
func (s *NatsClient) SendStreamingMessage(subj string, data []byte) error {

	if s.nsc == nil {
		return fmt.Errorf("Empty nats connection 2")
	}

	// пересылка сообщения через Nats
	if bmService == nil || len(data) < messages.BigMessageSize {
		return s.nsc.Publish(subj, data)
	}

	// пересылка через сервис больших сообщений
	log.Printf("Giant streaming message: %d bytes", len(data))
	proxyData, err := bmService.Put([]byte(subj), data)
	if err != nil {
		log.Println(errors.Wrap(err, "bigMsgService Put err"))
		return s.nsc.Publish(subj, data)
	}
	return s.nsc.Publish(subj, proxyData)
}

// SendRequest публикация сообщения с запросом ответа
func (s *NatsClient) SendRequest(subj string, data []byte, timeout time.Duration) (messages.Message, error) {

	if s.nsc == nil {
		return nil, fmt.Errorf("Empty nats connection 3")
	}

	// отправка запроса через Nats обычным способом
	if bmService == nil || len(data) < messages.BigMessageSize {
		return s.requestMsg(subj, data, timeout)
	}

	// пересылка через сервис больших сообщений
	log.Printf("Giant request message: %d bytes", len(data))
	proxyData, err := bmService.Put([]byte(subj), data)
	if err != nil {
		log.Println(errors.Wrap(err, "bigMsgService Put err"))
		// msg1, err := s.nsc.NatsConn().Request(subj, data, timeout)
		// msg2 := &NatsMsg{msg1}
		return s.requestMsg(subj, data, timeout)
	}
	return s.requestMsg(subj, proxyData, timeout)
}

func (s *NatsClient) requestMsg(subj string, data []byte, timeout time.Duration) (messages.Message, error) {

	msg1, err := s.nsc.NatsConn().Request(subj, data, timeout)
	if err != nil {
		return nil, errors.Wrap(err, "NatsConn Request err")
	}

	msg2 := &NatsMsg{msg1}
	err = bigMessageCheck(msg2, subj) // FIXME: ключ должен быть из сообщения
	if err != nil {
		return nil, errors.Wrap(err, "bigMessageCheck err")
	}

	go func() {
		<-time.After(timeout + time.Second*10)
		bmService.Remove([]byte(subj))
	}()

	return msg2, nil
}

func (s *NatsClient) HandleMessages(ctx context.Context, wg *sync.WaitGroup, subj string, handleMessageFunc func(ctx context.Context, wg *sync.WaitGroup, msg messages.Message)) {

	killed := ctx.Done()
	wg.Add(1)
	ctx, _ = context.WithCancel(ctx)

	go func() {
		defer wg.Done()

		if !s.waitingForConnectionWithNatsService(ctx) {
			return
		}

		// подписка на NATS
		_, ch, unsub := subToNats(ctx, s.nsc.NatsConn(), subj) // подписка на NATS
		defer unsub()
		var msg *nats.Msg

		// рабочий цикл хэндлера

		for {
			select {
			// прекращение работы хэндлера
			case <-killed: // контекст отменен
				return

			// новое сообщение
			case msg = <-ch:
				msg2 := &NatsMsg{msg}            // обертка исходного сообщения
				bigMessageCheck(msg2, subj)      // проверка на большое сообщение (с получением данных из сервиса больших сообщений)
				handleMessageFunc(ctx, wg, msg2) // передач сообщения хэндлеру
			}
		}
	}()
}

func (s *NatsClient) HandleStreamingMessages(ctx context.Context, wg *sync.WaitGroup, subj string, handleMessageFunc func(ctx context.Context, wg *sync.WaitGroup, msg messages.Message)) {

	// проверка пути подписки на wildcard, для NATS Stream не поддерживается
	if strings.Contains(subj, "*") {
		log.Fatal("NATS Streaming server does not support wildcard for channels, that is, one cannot subscribe on foo.*, or foo.>, etc...")
	}

	killed := ctx.Done()
	wg.Add(1)
	ctx, _ = context.WithCancel(ctx)

	go func() {
		defer wg.Done()

		if !s.waitingForConnectionWithNatsService(ctx) {
			return
		}

		// подписка на NATS Streaming
		// _, ch, unsub := SubToNatsStreaming(ctx, conn, params.NatsPath)
		natsSub := asyncSubToNatsStreaming(ctx, s.nsc, subj)
		defer natsSub.Unsub()
		var msg *stan.Msg

		// рабочий цикл хэндлера
	Loop1:
		for {
			select {
			// прекращение работы хэндлера
			case <-killed: // контекст отменен
				break Loop1
			// новое сообщение
			case msg = <-natsSub.MsgChan:
				msg2 := &StanMsg{msg}            // обертка исходного сообщения
				bigMessageCheck(msg2, subj)      // проверка на большое сообщение (с получением данных из сервиса больших сообщений)
				handleMessageFunc(ctx, wg, msg2) // исполняется переданная функция обработки сообщения
			}
		}
	}()
}

// Блокирующий метод ожидания установки соединения с сервером NATS.
// Возвращает true если соединение было установлено и false, если контекст был завершен до установки соединения.
func (s *NatsClient) waitingForConnectionWithNatsService(ctx context.Context) bool {

	killed := ctx.Done()

	// циклическая проверка s.nsc на nil
	for {
		if s.nsc != nil {
			return true
		}
		select {
		case <-time.After(time.Millisecond * 100):

		case <-killed:
			return false
		}
	}
}

func subToNats(ctx context.Context, nc *nats.Conn, natsPath string) (*nats.Subscription, chan *nats.Msg, func() error) {

	killed := ctx.Done()
	ch := make(chan *nats.Msg, 1024) // TODO: проверить какая емкость канала влияет
	var sub *nats.Subscription
	var err error

	// пустая функция имитирующая отмену отмену подписки в случае если сама подписка получена не была
	unsubMock := func() error { return nil }

Loop1:
	for {

		select {
		case <-killed:
			return nil, ch, unsubMock // возврат заглушки
		default:
		}

		if nc == nil {
			time.Sleep(time.Second)
			continue
		}

		sub, err = nc.ChanSubscribe(natsPath, ch)
		if err == nil {

			log.Printf("Get NATS subscription: %s", sub.Subject)
			break Loop1
		}
		log.Fatalf("Nats Subscribe Err: %s path:%s", err.Error(), natsPath)
		time.Sleep(time.Second)
	}

	unsub := sub.Unsubscribe

	return sub, ch, unsub
}

// asyncSubToNatsStreaming подписка на NATS Streaming.
// При подписке из NATS Streaming будет получено последнее по очереди сообщение.
// В канале всегда может находится только одно последнее сообщение, перед отправкой нового канал очищается от предыдущего сообщения.
func asyncSubToNatsStreaming(ctx context.Context, nsc stan.Conn, natsPath string) (natsSub *NatsSubscription) {

	killed := ctx.Done()
	var sub stan.Subscription
	var err error

	natsSub = &NatsSubscription{
		MsgChan:   make(chan *stan.Msg, 1),
		natsUnsub: func() error { return nil }, // пустая функция имитирующая отмену отмену подписки в случае если сама подписка получена не была
	}

	go func() {

	Loop1:
		for {

			select {
			case <-killed:
				break Loop1
			default:
			}

			if nsc == nil {
				time.Sleep(time.Second)
				continue
			}

			// Subscribe starting with most recently published value
			// новое сообщение будет отправлено в канал natsSub.msgChan
			sub, err = nsc.Subscribe(natsPath, func(m *stan.Msg) {

				// очистка канала от последнего сообщения
				select {
				case <-natsSub.MsgChan:
				default:
				}

				// отправка нового сообщения
				select {
				case natsSub.MsgChan <- m:
				default:

				}

			}, stan.StartWithLastReceived())

			if err == nil {
				break
			}
			time.Sleep(time.Second)
		}

		natsSub.mx.Lock()
		if sub != nil {
			natsSub.natsUnsub = sub.Unsubscribe
		}
		natsSub.mx.Unlock()
	}()

	return natsSub
}

type NatsSubscription struct {
	MsgChan   chan *stan.Msg
	mx        sync.RWMutex
	natsUnsub func() error
}

func (n *NatsSubscription) Unsub() error {
	n.mx.Lock()
	err := n.natsUnsub()
	n.mx.Unlock()
	return err
}

// обычная подпика
func (s *NatsClient) Sub(subj string) (<-chan messages.Message, func() error) {
	ctx := context.Background()
	wg := new(sync.WaitGroup)
	return s.SubWithContext(ctx, wg, subj)
}

// обычная подпика с контекстом
func (s *NatsClient) SubWithContext(ctx context.Context, wg *sync.WaitGroup, subj string) (<-chan messages.Message, func() error) {

	mx := new(sync.RWMutex)
	localCtx, cancel := context.WithCancel(ctx)

	ctxDone := localCtx.Done()
	natsMsgChan := make(chan *nats.Msg, 1024)           // TODO: проверить какая емкость канала влияет
	wrappedMsgChan := make(chan messages.Message, 1024) // TODO: проверить какая емкость канала влияет
	var sub *nats.Subscription
	var err error

	// пустая функция имитирующая отмену отмену подписки в случае если сама подписка получена не была
	unsubMock := func() error { return nil }
	var unsub func() error
	unsub = unsubMock

	callUnsubFunc := func() error {
		mx.RLock()
		defer mx.RUnlock()

		cancel()
		return unsub()
	}

	wg.Add(1)
	go func() {

		defer wg.Done()

		// подписка на subj

		for {

			select {
			case <-ctxDone:
				return
			default:
			}

			if s.nsc == nil {
				<-time.After(time.Second)
				continue
			}

			sub, err = s.nsc.NatsConn().ChanSubscribe(subj, natsMsgChan)
			if err == nil {
				log.Printf("Get NATS subscription: %s", sub.Subject)
				mx.Lock()
				unsub = sub.Unsubscribe
				mx.Unlock()
				break
			}
			log.Printf("Nats Subscribe Err: %s subject:%s", err.Error(), subj)
			<-time.After(time.Second)
		}

		// обертывание приходящих от сервера NATS сообщений и отправка в канал потребителям

		for {
			select {
			case <-ctxDone:
				return
			case natsMsg := <-natsMsgChan:
				// // очистка канала
				// select {
				// case <-wrappedMsgChan:
				// default:
				// }
				// обертывание и отправка сообщения
				select {
				case wrappedMsgChan <- &NatsMsg{natsMsg}:
				default:
				}
			}
		}

	}()

	return wrappedMsgChan, callUnsubFunc
}

// подпика на Streaming
func (s *NatsClient) SubStreaming(subj string) (<-chan messages.Message, func() error) {
	ctx := context.Background()
	wg := new(sync.WaitGroup)
	return s.SubStreamingWithContext(ctx, wg, subj)
}

// подпика на Streaming с контекстом
func (s *NatsClient) SubStreamingWithContext(ctx context.Context, wg *sync.WaitGroup, subj string) (<-chan messages.Message, func() error) {

	mx := new(sync.RWMutex)
	localCtx, cancel := context.WithCancel(ctx)
	_ = localCtx

	// ctxDone := localCtx.Done()
	// natsMsgChan := make(chan *nats.Msg, 1024)           // TODO: проверить какая емкость канала влияет
	wrappedMsgChan := make(chan messages.Message, 1024) // TODO: проверить какая емкость канала влияет
	// var sub stan.Subscription
	// var err error

	// // пустая функция имитирующая отмену отмену подписки в случае если сама подписка получена не была
	unsubMock := func() error { return nil }
	unsub := unsubMock

	callUnsubFunc := func() error {
		mx.RLock()
		defer mx.RUnlock()
		cancel()
		return unsub()
	}

	// wg.Add(1)
	// go func() {

	// 	defer wg.Done()

	// 	// подписка на subj

	// 	for {

	// 		select {
	// 		case <-ctxDone:
	// 			return
	// 		default:
	// 		}

	// 		if s.nsc == nil {
	// 			<-time.After(time.Second)
	// 			continue
	// 		}

	// 				// Subscribe starting with most recently published value
	// 		// новое сообщение будет отправлено в канал natsSub.msgChan
	// 		sub, err = s.nsc.Subscribe(subj, func(m *stan.Msg) {

	// 			// очистка канала от последнего сообщения
	// 			select {
	// 			case <-natsSub.MsgChan:
	// 			default:
	// 			}

	// 			// отправка нового сообщения
	// 			select {
	// 			case natsSub.MsgChan <- m:
	// 			default:

	// 			}

	// 		}, stan.StartWithLastReceived())
	// 		if err == nil {
	// 			log.Printf("Get NATS subscription: %s", sub.Subject)
	// 			mx.Lock()
	// 			unsub = sub.Unsubscribe
	// 			mx.Unlock()
	// 			break
	// 		}
	// 		log.Printf("Nats Subscribe Err: %s subject:%s", err.Error(), subj)
	// 		<-time.After(time.Second)
	// 	}

	// 	// обертывание приходящих от сервера NATS сообщений и отправка в канал потребителям

	// 	for {
	// 		select {
	// 		case <-ctxDone:
	// 			return
	// 		case natsMsg := <-natsMsgChan:
	// 			select {
	// 			case wrappedMsgChan <- &NatsMsg{natsMsg}:
	// 			default:
	// 			}
	// 		}
	// 	}

	// }()

	return wrappedMsgChan, callUnsubFunc
}

func bigMessageCheck(msg2 messages.Message, subj string) error {

	// проверка на большое сообщение
	if bmService == nil {
		return nil
	}

	if !bmService.DataIsProxy(msg2.Data()) {
		return nil
	}

	key := bmService.GetKey(msg2.Data())

	bigData, err := bmService.Get(key)
	if err != nil {
		return errors.Wrap(err, "Big message Get err")
	}
	msg2.SetData(bigData)

	return nil
}
