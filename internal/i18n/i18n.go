package i18n

// Dictionary holds the translations for a specific language.
type Dictionary map[string]string

var es = Dictionary{
	"Home":            "Inicio",
	"People":          "Gente",
	"Projects":        "Proyectos",
	"News":            "Noticias",
	"Blog":            "Blog",
	"ReadMore":        "Leer más",
	"AboutAuthor":     "Sobre el Autor",
	"Next":            "Siguiente",
	"Previous":        "Anterior",
	"BackToNews":      "← Volver a Noticias",
	"BackToBlog":      "← Volver al Blog",
	"ViewProfile":     "Ver Perfil",
	"PrincipalInvestigator": "Investigador Principal",
	"Coinvestigators": "Co-Investigadores",
	"ResearchAssistants": "Asistentes de Investigación",
	"RelatedPublications": "Publicaciones Relacionadas",
	"Description":     "Descripción",
	"Biography":       "Biografía",
	"Publications":    "Publicaciones",
	"Funding":         "Financiamiento",
	"Team":            "Equipo de Investigación",
	"LatestNews":      "Últimas Noticias",
	"RecentPapers":    "Publicaciones Recientes",
	"InstituteOfPhysics": "Instituto de Física",
	"UASD":            "UASD",
	"Search":            "Buscar",
	"SearchPlaceholder": "Escribe para buscar...",
}

var en = Dictionary{
	"Home":            "Home",
	"People":          "People",
	"Projects":        "Projects",
	"News":            "News",
	"Blog":            "Blog",
	"ReadMore":        "Read More",
	"AboutAuthor":     "About the Author",
	"Next":            "Next",
	"Previous":        "Previous",
	"BackToNews":      "← Back to News",
	"BackToBlog":      "← Back to Blog",
	"ViewProfile":     "View Profile",
	"PrincipalInvestigator": "Principal Investigator",
	"Coinvestigators": "Co-Investigators",
	"ResearchAssistants": "Research Assistants",
	"RelatedPublications": "Related Publications",
	"Description":     "Description",
	"Biography":       "Biography",
	"Publications":    "Publications",
	"Funding":         "Funding",
	"Team":            "Research Team",
	"LatestNews":      "Latest News",
	"RecentPapers":    "Recent Publications",
	"InstituteOfPhysics": "Institute of Physics",
	"UASD":            "UASD",
}

// GetDictionary returns the dictionary for the requested language.
func GetDictionary(lang string) Dictionary {
	if lang == "en" {
		return en
	}
	return es
}
