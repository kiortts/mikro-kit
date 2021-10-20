package alice

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/kiortts/mikro-kit/application"
	"github.com/kiortts/mikro-kit/components/messages"
)

type Alice struct {
	cancel context.CancelFunc
}

// статическая проверка реализаии интерфесов
var _ application.Runnable = (*Alice)(nil)
var msgServ messages.Service

func New(msgService messages.Service) *Alice {
	msgServ = msgService
	s := &Alice{}
	return s
}

func (s *Alice) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Alice) Run(mainParams *application.MainParams) error {

	log.Println("I'm Alice")

	mainParams.Wg.Add(1)
	go s.mainProcess(mainParams.Ctx, mainParams.Wg)

	return nil
}

func (s *Alice) mainProcess(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()

	<-time.After(time.Second * 2)
	s.sendMsgToBob("Bob1", "Hello from Alice!")

	<-time.After(time.Second * 2)
	s.sendMsgToBob("Bob1", "Hello from Alice!!")

	<-time.After(time.Second * 4)
	s.sendMsgToBob("Bob1", "Hello from Alice!!!")
}

func (s *Alice) sendMsgToBob(subj string, msg string) {
	err := msgServ.SendMessage(subj, []byte(msg))
	if err != nil {
		log.Printf("Alice SendMessage err: %s", err)
	} else {
		log.Printf(`Alice say to Bob: "%s"`, msg)
	}

}
