package storage_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/kiortts/mikro-kit/application"
	"github.com/kiortts/mikro-kit/examples/ghibli"
	"github.com/kiortts/mikro-kit/examples/ghibli/internal"
	"github.com/kiortts/mikro-kit/examples/ghibli/services/storage"
	"github.com/kiortts/mikro-kit/examples/ghibli/utils"
)

func getAnyFilm(stor *storage.Storage) (*ghibli.Film, error) {
	films, _ := stor.GetFilms()
	for _, film := range films {
		return film, nil

	}
	return nil, fmt.Errorf("Empty stor")
}

func getAnyPerson(stor *storage.Storage) (*ghibli.Person, error) {
	items, _ := stor.GetPeople()
	for _, item := range items {
		return item, nil

	}
	return nil, fmt.Errorf("Empty stor")
}

func getAnyLocation(stor *storage.Storage) (*ghibli.Location, error) {
	items, _ := stor.GetLocations()
	for _, item := range items {
		return item, nil

	}
	return nil, fmt.Errorf("Empty stor")
}

func getAnySpecies(stor *storage.Storage) (*ghibli.Species, error) {
	items, _ := stor.GetAllSpecies()
	for _, item := range items {
		return item, nil

	}
	return nil, fmt.Errorf("Empty stor")
}

func getAnyVehicle(stor *storage.Storage) (*ghibli.Vehicle, error) {
	items, _ := stor.GetVehicles()
	for _, item := range items {
		return item, nil

	}
	return nil, fmt.Errorf("Empty stor")
}

// корректный id
func GetFilmValidId(t *testing.T, stor *storage.Storage) {
	film1, err := getAnyFilm(stor)
	if err != nil {
		t.Errorf("getAnyFilmId err: %v", err)
	}
	film2, err := stor.GetFilm(film1.Id)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	if !reflect.DeepEqual(film1, film2) {
		t.Errorf("films not equal")
	}
}

