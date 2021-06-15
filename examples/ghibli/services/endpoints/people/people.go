package people

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kiortts/mikro-kit/examples/ghibli"
	"github.com/kiortts/mikro-kit/examples/ghibli/internal"
	"github.com/kiortts/mikro-kit/examples/ghibli/utils"
	"github.com/kiortts/mikro-kit/services/http/gorillaserver"
)

// Handler
type Handler struct {
	repo ghibli.PersonStorage
}

// проверка реализации типом требуемых интерфейсов
var _ gorillaserver.Router = (*Handler)(nil)

// New возвращает хэндлер
func New(repo ghibli.PersonStorage) *Handler {
	s := &Handler{
		repo: repo,
	}
	return s
}

// Routes реализация интерфейса HttpRouter, возвращает набор путей и хэндлеров для них
func (s *Handler) Routes() []gorillaserver.Route {
	routes := []gorillaserver.Route{

		{
			Name:    "GetPerson",
			Method:  http.MethodGet,
			Pattern: "/people/{id}",
			Handler: s.person,
		},

		{
			Name:    "GetPeople",
			Method:  http.MethodGet,
			Pattern: "/people",
			Handler: s.people,
		},
	}
	return routes
}

// countHandler хэндлер запроса "/person/{id}"
func (s *Handler) person(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	if !utils.IsValidUUID(id) {
		internal.WriteJSONResponse(w, http.StatusBadRequest, internal.EmptyJSON)
		return
	}

	item, err := s.repo.GetPerson(id)
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

// countHandler хэндлер запроса "/people"
func (s *Handler) people(w http.ResponseWriter, r *http.Request) {

	item, err := s.repo.GetPeople()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	internal.WriteItemAsJSON(w, item)
}
