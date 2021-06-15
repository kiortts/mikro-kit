package gorillaserver

import "net/http"

// Router интерфейс для сервисов, с помощью которых http сервер реализует определенный Api
type Router interface {
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
