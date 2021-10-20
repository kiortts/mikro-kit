package eva

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kiortts/mikro-kit/application"
	"github.com/kiortts/mikro-kit/components/messages"
)

type Eva struct {
	cancel context.CancelFunc
}

// статическая проверка реализаии интерфесов
var _ application.Runnable = (*Eva)(nil)
var msgServ messages.Service

func New(msgService messages.Service) *Eva {
	msgServ = msgService
	s := &Eva{}
	return s
}

func (s *Eva) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Eva) Run(mainParams *application.MainParams) error {

	log.Println("I'm Eva")

	mainParams.Wg.Add(1)
	go s.mainProcess(mainParams.Ctx, mainParams.Wg)

	return nil
}

func (s *Eva) mainProcess(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()

	<-time.After(time.Millisecond * 500)
	fmt.Println()
	fmt.Println("Step 1 -------------------------------------------------------------")
	fmt.Println()

	<-time.After(time.Second * 4)
	fmt.Println()

	<-time.After(time.Second * 2)
	fmt.Println()
}
