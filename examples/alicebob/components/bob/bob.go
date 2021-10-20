package bob

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/kiortts/mikro-kit/application"
	"github.com/kiortts/mikro-kit/components/messages"
)

type Bob struct {
	cancel context.CancelFunc
}

// статическая проверка реализаии интерфесов
var _ application.Runnable = (*Bob)(nil)
var msgServ messages.Service

func New(msgService messages.Service) *Bob {
	msgServ = msgService
	s := &Bob{}
	return s
}

func (s *Bob) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Bob) Run(mainParams *application.MainParams) error {
	log.Println("I'm Bob")

	mainParams.Wg.Add(1)
	go s.mainProcess(mainParams.Ctx, mainParams.Wg)

	return nil
}

func (s *Bob) mainProcess(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()
	ctxDone := ctx.Done()

	select {
	case <-ctxDone:
		return
	case <-time.After(time.Second * 6):
	}

	log.Println("Bob sub to messages for himself")
	in, _ := msgServ.Sub("Bob1")

	for {
		select {
		case <-ctxDone:
			return
		case msg := <-in:
			log.Printf(`Message for Bob: "%s"`, msg.Data())
		}
	}
}
