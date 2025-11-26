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
	AvatarAlt string            `yaml:"avatar_alt"` // Alt text for accessibility
	Social map[string]string `yaml:"social"` // { "scholar": "...", "twitter": "..." }

	// Contenido
	BioHTML template.HTML // El cuerpo del Markdown convertido a HTML

	// Relaciones
	Projects     []*Project
	Publications []*Publication // Obras donde es AUTOR
	Mentored     []*Publication // Obras donde es ASESOR (Tesis dirigidas)
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
	Tags      []string  `yaml:"tags"`       // Etiquetas adicionales
	StartDate time.Time `yaml:"start_date"` // YYYY-MM-DD
	EndDate   time.Time `yaml:"end_date"`

	// IDs para vinculación (Leídos del YAML)
	PrincipalInvestigatorID string   `yaml:"principal_investigator"`
	CoinvestigatorIDs        []string `yaml:"coinvestigator"`
	ResearchAssistantIDs     []string `yaml:"research_assistant"`

	// Contenido
	DescriptionHTML template.HTML

	// Relaciones (Calculadas)
	Publications          []*Publication
	PrincipalInvestigator *Person
	Coinvestigators       []*Person
	ResearchAssistants    []*Person
}

type Publication struct {
	// Identificadores
	ID   string // Citation Key
	Slug string
	DOI  string

	// Campos BibTeX
	Type      string // article, phdthesis, mastersthesis, etc.
	Title     string
	Authors   []string // Strings crudos (ej: "Juan Pérez")
	Journal   string
	Year      int
	Volume    string
	Number    string
	Pages     string
	Publisher string
	School    string // Importante para Tesis
	Booktitle string
	Abstract  string
	URL       string

	// Campos de Vinculación (x-fields)
	AuthorOrcids  []string // x-orcids
	AdvisorOrcids []string // x-advisors (NUEVO)
	ProjectID     string   // x-project

	// Relaciones (Punteros)
	LinkedAuthors  []*Person
	LinkedAdvisors []*Person // Investigadores locales que asesoraron
	Project        *Project
}

// NewsItem representa una noticia o evento.
type NewsItem struct {
	ID      string    `yaml:"id"`
	Slug    string
	Title   string    `yaml:"title"`
	Date    time.Time `yaml:"date"`
	Summary string    `yaml:"summary"`
	Image   string    `yaml:"image"`
	ImageAlt string   `yaml:"image_alt"` // Alt text for accessibility
	
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
	Publications []*Publication
	News      []*NewsItem
	BlogPosts []*BlogPost
	Tools     []*Tool
}

// NewDatabase inicializa los mapas para evitar nil pointer panics.
func NewDatabase() *Database {
	return &Database{
		People:    make(map[string]*Person),
		Projects:  make(map[string]*Project),
		Publications: make([]*Publication, 0),
		News:      make([]*NewsItem, 0),
		BlogPosts: make([]*BlogPost, 0),
		Tools:     make([]*Tool, 0),
	}
}

// Tool representa una aplicación o herramienta web.
type Tool struct {
	ID          string `yaml:"id"`
	Title       string `yaml:"title"`
	Summary     string `yaml:"summary"`
	Link        string `yaml:"link"`        // URL relativa, ej: "/apps/qr"
	Icon        string `yaml:"icon"`        // Nombre de icono (para usar en SVG) o ruta de img
	ButtonText  string `yaml:"button_text"` // Ej: "Abrir Generador"
}