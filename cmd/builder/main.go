package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/IfisUASD/ifisuasd.github.io/internal/i18n"
	"github.com/IfisUASD/ifisuasd.github.io/internal/linker"
	"github.com/IfisUASD/ifisuasd.github.io/internal/loader"
	"github.com/IfisUASD/ifisuasd.github.io/internal/types"
	"github.com/IfisUASD/ifisuasd.github.io/templates/layouts"
	"github.com/IfisUASD/ifisuasd.github.io/templates/pages"
)

var WarningCount int
var ErrorList []string


func main() {
	log.Println("🏗️  Iniciando generador del sitio...")

	// Verificar traducciones faltantes
	checkMissingTranslations("./content")

	// Generar sitio en Español (Default)
	if err := buildSite("es", "./output"); err != nil {
		log.Fatalf("❌ Error generando sitio en Español: %v", err)
	}

	// Generar sitio en Inglés
	if err := buildSite("en", "./output/en"); err != nil {
		log.Fatalf("❌ Error generando sitio en Inglés: %v", err)
	}

	log.Println("------------------------------------------------")
    log.Println("✅ Generación completa.")
    if WarningCount > 0 {
        log.Printf("⚠️  Se encontraron %d advertencias que requieren atención:\n", WarningCount)
        for _, err := range ErrorList {
            fmt.Println("   - " + err)
        }
    } else {
        log.Println("✨  Compilación limpia: 0 advertencias.")
    }

	log.Println("✅ Generación completa.")
}

