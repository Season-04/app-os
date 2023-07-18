package manifest

type Manifest struct {
	ID      string           `json:"id"`
	Name    string           `json:"name"`
	Image   string           `json:"image"`
	Schemas Schemas          `json:"schemas"`
	Routes  map[string]Route `json:"routes"`
}

type Schemas struct {
	Provides     []ProvidedSchema `json:"provides"`
	Dependencies []string         `json:"dependencies"`
}

type ProvidedSchema struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Port       uint16 `json:"port"`
	Definition string `json:"definition"`
}

type Route struct {
	Port uint16 `json:"port"`
}