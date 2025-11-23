package loader

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/IfisUASD/ifisuasd.github.io/internal/parsers"
	"github.com/IfisUASD/ifisuasd.github.io/internal/types"
)

// LoadContent recorre recursivamente el directorio rootPath y carga todos los contenidos
// en una nueva instancia de Database, filtrando por el idioma especificado (ej: "es", "en").
func LoadContent(rootPath string, lang string) (*types.Database, error) {
	db := types.NewDatabase()

	// Mapa temporal para agrupar archivos por ID base y tipo
	// Key: Tipo/ID -> Value: Map[Lang]FilePath
	type FileEntry struct {
		Path    string
		Content []byte
	}
	filesMap := make(map[string]map[string]FileEntry)

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		// Leer contenido
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", path, err)
		}

		relPath, _ := filepath.Rel(rootPath, path)
		relPath = filepath.ToSlash(relPath)

		// Identificar idioma del archivo
		fileLang := "default"
		baseName := filepath.Base(path)
		ext := filepath.Ext(path)
		nameWithoutExt := strings.TrimSuffix(baseName, ext)

		if strings.HasSuffix(nameWithoutExt, ".es") {
			fileLang = "es"
			nameWithoutExt = strings.TrimSuffix(nameWithoutExt, ".es")
		} else if strings.HasSuffix(nameWithoutExt, ".en") {
			fileLang = "en"
			nameWithoutExt = strings.TrimSuffix(nameWithoutExt, ".en")
		}

		// Clave única para agrupar: Directorio + NombreBase (sin idioma)
		// Ej: people/vladimir-perez
		dir := filepath.Dir(relPath)
		key := filepath.Join(dir, nameWithoutExt)

		if _, ok := filesMap[key]; !ok {
			filesMap[key] = make(map[string]FileEntry)
		}
		filesMap[key][fileLang] = FileEntry{Path: path, Content: content}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Procesar archivos seleccionados
	for _, variants := range filesMap {
		// Selección de mejor candidato
		var selected FileEntry
		found := false

		// 1. Intentar idioma exacto
		if entry, ok := variants[lang]; ok {
			selected = entry
			found = true
		} else if entry, ok := variants["default"]; ok {
			// 2. Fallback a default (sin sufijo)
			selected = entry
			found = true
		} else if lang == "es" && len(variants) > 0 {
             // 3. Fallback extra: si pedimos español y no hay exacto ni default,
             // pero hay variants (ej: solo .en), podríamos decidir qué hacer.
             // Por ahora, estricto: si no hay match ni default, se ignora o se toma cualquiera?
             // Regla de negocio: Si no existe en el idioma pedido ni default, ¿mostramos otro?
             // Asumamos que "default" es el fallback universal. Si solo existe .en y pedimos .es,
             // y no hay default, quizás deberíamos mostrar .en?
             // Implementación simple: Preferir Lang > Default > Primer disponible
             for _, v := range variants {
                 selected = v
                 found = true
                 break
             }
        } else {
             // Tomar cualquiera como último recurso
             for _, v := range variants {
                 selected = v
                 found = true
                 break
             }
        }

		if !found {
			continue
		}

		// Parsear según tipo
		path := selected.Path
		content := selected.Content
		
		// Normalizar path para switch
		relPath, _ := filepath.Rel(rootPath, path)
		relPath = filepath.ToSlash(relPath)

		switch {
		case strings.Contains(relPath, "people/") && strings.HasSuffix(path, ".md"):
			person, err := parsers.ParsePerson(path, content)
			if err != nil {
				return nil, fmt.Errorf("error parsing person %s: %w", path, err)
			}
			db.People[person.ID] = person

		case strings.Contains(relPath, "projects/") && strings.HasSuffix(path, ".md"):
			project, err := parsers.ParseProject(path, content)
			if err != nil {
				return nil, fmt.Errorf("error parsing project %s: %w", path, err)
			}
			db.Projects[project.ID] = project

		case strings.Contains(relPath, "references/") && strings.HasSuffix(path, ".bib"):
			// BibTeX no suele tener variantes por idioma, pero si las tuviera, funcionaría igual
			papers, err := parsers.ParseBibTeX(path, content)
			if err != nil {
				return nil, fmt.Errorf("error parsing bibtex %s: %w", path, err)
			}
			db.Papers = append(db.Papers, papers...)

		case strings.Contains(relPath, "news/") && strings.HasSuffix(path, ".md"):
			news, err := parsers.ParseNewsItem(path, content)
			if err != nil {
				return nil, fmt.Errorf("error parsing news %s: %w", path, err)
			}
			db.News = append(db.News, news)

		case strings.Contains(relPath, "blog/") && strings.HasSuffix(path, ".md"):
			post, err := parsers.ParseBlogPost(path, content)
			if err != nil {
				return nil, fmt.Errorf("error parsing blog post %s: %w", path, err)
			}
			db.BlogPosts = append(db.BlogPosts, post)
		}
	}

	return db, nil
}