func buildSite(lang string, outputDir string) error {
	log.Printf("🌍 Generando sitio para idioma: %s en %s", lang, outputDir)

	// 1. Cargar Contenido
	db, err := loader.LoadContent("./content", lang)
	if err != nil {
		return err
	}
	log.Printf("✅ Contenido cargado: %d Personas, %d Proyectos, %d publications", len(db.People), len(db.Projects), len(db.Publications))

	// 2. Vincular Datos
	linker.LinkData(db)

	// 3. Preparar Directorio de Salida
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Obtener diccionario
	dict := i18n.GetDictionary(lang)

	// 4. Generar Home Page
	latestNews := db.News
	sort.Slice(latestNews, func(i, j int) bool {
		return latestNews[i].Date.After(latestNews[j].Date)
	})
	if len(latestNews) > 3 {
		latestNews = latestNews[:3]
	}

	recentpublications := db.Publications
	sort.Slice(recentpublications, func(i, j int) bool {
		return recentpublications[i].Year > recentpublications[j].Year
	})
	if len(recentpublications) > 5 {
		recentpublications = recentpublications[:5]
	}

	homeData := pages.HomeData{
		Meta: layouts.MetaTags{
			Description: "Sitio oficial del Instituto de Física de la UASD",
			Keywords:    "Física, UASD, Investigación, Ciencia",
		},
		LatestNews:   latestNews,
		RecentPublications: recentpublications,
	}

	f, err := os.Create(outputDir + "/index.html")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := pages.Home(homeData, lang, dict).Render(context.Background(), f); err != nil {
		return err
	}

	// 5. Generar People Page
	peopleDir := outputDir + "/people"
	if err := os.MkdirAll(peopleDir, 0755); err != nil {
		return err
	}

	var peopleList []*types.Person
	for _, p := range db.People {
		peopleList = append(peopleList, p)
	}
	sort.Slice(peopleList, func(i, j int) bool {
		return peopleList[i].Name < peopleList[j].Name
	})

	peopleData := pages.PeopleData{
		Meta: layouts.MetaTags{
			Description: "Equipo del Instituto de Física de la UASD",
			Keywords:    "Investigadores, Profesores, Estudiantes, UASD",
		},
		People: peopleList,
	}

	fPeople, err := os.Create(peopleDir + "/index.html")
	if err != nil {
		return err
	}
	defer fPeople.Close()

	if err := pages.People(peopleData, lang, dict).Render(context.Background(), fPeople); err != nil {
		return err
	}

	// 5.1 Generar Páginas Individuales de Personas
	for _, p := range db.People {
		pDir := outputDir + "/people/" + p.Slug
		if err := os.MkdirAll(pDir, 0755); err != nil {
			continue
		}

		f, err := os.Create(pDir + "/index.html")
		if err != nil {
			continue
		}
		
		pData := pages.PersonData{
			Meta: layouts.MetaTags{
				Description: p.Name + " - " + p.Role,
				Keywords:    p.Name + ", " + p.Role + ", UASD",
				OGTitle:     p.Name,
				OGDesc:      p.Role,
				OGImage:     p.Avatar,
			},
			Person: p,
		}

		if p.Avatar != "" && p.AvatarAlt == "" {
			WarningCount++
			ErrorList = append(ErrorList, fmt.Sprintf("Persona %s tiene Avatar pero falta 'avatar_alt'", p.Name))
		}

		if err := pages.Person(pData, lang, dict).Render(context.Background(), f); err != nil {
			log.Printf("❌ Error renderizando persona %s: %v", p.Slug, err)
		}
		f.Close()
	}

	// 6. Generar Projects Page
	projectsDir := outputDir + "/projects"
	if err := os.MkdirAll(projectsDir, 0755); err != nil {
		return err
	}

	var projectsList []*types.Project
	for _, p := range db.Projects {
		projectsList = append(projectsList, p)
	}
	sort.Slice(projectsList, func(i, j int) bool {
		if !projectsList[i].StartDate.IsZero() && !projectsList[j].StartDate.IsZero() {
			return projectsList[i].StartDate.After(projectsList[j].StartDate)
		}
		return projectsList[i].Title < projectsList[j].Title
	})

	projectsData := pages.ProjectsData{
		Meta: layouts.MetaTags{
			Description: "Proyectos del Instituto de Física de la UASD",
			Keywords:    "Proyectos, Investigación, UASD",
		},
		Projects: projectsList,
	}

	fProjects, err := os.Create(projectsDir + "/index.html")
	if err != nil {
		return err
	}
	defer fProjects.Close()

	if err := pages.Projects(projectsData, lang, dict).Render(context.Background(), fProjects); err != nil {
		return err
	}

	// 6.1 Generar Páginas Individuales de Proyectos
	for _, p := range db.Projects {
		pDir := outputDir + "/projects/" + p.Slug
		if err := os.MkdirAll(pDir, 0755); err != nil {
			continue
		}

		f, err := os.Create(pDir + "/index.html")
		if err != nil {
			continue
		}

		pData := pages.ProjectData{
			Meta: layouts.MetaTags{
				Description: p.Title,
				Keywords:    "Proyecto, Investigación, UASD",
				OGTitle:     p.Title,
				OGDesc:      p.Status,
			},
			Project: p,
		}

		if err := pages.Project(pData, lang, dict).Render(context.Background(), f); err != nil {
			log.Printf("❌ Error renderizando proyecto %s: %v", p.Slug, err)
		}
		f.Close()
	}

	// 7. Generar News Page
	newsDir := outputDir + "/news"
	if err := os.MkdirAll(newsDir, 0755); err != nil {
		return err
	}

	var newsList []*types.NewsItem
	for _, n := range db.News {
		newsList = append(newsList, n)
	}
	sort.Slice(newsList, func(i, j int) bool {
		return newsList[i].Date.After(newsList[j].Date)
	})

	newsData := pages.NewsData{
		Meta: layouts.MetaTags{
			Description: "Noticias del Instituto de Física de la UASD",
			Keywords:    "Noticias, UASD, Investigación",
		},
		News: newsList,
	}

	fNews, err := os.Create(newsDir + "/index.html")
	if err != nil {
		return err
	}
	defer fNews.Close()

	if err := pages.News(newsData, lang, dict).Render(context.Background(), fNews); err != nil {
		return err
	}

	// 7.1 Generar Páginas Individuales de Noticias
	for i, n := range newsList {
		nDir := outputDir + "/news/" + n.Slug
		if err := os.MkdirAll(nDir, 0755); err != nil {
			continue
		}

		f, err := os.Create(nDir + "/index.html")
		if err != nil {
			continue
		}

		var next, prev *types.NewsItem
		if i > 0 {
			next = newsList[i-1]
		}
		if i < len(newsList)-1 {
			prev = newsList[i+1]
		}

		nData := pages.NewsItemData{
			Meta: layouts.MetaTags{
				Description: n.Summary,
				Keywords:    "Noticia, UASD",
				OGTitle:     n.Title,
				OGDesc:      n.Summary,
				OGImage:     n.Image,
			},
			News: n,
			Next: next,
			Prev: prev,
		}

		if n.Image != "" && n.ImageAlt == "" {
			WarningCount++
			ErrorList = append(ErrorList, fmt.Sprintf("Noticia %s tiene Image pero falta 'image_alt'", n.Title))
		}

		if err := pages.NewsItem(nData, lang, dict).Render(context.Background(), f); err != nil {
			log.Printf("❌ Error renderizando noticia %s: %v", n.Slug, err)
		}
		f.Close()
	}

	// 8. Generar Blog Page
	blogDir := outputDir + "/blog"
	if err := os.MkdirAll(blogDir, 0755); err != nil {
		return err
	}

	var blogList []*types.BlogPost
	for _, p := range db.BlogPosts {
		blogList = append(blogList, p)
	}
	sort.Slice(blogList, func(i, j int) bool {
		return blogList[i].Date.After(blogList[j].Date)
	})

	for _, post := range blogList {
		if author, ok := db.People[post.AuthorID]; ok {
			post.Author = author
		}
	}

	blogData := pages.BlogData{
		Meta: layouts.MetaTags{
			Description: "Blog del Instituto de Física de la UASD",
			Keywords:    "Blog, UASD, Investigación",
		},
		Posts: blogList,
	}

	fBlog, err := os.Create(blogDir + "/index.html")
	if err != nil {
		return err
	}
	defer fBlog.Close()

	if err := pages.Blog(blogData, lang, dict).Render(context.Background(), fBlog); err != nil {
		return err
	}

	// 8.1 Generar Páginas Individuales del Blog
	for i, p := range blogList {
		bDir := outputDir + "/blog/" + p.Slug
		if err := os.MkdirAll(bDir, 0755); err != nil {
			continue
		}

		f, err := os.Create(bDir + "/index.html")
		if err != nil {
			continue
		}

		var next, prev *types.BlogPost
		if i > 0 {
			next = blogList[i-1]
		}
		if i < len(blogList)-1 {
			prev = blogList[i+1]
		}

		bData := pages.BlogPostData{
			Meta: layouts.MetaTags{
				Description: p.Title,
				Keywords:    "Blog, UASD",
				OGTitle:     p.Title,
				OGDesc:      "Blog Post",
			},
			Post: p,
			Next: next,
			Prev: prev,
		}

		if err := pages.BlogPost(bData, lang, dict).Render(context.Background(), f); err != nil {
			log.Printf("❌ Error renderizando post %s: %v", p.Slug, err)
		}
		f.Close()
	}

	// 8.2 Generar Página de Índice de Publicaciones
	publicationsDir := outputDir + "/publications"
	if err := os.MkdirAll(publicationsDir, 0755); err != nil {
		return err
	}

	// Ordenar publicaciones por año (más recientes primero)
	sortedPublications := make([]*types.Publication, len(db.Publications))
	copy(sortedPublications, db.Publications)
	sort.Slice(sortedPublications, func(i, j int) bool {
		return sortedPublications[i].Year > sortedPublications[j].Year
	})

	publicationsData := pages.PublicationsData{
		Meta: layouts.MetaTags{
			Description: dict["PublicationsDescription"],
			Keywords:    "Publicaciones, Papers, Investigación, Física, UASD",
		},
		Publications: sortedPublications,
	}

	fPublications, err := os.Create(publicationsDir + "/index.html")
	if err != nil {
		return err
	}
	defer fPublications.Close()

	if err := pages.Publications(publicationsData, lang, dict).Render(context.Background(), fPublications); err != nil {
		return err
	}

	// 8.3 Generar Páginas Individuales de Publicaciones
	// El directorio ya fue creado en la sección anterior

	for _, pub := range db.Publications {
		pDir := outputDir + "/publications/" + pub.Slug
		if err := os.MkdirAll(pDir, 0755); err != nil {
			continue
		}

		f, err := os.Create(pDir + "/index.html")
		if err != nil {
			continue
		}

		// Usar el abstract como descripción si está disponible
		description := pub.Title
		if pub.Abstract != "" {
			// Limitar el abstract a 200 caracteres para meta description
			if len(pub.Abstract) > 200 {
				description = pub.Abstract[:197] + "..."
			} else {
				description = pub.Abstract
			}
		}

		pData := pages.PublicationData{
			Meta: layouts.MetaTags{
				Description: description,
				Keywords:    strings.Join(pub.Authors, ", ") + ", " + pub.Journal,
				OGTitle:     pub.Title,
				OGDesc:      description,
			},
			Publication: pub,
		}

		if err := pages.Publication(pData, lang, dict).Render(context.Background(), f); err != nil {
			log.Printf("❌ Error renderizando publicación %s: %v", pub.Slug, err)
		}
		f.Close()
	}

	

	// 9. Generar Índice de Búsqueda
	if err := generateSearchIndex(db, outputDir, lang); err != nil {
		log.Printf("⚠️ Error generando índice de búsqueda: %v", err)
	}

	// 10. Generar Página de Búsqueda
	searchDir := outputDir + "/search"
	if err := os.MkdirAll(searchDir, 0755); err != nil {
		return err
	}
	fSearch, err := os.Create(searchDir + "/index.html")
	if err != nil {
		return err
	}
	defer fSearch.Close()

	searchData := pages.SearchData{
		Meta: layouts.MetaTags{
			Description: dict["Search"],
			Keywords:    "Search, Buscador, UASD",
		},
	}
	if err := pages.Search(searchData, lang, dict).Render(context.Background(), fSearch); err != nil {
		return err
	}

	appsIndexDir := outputDir + "/apps"
    if err := os.MkdirAll(appsIndexDir, 0755); err != nil {
        return err
    }
    
    // Ordenar apps (opcional)
    sort.Slice(db.Tools, func(i, j int) bool {
        return db.Tools[i].Title < db.Tools[j].Title
    })

    appsData := pages.AppsData{
        Meta: layouts.MetaTags{
            Description: "Herramientas y Aplicaciones del Instituto de Física",
            Keywords:    "Apps, Herramientas, Física, QR, Markdown",
        },
        Tools: db.Tools,
    }

    fApps, err := os.Create(appsIndexDir + "/index.html")
    if err != nil {
        return err
    }
    defer fApps.Close()

    if err := pages.AppsIndex(appsData, lang, dict).Render(context.Background(), fApps); err != nil {
        return err
    }

	// 11. Generar Página de Apps / QR (NUEVO)
	appsDir := outputDir + "/apps/qr"
	if err := os.MkdirAll(appsDir, 0755); err != nil {
		return err
	}
	fQR, err := os.Create(appsDir + "/index.html")
	if err != nil {
		return err
	}
	defer fQR.Close()

	if err := pages.QRGenerator(lang, dict).Render(context.Background(), fQR); err != nil {
		return err
	}

	// 12. Generar Página de Markdown (App #2)
	mdDir := outputDir + "/apps/markdown"
	if err := os.MkdirAll(mdDir, 0755); err != nil {
		return err
	}
	fMD, err := os.Create(mdDir + "/index.html")
	if err != nil {
		return err
	}
	defer fMD.Close()

	if err := pages.MarkdownEditor(lang, dict).Render(context.Background(), fMD); err != nil {
		return err
	}

	return nil
}

