package parsers

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/IfisUASD/ifisuasd.github.io/internal/types"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
)

// --- Prototipos ---

// ParsePerson procesa el contenido crudo de un archivo Markdown de perfil.
// Extrae el YAML del encabezado y convierte el cuerpo a HTML.
func ParsePerson(filename string, content []byte) (*types.Person, error) {
	var p types.Person
	html, err := parseMarkdownAndYAML(content, &p)
	if err != nil {
		return nil, fmt.Errorf("error parsing markdown: %w", err)
	}

	p.BioHTML = html
	p.Slug = cleanSlug(filename)

	// Validaciones básicas
	if p.ID == "" {
		return nil, fmt.Errorf("missing required field: orcid")
	}
	if p.Name == "" {
		return nil, fmt.Errorf("missing required field: name")
	}

	return &p, nil
}

// ParseProject procesa el contenido crudo de un archivo Markdown de proyecto.
func ParseProject(filename string, content []byte) (*types.Project, error) {
	// Estructura temporal para manejar el parsing de fechas como strings
	type ProjectYAML struct {
		ID                      string   `yaml:"project_id"`
		Title                   string   `yaml:"title"`
		Status                  string   `yaml:"status"`
		Funding                 string   `yaml:"funding"`
		Tags                    []string `yaml:"tags"`
		StartDate               string   `yaml:"start_date"`
		EndDate                 string   `yaml:"end_date"`
		PrincipalInvestigatorID string   `yaml:"principal_investigator"`
		CoinvestigatorIDs       []string `yaml:"coinvestigator"`
		ResearchAssistantIDs    []string `yaml:"research_assistant"`
	}

	var py ProjectYAML
	html, err := parseMarkdownAndYAML(content, &py)
	if err != nil {
		return nil, fmt.Errorf("error parsing markdown: %w", err)
	}

	p := &types.Project{
		ID:                      py.ID,
		Slug:                    cleanSlug(filename),
		Title:                   py.Title,
		Status:                  py.Status,
		Funding:                 py.Funding,
		Tags:                    py.Tags,
		PrincipalInvestigatorID: py.PrincipalInvestigatorID,
		CoinvestigatorIDs:       py.CoinvestigatorIDs,
		ResearchAssistantIDs:    py.ResearchAssistantIDs,
		DescriptionHTML:         html,
	}

	if py.StartDate != "" {
		t, err := time.Parse("2006-01-02", py.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date format: %w", err)
		}
		p.StartDate = t
	}

	if py.EndDate != "" {
		t, err := time.Parse("2006-01-02", py.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format: %w", err)
		}
		p.EndDate = t
	}

	if p.ID == "" {
		return nil, fmt.Errorf("missing required field: project_id")
	}

	return p, nil
}

// ParseNewsItem procesa una noticia.
func ParseNewsItem(filename string, content []byte) (*types.NewsItem, error) {
	// Estructura temporal para fechas
	type NewsYAML struct {
		ID      string `yaml:"id"`
		Title   string `yaml:"title"`
		Date    string `yaml:"date"`
		Summary string `yaml:"summary"`
		Image   string `yaml:"image"`
		ImageAlt string `yaml:"image_alt"`
	}

	var ny NewsYAML
	html, err := parseMarkdownAndYAML(content, &ny)
	if err != nil {
		return nil, fmt.Errorf("error parsing markdown: %w", err)
	}

	n := &types.NewsItem{
		ID:          ny.ID,
		Slug:        cleanSlug(filename),
		Title:       ny.Title,
		Summary:     ny.Summary,
		Image:       ny.Image,
		ImageAlt:    ny.ImageAlt,
		ContentHTML: html,
	}

	if ny.Date != "" {
		t, err := time.Parse("2006-01-02", ny.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}
		n.Date = t
	}

	if n.ID == "" {
		// Fallback: usar slug como ID si no hay ID explícito
		n.ID = n.Slug
	}

	return n, nil
}

// ParseBlogPost procesa una entrada de blog.
func ParseBlogPost(filename string, content []byte) (*types.BlogPost, error) {
	type BlogYAML struct {
		ID       string   `yaml:"id"`
		Title    string   `yaml:"title"`
		Date     string   `yaml:"date"`
		AuthorID string   `yaml:"author_id"`
		Tags     []string `yaml:"tags"`
	}

	var by BlogYAML
	html, err := parseMarkdownAndYAML(content, &by)
	if err != nil {
		return nil, fmt.Errorf("error parsing markdown: %w", err)
	}

	b := &types.BlogPost{
		ID:          by.ID,
		Slug:        cleanSlug(filename),
		Title:       by.Title,
		AuthorID:    by.AuthorID,
		Tags:        by.Tags,
		ContentHTML: html,
	}

	if by.Date != "" {
		t, err := time.Parse("2006-01-02", by.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}
		b.Date = t
	}

	if b.ID == "" {
		b.ID = b.Slug
	}

	return b, nil
}

// --- Helper interno ---
func parseMarkdownAndYAML(content []byte, target interface{}) (template.HTML, error) {
	var buf bytes.Buffer
	
	// Configuración del parser con frontmatter
	fm := &frontmatter.Extender{
		Mode: frontmatter.SetMetadata,
	}
	
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			LatexMath,
			fm,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	
	// Contexto para extraer el frontmatter
	ctx := parser.NewContext()
	
	if err := md.Convert(content, &buf, parser.WithContext(ctx)); err != nil {
		return "", err
	}
	
	// Decodificar el frontmatter en el target
	d := frontmatter.Get(ctx)
	if d == nil {
		return "", fmt.Errorf("no frontmatter found")
	}
	
	if err := d.Decode(target); err != nil {
		return "", fmt.Errorf("error decoding frontmatter: %w", err)
	}
	
	return template.HTML(buf.String()), nil
}

func cleanSlug(filename string) string {
	base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	base = strings.TrimSuffix(base, ".es")
	base = strings.TrimSuffix(base, ".en")
	return base
}


// ParseTool procesa la definición de una herramienta.
func ParseTool(filename string, content []byte) (*types.Tool, error) {
	var t types.Tool
	// Reutilizamos la lógica de frontmatter, ignorando el contenido HTML por ahora
	// ya que la 'landing' de apps suele ser solo tarjetas.
	_, err := parseMarkdownAndYAML(content, &t)
	if err != nil {
		return nil, fmt.Errorf("error parsing tool markdown: %w", err)
	}
	
	if t.ID == "" {
		t.ID = cleanSlug(filename)
	}
	return &t, nil
}