/*
В пакете описаны общие для сервиса публичные типы и функции, предоставляемые как внешний API.
Типы и функции, общие для сервиса, но не предоставляемые публично описаны в пакете internal.
*/

package ghibli

// Film - all of the Studio Ghibli films.
type Film struct {
	Id                string `json:"id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Director          string `json:"director"`
	Producer          string `json:"producer"`
	ReleaseDate       string `json:"release_date"`
	RottenTomatoScore string `json:"rt_score"`
}

// Person includes all Ghibli characters, human and non-.
type Person struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Gender    string   `json:"gender"`
	Age       string   `json:"age"`
	EyeColor  string   `json:"eye_color"`
	HairColor string   `json:"hair_color"`
	Films     []string `json:"films"`
	Species   string   `json:"species"`
	Url       string   `json:"url"`
}

// Location includes lands, countries, and places.
type Location struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Climate      string   `json:"climate"`
	Terrain      string   `json:"terrain"`
	SurfaceWater string   `json:"surface_water"`
	Residents    []string `json:"residents"`
	Films        []string `json:"films"`
	Url          []string `json:"url"`
}

// Species includes humans, animals, and spirits et al.
type Species struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Classification string   `json:"classification"`
	EyeColors      string   `json:"eye_colors"`
	HairColors     string   `json:"hair_colors"`
	Url            string   `json:"url"`
	People         []string `json:"people"`
	Films          []string `json:"films"`
}

// Vehicle includes cars, ships, and planes.
type Vehicle struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	VehicleClass string `json:"vehicle_class"`
	Length       string `json:"length"`
	Pilot        string `json:"pilot"`
	Films        string `json:"films"`
	Url          string `json:"url"`
}

type FilmStorage interface {
	GetFilm(id string) (*Film, error) // Returns a film based on a single ID
	GetFilms() ([]*Film, error)       // Returns information about all of the Studio Ghibli films.
}

type PersonStorage interface {
	GetPerson(id string) (*Person, error)
	GetPeople() ([]*Person, error)
}

type LocationStorage interface {
	GetLocation(id string) (*Location, error)
	GetLocations() ([]*Location, error)
}

type SpeciesStorage interface {
	GetSpecies(id string) (*Species, error)
	GetAllSpecies() ([]*Species, error)
}

type VehicleStorage interface {
	GetVehicle(id string) (*Vehicle, error)
	GetVehicles() ([]*Vehicle, error)
}
