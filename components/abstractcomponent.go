package components

import "context"

type AbstractComponent struct {
	Cancel context.CancelFunc
}

func (s *AbstractComponent) Stop() {
	if s.Cancel != nil {
		s.Cancel()
	}
}
