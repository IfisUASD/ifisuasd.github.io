package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sort"

	"github.com/IfisUASD/ifisuasd.github.io/internal/i18n"
	"github.com/IfisUASD/ifisuasd.github.io/internal/linker"
	"github.com/IfisUASD/ifisuasd.github.io/internal/loader"
	"github.com/IfisUASD/ifisuasd.github.io/internal/types"
	"github.com/IfisUASD/ifisuasd.github.io/templates/layouts"
	"github.com/IfisUASD/ifisuasd.github.io/templates/pages"
)

func main() {
	log.Println("🏗️  Iniciando generador del sitio...")

	// Generar sitio en Español (Default)
	if err := buildSite("es", "./output"); err != nil {
		log.Fatalf("❌ Error generando sitio en Español: %v", err)
	}

	// Generar sitio en Inglés
	if err := buildSite("en", "./output/en"); err != nil {
		log.Fatalf("❌ Error generando sitio en Inglés: %v", err)
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
	log.Printf("✅ Contenido cargado: %d Personas, %d Proyectos, %d Papers", len(db.People), len(db.Projects), len(db.Papers))

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

	recentPapers := db.Papers
	sort.Slice(recentPapers, func(i, j int) bool {
		return recentPapers[i].Year > recentPapers[j].Year
	})
	if len(recentPapers) > 5 {
		recentPapers = recentPapers[:5]
	}

	homeData := pages.HomeData{
		Meta: layouts.MetaTags{
			Description: "Sitio oficial del Instituto de Física de la UASD",
			Keywords:    "Física, UASD, Investigación, Ciencia",
		},
		LatestNews:   latestNews,
		RecentPapers: recentPapers,
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

	return nil
}

type SearchItem struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
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
		})
	}

	// Proyectos
	for _, p := range db.Projects {
		items = append(items, SearchItem{
			Title:   p.Title,
			URL:     prefixPath("/projects/"+p.Slug, lang),
			Type:    "Project",
			Summary: p.Status,
		})
	}

	// Noticias
	for _, n := range db.News {
		items = append(items, SearchItem{
			Title:   n.Title,
			URL:     prefixPath("/news/"+n.Slug, lang),
			Type:    "News",
			Summary: n.Summary,
		})
	}

	// Blog
	for _, b := range db.BlogPosts {
		items = append(items, SearchItem{
			Title:   b.Title,
			URL:     prefixPath("/blog/"+b.Slug, lang),
			Type:    "Blog",
			Summary: "Blog Post",
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
