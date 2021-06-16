package locations

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kiortts/mikro-kit/components/httpserver"
	"github.com/kiortts/mikro-kit/examples/ghibli"
	"github.com/kiortts/mikro-kit/examples/ghibli/internal"
	"github.com/kiortts/mikro-kit/examples/ghibli/utils"
)

// Handler
type Handler struct {
	repo ghibli.LocationStorage
}

// проверка реализации типом требуемых интерфейсов
var _ httpserver.Router = (*Handler)(nil)

// New возвращает хэндлер
func New(repo ghibli.LocationStorage) *Handler {
	s := &Handler{
		repo: repo,
	}
	return s
}

// Routes реализация интерфейса HttpRouter, возвращает набор путей и хэндлеров для них
func (s *Handler) Routes() []httpserver.Route {
	routes := []httpserver.Route{

		{
			Name:    "GetLocation",
			Method:  http.MethodGet,
			Pattern: "/locations/{id}",
			Handler: s.location,
		},

		{
			Name:    "GetLocations",
			Method:  http.MethodGet,
			Pattern: "/locations",
			Handler: s.locations,
		},
	}
	return routes
}

// countHandler хэндлер запроса "/locations/{id}"
func (s *Handler) location(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	if !utils.IsValidUUID(id) {
		internal.WriteJSONResponse(w, http.StatusBadRequest, internal.EmptyJSON)
		return
	}

	item, err := s.repo.GetLocation(id)
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

// countHandler хэндлер запроса "/locations"
func (s *Handler) locations(w http.ResponseWriter, r *http.Request) {

	item, err := s.repo.GetLocations()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	internal.WriteItemAsJSON(w, item)
}

// func writeItem(w http.ResponseWriter, item interface{}) {

// 	data, err := json.MarshalIndent(item, "", "   ")
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.Write(data)
// }
