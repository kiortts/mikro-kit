package films

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/kiortts/mikro-kit/examples/ghibli"
	"github.com/kiortts/mikro-kit/examples/ghibli/services/storage"
	"github.com/kiortts/mikro-kit/examples/ghibli/utils"
	"github.com/kiortts/mikro-kit/services/httpserver/gorillaserver"
	"github.com/thoas/go-funk"
)

func beforeEach(t *testing.T) (*gorillaserver.GorillaServer, *storage.Storage) {
	log.SetFlags(log.Lshortfile)

	// тестовое хранилище, хэндлер, сервер
	repo := storage.NewMock()
	filmHandler := New(repo)
	server := gorillaserver.New(nil, filmHandler)

	return server, repo
}

func beforeTestGetFilm(t *testing.T) (*gorillaserver.GorillaServer, *ghibli.Film) {

	server, repo := beforeEach(t)

	// получение id любого фильма случайным образом
	var film1 *ghibli.Film
	films, _ := repo.GetFilms()
	for _, film := range films {
		film1 = film
		break
	}

	return server, film1
}

// ответ на корректный запрос
func TestGetFilm1(t *testing.T) {

	server, film1 := beforeTestGetFilm(t)

	req, err := http.NewRequest("GET", fmt.Sprintf("/films/%s", film1.Id), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)
	film2 := new(ghibli.Film)
	err = json.Unmarshal(resp.Body.Bytes(), film2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(film1, film2) {
		t.Errorf("film1 not equal film2")
	}
}

// ответ на запрос с id, который отсутствует в хранилище
func TestGetFilm2(t *testing.T) {

	server, _ := beforeTestGetFilm(t)

	req, err := http.NewRequest("GET", fmt.Sprintf("/films/%s", utils.NewUUID()), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)
	if resp.Code != http.StatusNotFound {
		t.Errorf("Wrong status code: %d", resp.Code)
	}

}

// ответ на запрос с id, не соответствующим uuid
func TestGetFilm3(t *testing.T) {

	server, film1 := beforeTestGetFilm(t)

	req, err := http.NewRequest("GET", fmt.Sprintf("/films/%s", film1.Id[1:]), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Errorf("Wrong status code: %d", resp.Code)
	}
}

// ответ на запрос с пустым id
func TestGetFilm4(t *testing.T) {

	server, _ := beforeTestGetFilm(t)
	req, err := http.NewRequest("GET", fmt.Sprintf("/films/%s", ""), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)
	if resp.Code != http.StatusNotFound {
		t.Errorf("Wrong status code: %d", resp.Code)
	}
}

// ответ на запрос с неправильным http методом
func TestGetFilm5(t *testing.T) {

	server, _ := beforeTestGetFilm(t)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/films/%s", ""), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)
	if resp.Code != http.StatusNotFound {
		t.Errorf("Wrong status code: %d", resp.Code)
	}

}

// ошибка хранилища
func TestGetFilm6(t *testing.T) {

	log.SetFlags(log.Lshortfile)

	// тестовое хранилище, хэндлер, сервер
	repo := storage.NewErr()
	filmHandler := New(repo)
	server := gorillaserver.New(nil, filmHandler)

	req, err := http.NewRequest("GET", fmt.Sprintf("/films/%s", utils.NewUUID()), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("Wrong status code: %d", resp.Code)
	}
}

func beforeTestGetFilms(t *testing.T) (*gorillaserver.GorillaServer, []*ghibli.Film) {

	server, repo := beforeEach(t)

	films1, err := repo.GetFilms()
	if err != nil {
		t.Fatal(err)
	}

	return server, films1
}

// ответ на корректный запрос
func TestGetFilms1(t *testing.T) {

	server, films1 := beforeTestGetFilms(t)

	req, err := http.NewRequest("GET", fmt.Sprintf("/films"), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)

	var films2 []*ghibli.Film
	err = json.Unmarshal(resp.Body.Bytes(), &films2)
	if err != nil {
		t.Fatal(err)
	}

	// писк каждого элемента из films1 в массиве films2
	for _, f1 := range films1 {
		idx := funk.IndexOf(films2, f1)
		if idx < 0 || idx > len(films2)-1 {
			t.Errorf("films1 not equal films2")
		}
	}
}

// ответ на некорректный запрос
func TestGetFilms2(t *testing.T) {

	server, _ := beforeTestGetFilms(t)

	req, err := http.NewRequest("GET", fmt.Sprintf("/fools"), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Wrong status code: %d", resp.Code)
	}
}

// ошибка хранилища
func TestGetFilms3(t *testing.T) {

	log.SetFlags(log.Lshortfile)

	// тестовое хранилище, хэндлер, сервер
	repo := storage.NewErr()
	filmHandler := New(repo)
	server := gorillaserver.New(nil, filmHandler)

	req, err := http.NewRequest("GET", fmt.Sprintf("/films"), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("Wrong status code: %d", resp.Code)
	}
}
