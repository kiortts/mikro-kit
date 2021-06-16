package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/kiortts/mikro-kit/components"
	"github.com/kiortts/mikro-kit/examples/ghibli"
	"github.com/kiortts/mikro-kit/examples/ghibli/internal"
	"github.com/pkg/errors"
)

// Storage сервис
type Storage struct {
	films     map[string]*ghibli.Film
	people    map[string]*ghibli.Person
	locations map[string]*ghibli.Location
	species   map[string]*ghibli.Species
	vehicles  map[string]*ghibli.Vehicle
}

// проверка реализации типом интерфейсов
var _ ghibli.FilmStorage = (*Storage)(nil)
var _ ghibli.PersonStorage = (*Storage)(nil)
var _ ghibli.LocationStorage = (*Storage)(nil)
var _ ghibli.SpeciesStorage = (*Storage)(nil)
var _ ghibli.VehicleStorage = (*Storage)(nil)

// конфигурация
type config struct {
	FilmsFile     string `arg:"env:FILMS_FILE"`
	PeopleFile    string `arg:"env:PEOPLE_FILE"`
	LocationsFile string `arg:"env:LOCATIONS_FILE"`
	SpeciesFile   string `arg:"env:SPECIES_FILE"`
	VehiclesFile  string `arg:"env:VEHICLES_FILE"`
}

func getConfig() *config {
	cfg := &config{}
	cfg.FilmsFile = os.Getenv("FILMS_FILE")
	cfg.PeopleFile = os.Getenv("PEOPLE_FILE")
	cfg.LocationsFile = os.Getenv("LOCATIONS_FILE")
	cfg.SpeciesFile = os.Getenv("SPECIES_FILE")
	cfg.VehiclesFile = os.Getenv("VEHICLES_FILE")
	return cfg
}

// NewLocal возвращает сервис
func NewLocal() *Storage {

	r := &Storage{
		films:     make(map[string]*ghibli.Film),
		people:    make(map[string]*ghibli.Person),
		locations: make(map[string]*ghibli.Location),
		species:   make(map[string]*ghibli.Species),
		vehicles:  make(map[string]*ghibli.Vehicle),
	}

	return r
}

func (s *Storage) Run(main *components.MainParams) error {

	// log.Printf(dict.LOG_SERVICE_RUN, "Storage")

	if err := s.fill(); err != nil {
		return errors.Wrap(err, "Run Storage err")
	}

	main.Wg.Add(1)
	go func() {
		defer main.Wg.Done()
		<-main.Ctx.Done()
		log.Printf("Storage DONE")
	}()

	return nil
}

// GetFilm реализация интерфейса FilmStorage.
// Возвращает один фильм по id.
func (s *Storage) GetFilm(id string) (*ghibli.Film, error) {

	film, exist := s.films[id]
	if !exist {
		return nil, internal.NotFoundError
	}

	return film, nil
}

// GetFilms реализация интерфейса FilmStorage
// Возвращает массив всех фильмов.
func (s *Storage) GetFilms() ([]*ghibli.Film, error) {

	films := []*ghibli.Film{}

	for _, film := range s.films {
		films = append(films, film)
	}

	return films, nil
}

// GetPerson реализация интерфейса PersonStorage.
// Возвращает одного персонажа фильм по id.
func (s *Storage) GetPerson(id string) (*ghibli.Person, error) {

	person, exist := s.people[id]
	if !exist {
		return nil, internal.NotFoundError
	}

	return person, nil
}

// GetPeople реализация интерфейса PersonStorage
// Возвращает массив всех персонажей.
func (s *Storage) GetPeople() ([]*ghibli.Person, error) {

	people := []*ghibli.Person{}

	for _, person := range s.people {
		people = append(people, person)
	}

	return people, nil
}

// GetLocation реализация интерфейса PersonStorage.
// Возвращает одну локацию.
func (s *Storage) GetLocation(id string) (*ghibli.Location, error) {

	item, exist := s.locations[id]
	if !exist {
		return nil, internal.NotFoundError
	}

	return item, nil
}

