package api

import (
	"context"
	"net/http"
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

type RawData struct {
	StatusCode int                  `json:"status_code"`
	Message    string               `json:"message"`
	Data       map[int]CategoryData `json:"data"`
	// Data       map[int]map[int]map[string]int `json:"data"`
}

type CategoryData map[int]SubCategoryData
type SubCategoryData map[string]Position
type Position int

type ResultData struct {
	StatusCode int              `json:"status_code"`
	Message    string           `json:"message"`
	Data       map[int]Position `json:"data"`
}

type DataSource interface {
	Get(applicationId, countryId, from, to string) (*RawData, error)
}

type Storage interface {
	Put(key, data []byte) error
	Get(key []byte) ([]byte, bool, error)
}

type DataProcessor interface {
	ProcessRawData(*RawData) (map[string]*ResultData, error)
}

// HttpRouter интерфейс для сервисов, с помощью которых http сервер реализует определенный Api
type HttpRouter interface {
	Routes() []Route // возвращает коллекцию маршрутов
	// RootPath() string // возвращает корневой маршрут
}

// Route определяет один маршрут
type Route struct {
	Name       string           // человекочитаемое название
	Method     string           // тип http метода
	Pattern    string           // путь к эндпоинту
	Handler    http.HandlerFunc // исполняемая функция
	QueryPairs []string
}

// общие параметры конфигурации
type CommonParams struct {
	ApplicationId int `arg:"env:APPLICATION_ID"`
	CountryId     int `arg:"env:COUNTRY_ID"`
}

// ограничитель количества запросов для определенного id (например ip адреса)
type RateLimiter interface {
	Allow(id string) bool
}
