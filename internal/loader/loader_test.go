package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadContent(t *testing.T) {
	// 1. Crear estructura de directorios temporal
	tmpDir := t.TempDir()
	
	dirs := []string{"people", "projects", "references", "news", "blog"}
	for _, d := range dirs {
		if err := os.Mkdir(filepath.Join(tmpDir, d), 0755); err != nil {
			t.Fatal(err)
		}
	}

	// 2. Crear archivos de prueba
	createFile(t, tmpDir, "people/person.md", `---
orcid: "0000-0000"
name: "Test Person"
---
Bio`)
	createFile(t, tmpDir, "projects/project.md", `---
project_id: "PROJ-1"
title: "Test Project"
---
Desc`)
	createFile(t, tmpDir, "references/refs.bib", `@article{key1, title={Paper 1}}`)
	createFile(t, tmpDir, "news/news.md", `---
title: "News 1"
date: "2024-01-01"
---
News content`)
	createFile(t, tmpDir, "blog/post.md", `---
title: "Blog Post 1"
date: "2024-01-01"
---
Blog content`)

	// 3. Ejecutar Loader (Default lang "es")
	db, err := LoadContent(tmpDir, "es")
	if err != nil {
		t.Fatalf("LoadContent falló: %v", err)
	}

	// 4. Verificar resultados
	if len(db.People) != 1 {
		t.Errorf("Esperaba 1 persona, obtuvo %d", len(db.People))
	}
	if _, ok := db.People["0000-0000"]; !ok {
		t.Error("Persona no encontrada por ID")
	}

	if len(db.Projects) != 1 {
		t.Errorf("Esperaba 1 proyecto, obtuvo %d", len(db.Projects))
	}

	if len(db.Papers) != 1 {
		t.Errorf("Esperaba 1 paper, obtuvo %d", len(db.Papers))
	}

	if len(db.News) != 1 {
		t.Errorf("Esperaba 1 noticia, obtuvo %d", len(db.News))
	}

	if len(db.BlogPosts) != 1 {
		t.Errorf("Esperaba 1 post, obtuvo %d", len(db.BlogPosts))
	}
}

func TestLoadContent_I18n(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Crear estructura
	os.Mkdir(filepath.Join(tmpDir, "people"), 0755)

	// Crear archivo base (español implícito o default)
	createFile(t, tmpDir, "people/bio.md", `---
orcid: "0000-BIO"
name: "Bio Default"
---
Bio Default`)

	// Crear archivo inglés
	createFile(t, tmpDir, "people/bio.en.md", `---
orcid: "0000-BIO"
name: "Bio English"
---
Bio English`)

	// Caso 1: Pedir "en" -> Debe cargar bio.en.md
	dbEn, err := LoadContent(tmpDir, "en")
	if err != nil {
		t.Fatal(err)
	}
	if p, ok := dbEn.People["0000-BIO"]; !ok || p.Name != "Bio English" {
		t.Errorf("Esperaba 'Bio English', obtuvo '%s'", p.Name)
	}

	// Caso 2: Pedir "es" -> Debe cargar bio.md (default)
	dbEs, err := LoadContent(tmpDir, "es")
	if err != nil {
		t.Fatal(err)
	}
	if p, ok := dbEs.People["0000-BIO"]; !ok || p.Name != "Bio Default" {
		t.Errorf("Esperaba 'Bio Default', obtuvo '%s'", p.Name)
	}
}

func createFile(t *testing.T, root, path, content string) {
	fullPath := filepath.Join(root, path)
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
