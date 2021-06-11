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

func getAnyFilm(repo *storage.Storage) (*ghibli.Film, error) {
	films, _ := repo.GetFilms()
	for _, film := range films {
		return film, nil

	}
	return nil, fmt.Errorf("Empty repo")
}

func getAnyPerson(repo *storage.Storage) (*ghibli.Person, error) {
	items, _ := repo.GetPeople()
	for _, item := range items {
		return item, nil

	}
	return nil, fmt.Errorf("Empty repo")
}

func getAnyLocation(repo *storage.Storage) (*ghibli.Location, error) {
	items, _ := repo.GetLocations()
	for _, item := range items {
		return item, nil

	}
	return nil, fmt.Errorf("Empty repo")
}

func getAnySpecies(repo *storage.Storage) (*ghibli.Species, error) {
	items, _ := repo.GetAllSpecies()
	for _, item := range items {
		return item, nil

	}
	return nil, fmt.Errorf("Empty repo")
}

func getAnyVehicle(repo *storage.Storage) (*ghibli.Vehicle, error) {
	items, _ := repo.GetVehicles()
	for _, item := range items {
		return item, nil

	}
	return nil, fmt.Errorf("Empty repo")
}

// корректный id
func GetFilmValidId(t *testing.T, repo *storage.Storage) {
	film1, err := getAnyFilm(repo)
	if err != nil {
		t.Errorf("getAnyFilmId err: %v", err)
	}
	film2, err := repo.GetFilm(film1.Id)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	if !reflect.DeepEqual(film1, film2) {
		t.Errorf("films not equal")
	}
}

