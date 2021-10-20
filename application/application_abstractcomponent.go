package application

import (
	"context"
	"sync"
)

type AbstractComponent struct {
	Ctx        context.Context
	Cancel     context.CancelFunc
	Wg         *sync.WaitGroup
	mainParams *MainParams
}

func (s *AbstractComponent) Stop() {
	if s.Cancel != nil {
		s.Cancel()
	}
}

func (s *AbstractComponent) MakeLocalCtxAndWg(mainParams *MainParams) {
	s.mainParams = mainParams
	s.Ctx, s.Cancel = context.WithCancel(mainParams.Ctx)
	s.Wg = new(sync.WaitGroup)
}

func (s *AbstractComponent) WaitAndDo(doAfter ...func()) {
	s.mainParams.Wg.Add(1)
	go func() {
		s.Wg.Wait()
		for _, f := range doAfter {
			f()
		}
		s.mainParams.Wg.Done()
	}()
}
