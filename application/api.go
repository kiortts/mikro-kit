package application

import (
	"context"
	"sync"
)

type Runnable interface {
	Run(*MainParams) error
}

type MainParams struct {
	Ctx  context.Context
	Wg   *sync.WaitGroup
	Kill func()
}

// общие параметры конфигурации
type CommonParams struct {
	ApplicationId int `arg:"env:APPLICATION_ID"`
	CountryId     int `arg:"env:COUNTRY_ID"`
}
