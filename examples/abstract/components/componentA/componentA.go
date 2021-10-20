package componenta

import (
	"context"

	"github.com/kiortts/mikro-kit/application"
)

//
type ComponentA struct {
	application.AbstractComponent
}

// статическая проверка реализаии интерфесов
var _ application.Runnable = (*ComponentA)(nil)

// статические переменные
var cfg *Config

// Constructor
func New(config *Config) *ComponentA {
	checkConfig(config)
	s := &ComponentA{}
	return s
}

// Запуск компонента в работу.
// Реализация интерфейса application.Runnable.
func (s *ComponentA) Run(mainParams *application.MainParams) error {

	// make local context
	wg := mainParams.Wg
	var localCtx context.Context
	localCtx, s.Cancel = context.WithCancel(mainParams.Ctx)

	_ = localCtx
	_ = wg

	return nil
}