// GetLocations реализация интерфейса PersonStorage
// Возвращает массив всех локаций.
func (s *Storage) GetLocations() ([]*ghibli.Location, error) {

	items := []*ghibli.Location{}

	for _, item := range s.locations {
		items = append(items, item)
	}

	return items, nil
}

// GetLocation реализация интерфейса PersonStorage.
// Возвращает одну локацию.
func (s *Storage) GetSpecies(id string) (*ghibli.Species, error) {

	item, exist := s.species[id]
	if !exist {
		return nil, internal.NotFoundError
	}

	return item, nil
}

// GetLocations реализация интерфейса PersonStorage
// Возвращает массив всех локаций.
func (s *Storage) GetAllSpecies() ([]*ghibli.Species, error) {

	items := []*ghibli.Species{}

	for _, item := range s.species {
		items = append(items, item)
	}

	return items, nil
}

// GetLocation реализация интерфейса PersonStorage.
// Возвращает одну локацию.
func (s *Storage) GetVehicle(id string) (*ghibli.Vehicle, error) {

	item, exist := s.vehicles[id]
	if !exist {
		return nil, internal.NotFoundError
	}

	return item, nil
}

// GetLocations реализация интерфейса PersonStorage
// Возвращает массив всех локаций.
func (s *Storage) GetVehicles() ([]*ghibli.Vehicle, error) {

	items := []*ghibli.Vehicle{}

	for _, item := range s.vehicles {
		items = append(items, item)
	}

	return items, nil
}

// заполнение хранилища данными
func (s *Storage) fill() error {

	cfg := getConfig()

	// заполнение коллекции фильмов
	if err := s.fillFilms(cfg.FilmsFile); err != nil {
		return errors.Wrap(err, "fillFilms err")
	}

	// заполнение коллекции персонажей
	if err := s.fillPeople(cfg.PeopleFile); err != nil {
		return errors.Wrap(err, "fillPeople err")
	}

	// заполнение коллекции мест
	if err := s.fillLocations(cfg.LocationsFile); err != nil {
		return errors.Wrap(err, "fillLocations err")
	}

	// заполнение коллекции разновидностей
	if err := s.fillSpecies(cfg.SpeciesFile); err != nil {
		return errors.Wrap(err, "fillSpecies err")
	}

	// заполнение коллекции транспорта
	if err := s.fillVehicles(cfg.VehiclesFile); err != nil {
		return errors.Wrap(err, "vehicles err")
	}

	return nil
}

// заполнение коллекции фильмов
func (s *Storage) fillFilms(file string) error {

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "ReadFile err")
	}

	var items []*ghibli.Film
	err = json.Unmarshal(dat, &items)
	if err != nil {
		return errors.Wrap(err, "Unmarshal err")
	}

	for _, item := range items {
		s.films[item.Id] = item
	}

	return nil
}

// заполнение коллекции персонажей
func (s *Storage) fillPeople(file string) error {

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "ReadFile err")
	}

	var items []*ghibli.Person
	err = json.Unmarshal(dat, &items)
	if err != nil {
		return errors.Wrap(err, "Unmarshal err")
	}

	for _, item := range items {
		s.people[item.Id] = item
	}

	return nil
}

// заполнение коллекции мест
func (s *Storage) fillLocations(file string) error {

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "ReadFile err")
	}

	var items []*ghibli.Location
	err = json.Unmarshal(dat, &items)
	if err != nil {
		return errors.Wrap(err, "Unmarshal err")
	}

	for _, item := range items {
		s.locations[item.Id] = item
	}

	return nil
}

// заполнение коллекции разновидностей
func (s *Storage) fillSpecies(file string) error {

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "ReadFile err")
	}

	var items []*ghibli.Species
	err = json.Unmarshal(dat, &items)
	if err != nil {
		return errors.Wrap(err, "Unmarshal err")
	}

	for _, item := range items {
		s.species[item.Id] = item
	}

	return nil
}

// заполнение коллекции транспорта
func (s *Storage) fillVehicles(file string) error {

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "ReadFile err")
	}

	var items []*ghibli.Vehicle
	err = json.Unmarshal(dat, &items)
	if err != nil {
		return errors.Wrap(err, "Unmarshal err")
	}

	for _, item := range items {
		s.vehicles[item.Id] = item
	}

	return nil
}
