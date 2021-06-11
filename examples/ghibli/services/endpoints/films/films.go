package films

import (
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/kiortts/mikro-kit/examples/ghibli"
	"github.com/kiortts/mikro-kit/examples/ghibli/internal"
	"github.com/kiortts/mikro-kit/services/httpserver"
)

// Handler
type Handler struct {
	repo ghibli.FilmRepo
}

// проверка реализации типом требуемых интерфейсов
var _ httpserver.Router = (*Handler)(nil)

// New возвращает хэндлер
func New(repo ghibli.FilmRepo) *Handler {
	s := &Handler{
		repo: repo,
	}
	return s
}

// Routes реализация интерфейса HttpRouter, возвращает набор путей и хэндлеров для них
func (s *Handler) Routes() []httpserver.Route {
	routes := []httpserver.Route{

		{
			Name:    "GetFilm",
			Method:  http.MethodGet,
			Pattern: "/films/{id}",
			Handler: s.film,
		},

		{
			Name:    "GetFilms",
			Method:  http.MethodGet,
			Pattern: "/films",
			Handler: s.films,
		},
	}
	return routes
}

// countHandler хэндлер запроса "/film/{id}"
func (s *Handler) film(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	if !IsValidUUID(id) {
		internal.WriteJSONResponse(w, http.StatusBadRequest, internal.EmptyJSON)
		return
	}

	item, err := s.repo.GetFilm(id)
	if err != nil {
		switch err {
		case internal.NotFoundError:
			internal.WriteJSONResponse(w, http.StatusNotFound, internal.EmptyJSON)
		default:
			internal.WriteJSONResponse(w, http.StatusInternalServerError, internal.EmptyJSON)
		}
		return
	}

	internal.WriteItemAsJSON(w, item)
}

// countHandler хэндлер запроса "/films"
func (s *Handler) films(w http.ResponseWriter, r *http.Request) {

	item, err := s.repo.GetFilms()
	if err != nil {
		internal.WriteJSONResponse(w, http.StatusInternalServerError, internal.EmptyJSON)
		return
	}

	internal.WriteItemAsJSON(w, item)
}

// сравнение строки с шаблоном uuid
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
