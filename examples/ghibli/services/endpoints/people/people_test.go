package people_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/kiortts/mikro-kit/examples/ghibli"
	"github.com/kiortts/mikro-kit/examples/ghibli/services/endpoints/people"
	"github.com/kiortts/mikro-kit/examples/ghibli/services/storage"
	"github.com/kiortts/mikro-kit/examples/ghibli/utils"
	"github.com/kiortts/mikro-kit/services/httpserver"
	"github.com/thoas/go-funk"
)

func beforeEach(t *testing.T) (*httpserver.HttpServer, *storage.Storage) {
	log.SetFlags(log.Lshortfile)

	// тестовое хранилище, хэндлер, сервер
	repo := storage.NewMock()
	handler := people.New(repo)
	server := httpserver.New(nil, handler)

	return server, repo
}

func beforeTestGetPerson(t *testing.T) (*httpserver.HttpServer, *ghibli.Person) {

	server, repo := beforeEach(t)

	// получение id любого фильма случайным образом
	var person1 *ghibli.Person
	people, _ := repo.GetPeople()
	for _, person := range people {
		person1 = person
		break
	}

	return server, person1
}

// ответ на корректный запрос
func TestGetPerson1(t *testing.T) {

	server, person1 := beforeTestGetPerson(t)

	req, err := http.NewRequest("GET", fmt.Sprintf("/people/%s", person1.Id), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)
	person2 := new(ghibli.Person)
	err = json.Unmarshal(resp.Body.Bytes(), person2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(person1, person2) {
		t.Errorf("person1 not equal person2")
	}
}

// ответ на запрос с id, который отсутствует в хранилище
func TestGetPerson2(t *testing.T) {

	server, _ := beforeTestGetPerson(t)

	req, err := http.NewRequest("GET", fmt.Sprintf("/people/%s", utils.NewUUID()), nil)
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
func TestGetPerson3(t *testing.T) {

	server, person1 := beforeTestGetPerson(t)

	req, err := http.NewRequest("GET", fmt.Sprintf("/people/%s", person1.Id[1:]), nil)
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
func TestGetPerson4(t *testing.T) {

	server, _ := beforeTestGetPerson(t)
	req, err := http.NewRequest("GET", fmt.Sprintf("/people/%s", ""), nil)
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
func TestGetPerson5(t *testing.T) {

	server, _ := beforeTestGetPerson(t)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/people/%s", ""), nil)
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
func TestGetPerson6(t *testing.T) {

	log.SetFlags(log.Lshortfile)

	// тестовое хранилище, хэндлер, сервер
	repo := storage.NewErr()
	handler := people.New(repo)
	server := httpserver.New(nil, handler)

	req, err := http.NewRequest("GET", fmt.Sprintf("/people/%s", utils.NewUUID()), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("Wrong status code: %d", resp.Code)
	}
}

func beforeTestGetPersons(t *testing.T) (*httpserver.HttpServer, []*ghibli.Person) {

	server, repo := beforeEach(t)

	people1, err := repo.GetPeople()
	if err != nil {
		t.Fatal(err)
	}

	return server, people1
}

// ответ на корректный запрос
func TestGetPersons1(t *testing.T) {

	server, people1 := beforeTestGetPersons(t)

	req, err := http.NewRequest("GET", fmt.Sprintf("/people"), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)

	var films2 []*ghibli.Person
	err = json.Unmarshal(resp.Body.Bytes(), &films2)
	if err != nil {
		t.Fatal(err)
	}

	// писк каждого элемента из people1 в массиве films2
	for _, f1 := range people1 {
		idx := funk.IndexOf(films2, f1)
		if idx < 0 || idx > len(films2)-1 {
			t.Errorf("people1 not equal films2")
		}
	}
}

// ответ на некорректный запрос
func TestGetPersons2(t *testing.T) {

	server, _ := beforeTestGetPersons(t)

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
func TestGetPersons3(t *testing.T) {

	log.SetFlags(log.Lshortfile)

	// тестовое хранилище, хэндлер, сервер
	repo := storage.NewErr()
	handler := people.New(repo)
	server := httpserver.New(nil, handler)

	req, err := http.NewRequest("GET", fmt.Sprintf("/people"), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("Wrong status code: %d", resp.Code)
	}
}
