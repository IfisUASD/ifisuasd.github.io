package types

import (
	"html/template"
	"time"
)

// --- Entidades Principales ---

// Person representa a un investigador, administrativo o estudiante.
type Person struct {
	// Identificadores
	ID    string `yaml:"orcid"` // El ORCID es el ID principal
	Slug  string // El nombre del archivo (ej: vladimir-perez) para la URL

	// Metadatos (YAML)
	Name   string            `yaml:"name"`
	Role   string            `yaml:"role"`
	Type   string            `yaml:"type"` // academic, staff, student
	Email  string            `yaml:"email"`
	Avatar string            `yaml:"avatar"`
	Social map[string]string `yaml:"social"` // { "scholar": "...", "twitter": "..." }

	// Contenido
	BioHTML template.HTML // El cuerpo del Markdown convertido a HTML

	// Relaciones (Calculadas en tiempo de ejecución)
	Projects     []*Project
	Publications []*Paper
	Theses       []*Paper // Tesis dirigidas
}

// Project representa un proyecto de investigación.
type Project struct {
	// Identificadores
	ID   string `yaml:"project_id"` // Ej: FONDOCYT-2024-10
	Slug string // Ej: red-monitoreo-aire

	// Metadatos (YAML)
	Title     string    `yaml:"title"`
	Status    string    `yaml:"status"`     // active, finished
	Funding   string    `yaml:"funding"`    // FONDOCYT, MESCYT, UASD
	StartDate time.Time `yaml:"start_date"` // YYYY-MM-DD
	EndDate   time.Time `yaml:"end_date"`

	// IDs para vinculación (Leídos del YAML)
	PrincipalInvestigatorID string   `yaml:"principal_investigator"`
	CoinvestigatorIDs        []string `yaml:"coinvestigator"`
	ResearchAssistantIDs     []string `yaml:"research_assistant"`

	// Contenido
	DescriptionHTML template.HTML

	// Relaciones (Calculadas)
	Publications          []*Paper
	PrincipalInvestigator *Person
	Coinvestigators       []*Person
	ResearchAssistants    []*Person
}

// Paper representa una publicación académica (desde BibTeX).
type Paper struct {
	// Identificadores
	ID  string // DOI o Citation Key
	DOI string

	// BibTeX Standard Fields
	Type    string // article, phdthesis, inproceedings
	Title   string
	Authors []string // Strings crudos del BibTeX (ej: "Pérez, V. and Montero, E.")
	Journal string
	Year    int
	Volume  string
	URL     string

	// Custom Fields (x-fields para vinculación)
	AuthorOrcids []string // x-orcids
	ProjectID    string   // x-project
}

// NewsItem representa una noticia o evento.
type NewsItem struct {
	ID      string    `yaml:"id"`
	Slug    string
	Title   string    `yaml:"title"`
	Date    time.Time `yaml:"date"`
	Summary string    `yaml:"summary"`
	Image   string    `yaml:"image"`
	
	ContentHTML template.HTML
}

// BlogPost representa una entrada del blog.
type BlogPost struct {
	ID       string    `yaml:"id"`
	Slug     string
	Title    string    `yaml:"title"`
	Date     time.Time `yaml:"date"`
	AuthorID string    `yaml:"author_id"`
	Tags     []string  `yaml:"tags"`
	
	ContentHTML template.HTML
	
	// Relaciones
	Author *Person
}

// Database es el contenedor global de todos los datos en memoria.
type Database struct {
	People    map[string]*Person  // Key: ORCID
	Projects  map[string]*Project // Key: ProjectID
	Papers    []*Paper
	News      []*NewsItem
	BlogPosts []*BlogPost
}

// NewDatabase inicializa los mapas para evitar nil pointer panics.
func NewDatabase() *Database {
	return &Database{
		People:    make(map[string]*Person),
		Projects:  make(map[string]*Project),
		Papers:    make([]*Paper, 0),
		News:      make([]*NewsItem, 0),
		BlogPosts: make([]*BlogPost, 0),
	}
}