// несуществующий id
func GetFilmNonexistId(t *testing.T, stor *storage.Storage) {
	_, err := stor.GetFilm(utils.NewUUID())
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// некорректный id
func GetFilmInvalidId(t *testing.T, stor *storage.Storage) {
	_, err := stor.GetFilm("")
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// корректный id
func GetPersonValidId(t *testing.T, stor *storage.Storage) {
	item1, err := getAnyPerson(stor)
	if err != nil {
		t.Errorf("getAnyPerson err: %v", err)
	}
	item2, err := stor.GetPerson(item1.Id)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	if !reflect.DeepEqual(item1, item2) {
		t.Errorf("items not equal")
	}
}

// несуществующий id
func GetPersonNonexistId(t *testing.T, stor *storage.Storage) {
	_, err := stor.GetPerson(utils.NewUUID())
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// некорректный id
func GetPersonInvalidId(t *testing.T, stor *storage.Storage) {
	_, err := stor.GetPerson("")
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// корректный id
func GetLocationValidId(t *testing.T, stor *storage.Storage) {
	item1, err := getAnyLocation(stor)
	if err != nil {
		t.Errorf("getAnyPerson err: %v", err)
	}
	item2, err := stor.GetLocation(item1.Id)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	if !reflect.DeepEqual(item1, item2) {
		t.Errorf("items not equal")
	}
}

// несуществующий id
func GetLocationNonexistId(t *testing.T, stor *storage.Storage) {
	_, err := stor.GetLocation(utils.NewUUID())
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// некорректный id
func GetLocationInvalidId(t *testing.T, stor *storage.Storage) {
	_, err := stor.GetLocation("")
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// корректный id
func GetSpeciesValidId(t *testing.T, stor *storage.Storage) {
	item1, err := getAnySpecies(stor)
	if err != nil {
		t.Errorf("getAnyPerson err: %v", err)
	}
	item2, err := stor.GetSpecies(item1.Id)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	if !reflect.DeepEqual(item1, item2) {
		t.Errorf("items not equal")
	}
}

// несуществующий id
func GetSpeciesNonexistId(t *testing.T, stor *storage.Storage) {
	_, err := stor.GetSpecies(utils.NewUUID())
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// некорректный id
func GetSpeciesInvalidId(t *testing.T, stor *storage.Storage) {
	_, err := stor.GetSpecies("")
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// корректный id
func GetVehicleValidId(t *testing.T, stor *storage.Storage) {
	item1, err := getAnyVehicle(stor)
	if err != nil {
		t.Errorf("getAnyPerson err: %v", err)
	}
	item2, err := stor.GetVehicle(item1.Id)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	if !reflect.DeepEqual(item1, item2) {
		t.Errorf("items not equal")
	}
}

// несуществующий id
func GetVehicleNonexistId(t *testing.T, stor *storage.Storage) {
	_, err := stor.GetVehicle(utils.NewUUID())
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// некорректный id
func GetVehicleInvalidId(t *testing.T, stor *storage.Storage) {
	_, err := stor.GetVehicle("")
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

func TestGetFilmMock1(t *testing.T) {
	stor := storage.NewMock()
	GetFilmValidId(t, stor)
}

func TestGetFilmMock2(t *testing.T) {
	stor := storage.NewMock()
	GetFilmNonexistId(t, stor)
}

func TestGetFilmMock3(t *testing.T) {
	stor := storage.NewMock()
	GetFilmInvalidId(t, stor)
}

func setEnv() {

	log.SetFlags(log.Lshortfile)

	os.Setenv("FILMS_FILE", "./src/films.js")
	os.Setenv("PEOPLE_FILE", "./src/people.js")
	os.Setenv("LOCATIONS_FILE", "./src/locations.js")
	os.Setenv("SPECIES_FILE", "./src/species.js")
	os.Setenv("VEHICLES_FILE", "./src/vehicles.js")
}

func getStorage(t *testing.T) (*storage.Storage, error) {

	stor := storage.NewLocal()
	main := &application.MainParams{Ctx: context.Background(), Wg: new(sync.WaitGroup), Kill: func() {}}

	err := stor.Run(main)
	if err != nil {
		return nil, fmt.Errorf("stor.Run err: %v", err)
	}

	return stor, nil
}

// корректный id
func TestGetFilm1(t *testing.T) {
	setEnv()
	stor, err := getStorage(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetFilmValidId(t, stor)
}

// несуществующий id
func TestGetFilm2(t *testing.T) {
	setEnv()
	stor, err := getStorage(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetFilmNonexistId(t, stor)
}

// некорректный id
func TestGetFilm3(t *testing.T) {
	setEnv()
	stor, err := getStorage(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetFilmInvalidId(t, stor)
}

// некорректно указан файл фильмов

func TestGetFilm4(t *testing.T) {
	setEnv()
	os.Setenv("FILMS_FILE", "")
	_, err := getStorage(t)
	if err == nil {
		t.Errorf("Nil err")
	}
}

func TestGetFilm5(t *testing.T) {
	setEnv()
	os.Setenv("PEOPLE_FILE", "")
	_, err := getStorage(t)
	if err == nil {
		t.Errorf("Nil err")
	}
}

func TestGetFilm6(t *testing.T) {
	setEnv()
	os.Setenv("LOCATIONS_FILE", "")
	_, err := getStorage(t)
	if err == nil {
		t.Errorf("Nil err")
	}
}

func TestGetFilm7(t *testing.T) {
	setEnv()
	os.Setenv("SPECIES_FILE", "")
	_, err := getStorage(t)
	if err == nil {
		t.Errorf("Nil err")
	}
}

func TestGetFilm8(t *testing.T) {
	setEnv()
	os.Setenv("VEHICLES_FILE", "")
	_, err := getStorage(t)
	if err == nil {
		t.Errorf("Nil err")
	}
}

// =============================================================

func TestGetPersonMock1(t *testing.T) {
	stor := storage.NewMock()
	GetPersonValidId(t, stor)
}

func TestGetPersonMock2(t *testing.T) {
	stor := storage.NewMock()
	GetPersonNonexistId(t, stor)
}

func TestGetPersonMock3(t *testing.T) {
	stor := storage.NewMock()
	GetPersonInvalidId(t, stor)
}

// корректный id
func TestGetPerson1(t *testing.T) {
	setEnv()
	stor, err := getStorage(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetPersonValidId(t, stor)
}

// несуществующий id
func TestGetPerson2(t *testing.T) {
	setEnv()
	stor, err := getStorage(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetPersonNonexistId(t, stor)
}

// некорректный id
func TestGetPerson3(t *testing.T) {
	setEnv()
	stor, err := getStorage(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetPersonInvalidId(t, stor)
}

// =============================================================

func TestGetLocationMock1(t *testing.T) {
	stor := storage.NewMock()
	GetLocationValidId(t, stor)
}

func TestGetLocationMock2(t *testing.T) {
	stor := storage.NewMock()
	GetLocationNonexistId(t, stor)
}

func TestGetLocationMock3(t *testing.T) {
	stor := storage.NewMock()
	GetLocationInvalidId(t, stor)
}

// корректный id
func TestGetLocation1(t *testing.T) {
	setEnv()
	stor, err := getStorage(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetLocationValidId(t, stor)
}

// несуществующий id
func TestGetLocation2(t *testing.T) {
	setEnv()
	stor, err := getStorage(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetLocationNonexistId(t, stor)
}

// некорректный id
func TestGetLocation3(t *testing.T) {
	setEnv()
	stor, err := getStorage(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetLocationInvalidId(t, stor)
}

// =============================================================

// func TestGetSpeciesMock1(t *testing.T) {
// 	stor := stor.NewMock()
// 	GetSpeciesValidId(t, stor)
// }

// func TestGetSpeciesMock2(t *testing.T) {
// 	stor := stor.NewMock()
// 	GetSpeciesNonexistId(t, stor)
// }

// func TestGetSpeciesMock3(t *testing.T) {
// 	stor := stor.NewMock()
// 	GetSpeciesInvalidId(t, stor)
// }

// корректный id
func TestGetSpecies1(t *testing.T) {
	setEnv()
	stor, err := getStorage(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetSpeciesValidId(t, stor)
}

// несуществующий id
func TestGetSpecies2(t *testing.T) {
	setEnv()
	stor, err := getStorage(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetSpeciesNonexistId(t, stor)
}

// некорректный id
func TestGetSpecies3(t *testing.T) {
	setEnv()
	stor, err := getStorage(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetSpeciesInvalidId(t, stor)
}

// =============================================================

// func TestGetVehicleMock1(t *testing.T) {
// 	stor := stor.NewMock()
// 	GetVehicleValidId(t, stor)
// }

// func TestGetVehicleMock2(t *testing.T) {
// 	stor := stor.NewMock()
// 	GetVehicleNonexistId(t, stor)
// }

// func TestGetVehicleMock3(t *testing.T) {
// 	stor := stor.NewMock()
// 	GetVehicleInvalidId(t, stor)
// }

// корректный id
func TestGetVehicle1(t *testing.T) {
	setEnv()
	stor, err := getStorage(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetVehicleValidId(t, stor)
}

// несуществующий id
func TestGetVehicle2(t *testing.T) {
	setEnv()
	stor, err := getStorage(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetVehicleNonexistId(t, stor)
}

// некорректный id
func TestGetVehicle3(t *testing.T) {
	setEnv()
	stor, err := getStorage(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetVehicleInvalidId(t, stor)
}

// =====================================================

// ErrorStorage должен возвращать только ошибки internal.MockError
func TestErrorStorage(t *testing.T) {
	stor := storage.NewErr()
	if _, err := stor.GetFilm(utils.NewUUID()); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := stor.GetFilms(); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := stor.GetPerson(utils.NewUUID()); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := stor.GetPeople(); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := stor.GetLocation(utils.NewUUID()); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := stor.GetLocations(); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := stor.GetSpecies(utils.NewUUID()); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := stor.GetAllSpecies(); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := stor.GetVehicle(utils.NewUUID()); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := stor.GetVehicles(); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
}
