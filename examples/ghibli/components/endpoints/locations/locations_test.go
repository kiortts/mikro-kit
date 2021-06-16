package locations_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/kiortts/mikro-kit/components/httpserver/gorillaserver"
	"github.com/kiortts/mikro-kit/examples/ghibli"
	"github.com/kiortts/mikro-kit/examples/ghibli/components/endpoints/locations"
	"github.com/kiortts/mikro-kit/examples/ghibli/components/storage"
	"github.com/kiortts/mikro-kit/examples/ghibli/utils"
	"github.com/thoas/go-funk"
)

func beforeEach(t *testing.T) (*gorillaserver.GorillaServer, *storage.Storage) {
	log.SetFlags(log.Lshortfile)

	// тестовое хранилище, хэндлер, сервер
	repo := storage.NewMock()
	handler := locations.New(repo)
	server := gorillaserver.New(nil, handler)

	return server, repo
}

func beforeTestGetLocation(t *testing.T) (*gorillaserver.GorillaServer, *ghibli.Location) {

	server, stor := beforeEach(t)

	// получение id любого фильма случайным образом
	var location1 *ghibli.Location
	locations, _ := stor.GetLocations()
	for _, location := range locations {
		location1 = location
		break
	}

	return server, location1
}

// ответ на корректный запрос
func TestGetLocation1(t *testing.T) {

	server, location1 := beforeTestGetLocation(t)

	req, err := http.NewRequest("GET", fmt.Sprintf("/locations/%s", location1.Id), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)
	location2 := new(ghibli.Location)
	err = json.Unmarshal(resp.Body.Bytes(), location2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(location1, location2) {
		t.Errorf("location1 not equal location2")
	}
}

// ответ на запрос с id, который отсутствует в хранилище
func TestGetLocation2(t *testing.T) {

	server, _ := beforeTestGetLocation(t)

	req, err := http.NewRequest("GET", fmt.Sprintf("/locations/%s", utils.NewUUID()), nil)
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
func TestGetLocation3(t *testing.T) {

	server, location1 := beforeTestGetLocation(t)

	req, err := http.NewRequest("GET", fmt.Sprintf("/locations/%s", location1.Id[1:]), nil)
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
func TestGetLocation4(t *testing.T) {

	server, _ := beforeTestGetLocation(t)
	req, err := http.NewRequest("GET", fmt.Sprintf("/locations/%s", ""), nil)
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
func TestGetLocation5(t *testing.T) {

	server, _ := beforeTestGetLocation(t)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/locations/%s", ""), nil)
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
func TestGetLocation6(t *testing.T) {

	log.SetFlags(log.Lshortfile)

	// тестовое хранилище, хэндлер, сервер
	repo := storage.NewErr()
	handler := locations.New(repo)
	server := gorillaserver.New(nil, handler)

	req, err := http.NewRequest("GET", fmt.Sprintf("/locations/%s", utils.NewUUID()), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("Wrong status code: %d", resp.Code)
	}
}

func beforeTestGetLocations(t *testing.T) (*gorillaserver.GorillaServer, []*ghibli.Location) {

	server, stor := beforeEach(t)

	location1, err := stor.GetLocations()
	if err != nil {
		t.Fatal(err)
	}

	return server, location1
}

// ответ на корректный запрос
func TestGetLocations1(t *testing.T) {

	server, location1 := beforeTestGetLocations(t)

	req, err := http.NewRequest("GET", fmt.Sprintf("/locations"), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)

	var films2 []*ghibli.Location
	err = json.Unmarshal(resp.Body.Bytes(), &films2)
	if err != nil {
		t.Fatal(err)
	}

	// писк каждого элемента из location1 в массиве films2
	for _, f1 := range location1 {
		idx := funk.IndexOf(films2, f1)
		if idx < 0 || idx > len(films2)-1 {
			t.Errorf("location1 not equal films2")
		}
	}
}

// ответ на некорректный запрос
func TestGetLocations2(t *testing.T) {

	server, _ := beforeTestGetLocations(t)

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
func TestGetLocations3(t *testing.T) {

	log.SetFlags(log.Lshortfile)

	// тестовое хранилище, хэндлер, сервер
	repo := storage.NewErr()
	handler := locations.New(repo)
	server := gorillaserver.New(nil, handler)

	req, err := http.NewRequest("GET", fmt.Sprintf("/locations"), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	server.Router().ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("Wrong status code: %d", resp.Code)
	}
}