type SearchItem struct {
	Title   string   `json:"title"`
	URL     string   `json:"url"`
	Type    string   `json:"type"`
	Summary string   `json:"summary"`
	Tags    []string `json:"tags"`
}

func generateSearchIndex(db *types.Database, outputDir, lang string) error {
	var items []SearchItem

	// Personas
	for _, p := range db.People {
		items = append(items, SearchItem{
			Title:   p.Name,
			URL:     prefixPath("/people/"+p.Slug, lang),
			Type:    "Person",
			Summary: p.Role,
			Tags:    []string{p.Role, p.Type},
		})
	}

	// Proyectos
	for _, p := range db.Projects {
		items = append(items, SearchItem{
			Title:   p.Title,
			URL:     prefixPath("/projects/"+p.Slug, lang),
			Type:    "Project",
			Summary: p.Status,
			Tags:    append(p.Tags, p.Funding, p.Status),
		})
	}

	// Noticias
	for _, n := range db.News {
		items = append(items, SearchItem{
			Title:   n.Title,
			URL:     prefixPath("/news/"+n.Slug, lang),
			Type:    "News",
			Summary: n.Summary,
			Tags:    []string{"News", "Noticia"},
		})
	}

	// Blog
	for _, b := range db.BlogPosts {
		items = append(items, SearchItem{
			Title:   b.Title,
			URL:     prefixPath("/blog/"+b.Slug, lang),
			Type:    "Blog",
			Summary: "Blog Post",
			Tags:    b.Tags,
		})
	}

	// Escribir JSON
	f, err := os.Create(outputDir + "/search.json")
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	// encoder.SetIndent("", "  ") // Opcional: indentar para debug
	return encoder.Encode(items)
}

func prefixPath(path, lang string) string {
	if lang == "en" {
		return "/en" + path
	}
	return path
}

func checkMissingTranslations(contentDir string) {
	log.Println("🔍 Verificando traducciones faltantes...")
	
	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Ignorar directorio de referencias
			if info.Name() == "references" {
				return filepath.SkipDir
			}
			return nil
		}

		// Solo procesar archivos .md
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		dir := filepath.Dir(path)
		filename := info.Name()

		// Verificar paridad ES <-> EN
		if strings.HasSuffix(filename, ".es.md") {
			baseName := strings.TrimSuffix(filename, ".es.md")
			enFile := filepath.Join(dir, baseName+".en.md")
			if _, err := os.Stat(enFile); os.IsNotExist(err) {
				log.Printf("⚠️  Falta traducción al INGLÉS: %s", path)
			}
		} else if strings.HasSuffix(filename, ".en.md") {
			baseName := strings.TrimSuffix(filename, ".en.md")
			esFile := filepath.Join(dir, baseName+".es.md")
			if _, err := os.Stat(esFile); os.IsNotExist(err) {
				log.Printf("⚠️  Falta traducción al ESPAÑOL: %s", path)
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("Error verificando traducciones: %v", err)
	}
}
