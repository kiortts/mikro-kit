package componenta

import (
	"context"

	"github.com/kiortts/mikro-kit/components"
)

//
type ComponentA struct {
	components.AbstractComponent
}

// статическая проверка реализаии интерфесов
var _ components.Runnable = (*ComponentA)(nil)

// статические переменные
var cfg *Config

// Constructor
func New(config *Config) *ComponentA {
	checkConfig(config)
	s := &ComponentA{}
	return s
}

// Запуск компонента в работу.
// Реализация интерфейса components.Runnable.
func (s *ComponentA) Run(mainParams *components.MainParams) error {

	// make local context
	wg := mainParams.Wg
	var localCtx context.Context
	localCtx, s.Cancel = context.WithCancel(mainParams.Ctx)

	_ = localCtx
	_ = wg

	return nil
}