// несуществующий id
func GetFilmNonexistId(t *testing.T, repo *storage.Storage) {
	_, err := repo.GetFilm(utils.NewUUID())
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// некорректный id
func GetFilmInvalidId(t *testing.T, repo *storage.Storage) {
	_, err := repo.GetFilm("")
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// корректный id
func GetPersonValidId(t *testing.T, repo *storage.Storage) {
	item1, err := getAnyPerson(repo)
	if err != nil {
		t.Errorf("getAnyPerson err: %v", err)
	}
	item2, err := repo.GetPerson(item1.Id)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	if !reflect.DeepEqual(item1, item2) {
		t.Errorf("items not equal")
	}
}

// несуществующий id
func GetPersonNonexistId(t *testing.T, repo *storage.Storage) {
	_, err := repo.GetPerson(utils.NewUUID())
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// некорректный id
func GetPersonInvalidId(t *testing.T, repo *storage.Storage) {
	_, err := repo.GetPerson("")
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// корректный id
func GetLocationValidId(t *testing.T, repo *storage.Storage) {
	item1, err := getAnyLocation(repo)
	if err != nil {
		t.Errorf("getAnyPerson err: %v", err)
	}
	item2, err := repo.GetLocation(item1.Id)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	if !reflect.DeepEqual(item1, item2) {
		t.Errorf("items not equal")
	}
}

// несуществующий id
func GetLocationNonexistId(t *testing.T, repo *storage.Storage) {
	_, err := repo.GetLocation(utils.NewUUID())
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// некорректный id
func GetLocationInvalidId(t *testing.T, repo *storage.Storage) {
	_, err := repo.GetLocation("")
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// корректный id
func GetSpeciesValidId(t *testing.T, repo *storage.Storage) {
	item1, err := getAnySpecies(repo)
	if err != nil {
		t.Errorf("getAnyPerson err: %v", err)
	}
	item2, err := repo.GetSpecies(item1.Id)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	if !reflect.DeepEqual(item1, item2) {
		t.Errorf("items not equal")
	}
}

// несуществующий id
func GetSpeciesNonexistId(t *testing.T, repo *storage.Storage) {
	_, err := repo.GetSpecies(utils.NewUUID())
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// некорректный id
func GetSpeciesInvalidId(t *testing.T, repo *storage.Storage) {
	_, err := repo.GetSpecies("")
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// корректный id
func GetVehicleValidId(t *testing.T, repo *storage.Storage) {
	item1, err := getAnyVehicle(repo)
	if err != nil {
		t.Errorf("getAnyPerson err: %v", err)
	}
	item2, err := repo.GetVehicle(item1.Id)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	if !reflect.DeepEqual(item1, item2) {
		t.Errorf("items not equal")
	}
}

// несуществующий id
func GetVehicleNonexistId(t *testing.T, repo *storage.Storage) {
	_, err := repo.GetVehicle(utils.NewUUID())
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

// некорректный id
func GetVehicleInvalidId(t *testing.T, repo *storage.Storage) {
	_, err := repo.GetVehicle("")
	if err != internal.NotFoundError {
		t.Errorf("Wrong err: %v", err)
	}
}

func TestGetFilmMock1(t *testing.T) {
	repo := storage.NewMock()
	GetFilmValidId(t, repo)
}

func TestGetFilmMock2(t *testing.T) {
	repo := storage.NewMock()
	GetFilmNonexistId(t, repo)
}

func TestGetFilmMock3(t *testing.T) {
	repo := storage.NewMock()
	GetFilmInvalidId(t, repo)
}

func setEnv() {

	log.SetFlags(log.Lshortfile)

	os.Setenv("FILMS_FILE", "./src/films.js")
	os.Setenv("PEOPLE_FILE", "./src/people.js")
	os.Setenv("LOCATIONS_FILE", "./src/locations.js")
	os.Setenv("SPECIES_FILE", "./src/species.js")
	os.Setenv("VEHICLES_FILE", "./src/vehicles.js")
}

func getRepo(t *testing.T) (*storage.Storage, error) {

	repo := storage.NewLocal()
	main := &application.MainParams{Ctx: context.Background(), Wg: new(sync.WaitGroup), Kill: func() {}}

	err := repo.Run(main)
	if err != nil {
		return nil, fmt.Errorf("repo.Run err: %v", err)
	}

	return repo, nil
}

// корректный id
func TestGetFilm1(t *testing.T) {
	setEnv()
	repo, err := getRepo(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetFilmValidId(t, repo)
}

// несуществующий id
func TestGetFilm2(t *testing.T) {
	setEnv()
	repo, err := getRepo(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetFilmNonexistId(t, repo)
}

// некорректный id
func TestGetFilm3(t *testing.T) {
	setEnv()
	repo, err := getRepo(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetFilmInvalidId(t, repo)
}

// некорректно указан файл фильмов

func TestGetFilm4(t *testing.T) {
	setEnv()
	os.Setenv("FILMS_FILE", "")
	_, err := getRepo(t)
	if err == nil {
		t.Errorf("Nil err")
	}
}

func TestGetFilm5(t *testing.T) {
	setEnv()
	os.Setenv("PEOPLE_FILE", "")
	_, err := getRepo(t)
	if err == nil {
		t.Errorf("Nil err")
	}
}

func TestGetFilm6(t *testing.T) {
	setEnv()
	os.Setenv("LOCATIONS_FILE", "")
	_, err := getRepo(t)
	if err == nil {
		t.Errorf("Nil err")
	}
}

func TestGetFilm7(t *testing.T) {
	setEnv()
	os.Setenv("SPECIES_FILE", "")
	_, err := getRepo(t)
	if err == nil {
		t.Errorf("Nil err")
	}
}

func TestGetFilm8(t *testing.T) {
	setEnv()
	os.Setenv("VEHICLES_FILE", "")
	_, err := getRepo(t)
	if err == nil {
		t.Errorf("Nil err")
	}
}

// =============================================================

func TestGetPersonMock1(t *testing.T) {
	repo := storage.NewMock()
	GetPersonValidId(t, repo)
}

func TestGetPersonMock2(t *testing.T) {
	repo := storage.NewMock()
	GetPersonNonexistId(t, repo)
}

func TestGetPersonMock3(t *testing.T) {
	repo := storage.NewMock()
	GetPersonInvalidId(t, repo)
}

// корректный id
func TestGetPerson1(t *testing.T) {
	setEnv()
	repo, err := getRepo(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetPersonValidId(t, repo)
}

// несуществующий id
func TestGetPerson2(t *testing.T) {
	setEnv()
	repo, err := getRepo(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetPersonNonexistId(t, repo)
}

// некорректный id
func TestGetPerson3(t *testing.T) {
	setEnv()
	repo, err := getRepo(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetPersonInvalidId(t, repo)
}

// =============================================================

func TestGetLocationMock1(t *testing.T) {
	repo := storage.NewMock()
	GetLocationValidId(t, repo)
}

func TestGetLocationMock2(t *testing.T) {
	repo := storage.NewMock()
	GetLocationNonexistId(t, repo)
}

func TestGetLocationMock3(t *testing.T) {
	repo := storage.NewMock()
	GetLocationInvalidId(t, repo)
}

// корректный id
func TestGetLocation1(t *testing.T) {
	setEnv()
	repo, err := getRepo(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetLocationValidId(t, repo)
}

// несуществующий id
func TestGetLocation2(t *testing.T) {
	setEnv()
	repo, err := getRepo(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetLocationNonexistId(t, repo)
}

// некорректный id
func TestGetLocation3(t *testing.T) {
	setEnv()
	repo, err := getRepo(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetLocationInvalidId(t, repo)
}

// =============================================================

// func TestGetSpeciesMock1(t *testing.T) {
// 	repo := repo.NewMock()
// 	GetSpeciesValidId(t, repo)
// }

// func TestGetSpeciesMock2(t *testing.T) {
// 	repo := repo.NewMock()
// 	GetSpeciesNonexistId(t, repo)
// }

// func TestGetSpeciesMock3(t *testing.T) {
// 	repo := repo.NewMock()
// 	GetSpeciesInvalidId(t, repo)
// }

// корректный id
func TestGetSpecies1(t *testing.T) {
	setEnv()
	repo, err := getRepo(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetSpeciesValidId(t, repo)
}

// несуществующий id
func TestGetSpecies2(t *testing.T) {
	setEnv()
	repo, err := getRepo(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetSpeciesNonexistId(t, repo)
}

// некорректный id
func TestGetSpecies3(t *testing.T) {
	setEnv()
	repo, err := getRepo(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetSpeciesInvalidId(t, repo)
}

// =============================================================

// func TestGetVehicleMock1(t *testing.T) {
// 	repo := repo.NewMock()
// 	GetVehicleValidId(t, repo)
// }

// func TestGetVehicleMock2(t *testing.T) {
// 	repo := repo.NewMock()
// 	GetVehicleNonexistId(t, repo)
// }

// func TestGetVehicleMock3(t *testing.T) {
// 	repo := repo.NewMock()
// 	GetVehicleInvalidId(t, repo)
// }

// корректный id
func TestGetVehicle1(t *testing.T) {
	setEnv()
	repo, err := getRepo(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetVehicleValidId(t, repo)
}

// несуществующий id
func TestGetVehicle2(t *testing.T) {
	setEnv()
	repo, err := getRepo(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetVehicleNonexistId(t, repo)
}

// некорректный id
func TestGetVehicle3(t *testing.T) {
	setEnv()
	repo, err := getRepo(t)
	if err != nil {
		t.Errorf("Not nil err: %v", err)
	}
	GetVehicleInvalidId(t, repo)
}

// =====================================================

// ErrorRepo должен возвращать только ошибки internal.MockError
func TestErrorRepo(t *testing.T) {
	repo := storage.NewErr()
	if _, err := repo.GetFilm(utils.NewUUID()); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := repo.GetFilms(); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := repo.GetPerson(utils.NewUUID()); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := repo.GetPeople(); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := repo.GetLocation(utils.NewUUID()); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := repo.GetLocations(); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := repo.GetSpecies(utils.NewUUID()); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := repo.GetAllSpecies(); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := repo.GetVehicle(utils.NewUUID()); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
	if _, err := repo.GetVehicles(); err != internal.MockError {
		t.Errorf("Wrong error: %v", err)
	}
}
