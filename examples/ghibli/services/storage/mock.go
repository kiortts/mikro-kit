package storage

import (
	"github.com/kiortts/mikro-kit/examples/ghibli"
	"github.com/kiortts/mikro-kit/examples/ghibli/internal"
)

// Возвращает хранилище заполненное минимальным количеством данных, необходимых для тестирвания других сервисов.
// По два объекта каждого типа.
func NewMock() *Storage {
	s := NewLocal()
	s.fillMockData()
	return s
}

// Заполнение хранилища тестовыми данными
func (s *Storage) fillMockData() {
	s.films[film1.Id] = &film1
	s.films[film2.Id] = &film2
	s.people[person1.Id] = &person1
	s.people[person2.Id] = &person2
	s.locations[location1.Id] = &location1
	s.locations[location2.Id] = &location2
	s.species[species1.Id] = &species1
	s.species[species2.Id] = &species2
	s.vehicles[vehicle1.Id] = &vehicle1
	s.vehicles[vehicle2.Id] = &vehicle2
}

var film1 ghibli.Film = ghibli.Film{
	Id:                "2baf70d1-42bb-4437-b551-e5fed5a87abe",
	Title:             "Castle in the Sky",
	Description:       "The orphan Sheeta inherited a mysterious crystal that links her to the mythical sky-kingdom of Laputa. With the help of resourceful Pazu and a rollicking band of sky pirates, she makes her way to the ruins of the once-great civilization. Sheeta and Pazu must outwit the evil Muska, who plans to use Laputa's science to make himself ruler of the world.",
	Director:          "Hayao Miyazaki",
	Producer:          "Isao Takahata",
	ReleaseDate:       "1986",
	RottenTomatoScore: "95",
}

var film2 ghibli.Film = ghibli.Film{
	Id:                "12cfb892-aac0-4c5b-94af-521852e46d6a",
	Title:             "Grave of the Fireflies",
	Description:       "In the latter part of World War II, a boy and his sister, orphaned when their mother is killed in the firebombing of Tokyo, are left to survive on their own in what remains of civilian life in Japan. The plot follows this boy and his sister as they do their best to survive in the Japanese countryside, battling hunger, prejudice, and pride in their own quiet, personal battle.",
	Director:          "Isao Takahata",
	Producer:          "Toru Hara",
	ReleaseDate:       "1988",
	RottenTomatoScore: "97",
}

var person1 ghibli.Person = ghibli.Person{
	Id:        "ba924631-068e-4436-b6de-f3283fa848f0",
	Name:      "Ashitaka",
	Gender:    "male",
	Age:       "late teens",
	EyeColor:  "brown",
	HairColor: "brown",
	Films:     []string{"https://ghiblighibli.herokuapp.com/films/030555b3-4c92-4fce-93fb-e70c3ae3df8b"},
	Species:   "https://ghiblighibli.herokuapp.com/species/af3910a6-429f-4c74-9ad5-dfe1c4aa04f2",
	Url:       "https://ghiblighibli.herokuapp.com/people/ba924631-068e-4436-b6de-f3283fa848f0",
}

var person2 ghibli.Person = ghibli.Person{
	Id:        "030555b3-4c92-4fce-93fb-e70c3ae3df8b",
	Name:      "Yakul",
	Gender:    "male",
	Age:       "Unknown",
	EyeColor:  "Grey",
	HairColor: "Brown",
	Films:     []string{"https://ghiblighibli.herokuapp.com/films/0440483e-ca0e-4120-8c50-4c8cd9b965d6"},
	Species:   "https://ghiblighibli.herokuapp.com/species/6bc92fdd-b0f4-4286-ad71-1f99fb4a0d1e",
	Url:       "https://ghiblighibli.herokuapp.com/people/030555b3-4c92-4fce-93fb-e70c3ae3df8b",
}

var location1 ghibli.Location = ghibli.Location{
	Id:           "11014596-71b0-4b3e-b8c0-1c4b15f28b9a",
	Name:         "Irontown",
	Climate:      "Continental",
	Terrain:      "Mountain",
	SurfaceWater: "40",
	Residents: []string{"https://ghiblighibli.herokuapp.com/people/ba924631-068e-4436-b6de-f3283fa848f0",
		"https://ghiblighibli.herokuapp.com/people/030555b3-4c92-4fce-93fb-e70c3ae3df8b"},
	Films: []string{"https://ghiblighibli.herokuapp.com/films/0440483e-ca0e-4120-8c50-4c8cd9b965d6"},
	Url:   []string{"https://ghiblighibli.herokuapp.com/locations/11014596-71b0-4b3e-b8c0-1c4b15f28b9a"},
}

var location2 ghibli.Location = ghibli.Location{
	Id:           "11014596-71b0-4b3e-b8c0-1c4b15f28b9a",
	Name:         "Gutiokipanja",
	Climate:      "Continental",
	Terrain:      "Hill",
	SurfaceWater: "50",
	Residents: []string{"https://ghiblighibli.herokuapp.com/people/ba924631-068e-4436-b6de-f3283fa848f0",
		"https://ghiblighibli.herokuapp.com/people/030555b3-4c92-4fce-93fb-e70c3ae3df8b"},
	Films: []string{"https://ghiblighibli.herokuapp.com/films/0440483e-ca0e-4120-8c50-4c8cd9b965d6"},
	Url:   []string{"https://ghiblighibli.herokuapp.com/locations/11014596-71b0-4b3e-b8c0-1c4b15f28b9a"},
}

var species1 ghibli.Species = ghibli.Species{
	Id: "",
}

var species2 ghibli.Species = ghibli.Species{
	Id: "",
}

var vehicle1 ghibli.Vehicle = ghibli.Vehicle{
	Id: "",
}

var vehicle2 ghibli.Vehicle = ghibli.Vehicle{
	Id: "",
}

// хранилище, возвращающее только ошибки
type errRepo struct{}

// реализация интерфейсов всех repo
func (r *errRepo) GetFilm(id string) (*ghibli.Film, error)         { return nil, internal.MockError }
func (r *errRepo) GetFilms() ([]*ghibli.Film, error)               { return nil, internal.MockError }
func (r *errRepo) GetPerson(id string) (*ghibli.Person, error)     { return nil, internal.MockError }
func (r *errRepo) GetPeople() ([]*ghibli.Person, error)            { return nil, internal.MockError }
func (r *errRepo) GetLocation(id string) (*ghibli.Location, error) { return nil, internal.MockError }
func (r *errRepo) GetLocations() ([]*ghibli.Location, error)       { return nil, internal.MockError }
func (r *errRepo) GetSpecies(id string) (*ghibli.Species, error)   { return nil, internal.MockError }
func (r *errRepo) GetAllSpecies() ([]*ghibli.Species, error)       { return nil, internal.MockError }
func (r *errRepo) GetVehicle(id string) (*ghibli.Vehicle, error)   { return nil, internal.MockError }
func (r *errRepo) GetVehicles() ([]*ghibli.Vehicle, error)         { return nil, internal.MockError }

// Возвращает хранилище, имитирующее ошибки при работе.
func NewErr() *errRepo {
	return new(errRepo)
}
