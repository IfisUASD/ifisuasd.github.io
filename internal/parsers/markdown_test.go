package parsers

import (
	"testing"
	"time"
)

// TestParsePerson verifica que un archivo .md de persona se convierta correctamente en el struct Person.
func TestParsePerson(t *testing.T) {
	// 1. Datos de entrada simulados (Mock)
	filename := "vladimir-perez.md"
	rawContent := []byte(`---
orcid: "0000-0002-1825-0097"
name: "Vladimir Pérez"
role: "Director"
type: "academic"
social:
  scholar: "https://scholar.google.com/..."
---
# Biografía
Hola, soy físico.
`)

	// 2. Ejecución
	person, err := ParsePerson(filename, rawContent)

	// 3. Aserciones (Validar resultados)
	if err != nil {
		t.Fatalf("ParsePerson devolvió error inesperado: %v", err)
	}
	if person == nil {
		t.Fatal("ParsePerson devolvió nil")
	}

	// Validar campos mapeados del YAML
	if person.ID != "0000-0002-1825-0097" {
		t.Errorf("Esperaba ORCID '0000-0002-1825-0097', obtuvo '%s'", person.ID)
	}
	if person.Name != "Vladimir Pérez" {
		t.Errorf("Esperaba Name 'Vladimir Pérez', obtuvo '%s'", person.Name)
	}
	if person.Slug != "vladimir-perez" {
		t.Errorf("Esperaba Slug 'vladimir-perez', obtuvo '%s'", person.Slug)
	}
	
	// Validar mapa anidado
	if person.Social["scholar"] == "" {
		t.Error("No se parseó correctamente el mapa 'social'")
	}

	// Validar conversión de Markdown (básica)
	if person.BioHTML == "" {
		t.Error("El BioHTML está vacío, se esperaba contenido parseado")
	}
}

// TestParseProject verifica la lectura de proyectos con los nuevos campos.
func TestParseProject(t *testing.T) {
	filename := "red-aire.md"
	rawContent := []byte(`---
project_id: "FONDOCYT-2024"
title: "Red de Aire"
start_date: "2024-01-01"
principal_investigator: "0000-0002-1825-0097"
coinvestigator:
  - "0000-0001-COINV-1111"
research_assistant:
  - "0000-0003-ASIST-2222"
---
Descripción del proyecto.`)

	proj, err := ParseProject(filename, rawContent)

	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	if proj.ID != "FONDOCYT-2024" {
		t.Errorf("ID incorrecto: %s", proj.ID)
	}

	// Validar fechas
	expectedDate, _ := time.Parse("2006-01-02", "2024-01-01")
	if !proj.StartDate.Equal(expectedDate) {
		t.Errorf("Fecha de inicio incorrecta. Esperaba %v, obtuvo %v", expectedDate, proj.StartDate)
	}

	// Validar CoInvestigadores
	if len(proj.CoinvestigatorIDs) != 1 {
		t.Errorf("Esperaba 1 coinvestigador, obtuvo %d", len(proj.CoinvestigatorIDs))
	}
	if proj.CoinvestigatorIDs[0] != "0000-0001-COINV-1111" {
		t.Errorf("ID de coinvestigador incorrecto: %s", proj.CoinvestigatorIDs[0])
	}

	// Validar Asistentes
	if len(proj.ResearchAssistantIDs) != 1 {
		t.Errorf("Esperaba 1 asistente, obtuvo %d", len(proj.ResearchAssistantIDs))
	}
	if proj.ResearchAssistantIDs[0] != "0000-0003-ASIST-2222" {
		t.Errorf("ID de asistente incorrecto: %s", proj.ResearchAssistantIDs[0])
	}
}

func TestParseNewsItem(t *testing.T) {
	filename := "evento-fisica.md"
	rawContent := []byte(`---
title: "Evento de Física"
date: "2024-03-15"
summary: "Un gran evento."
---
Detalles del evento.`)

	news, err := ParseNewsItem(filename, rawContent)
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	if news.Title != "Evento de Física" {
		t.Errorf("Título incorrecto: %s", news.Title)
	}
	expectedDate, _ := time.Parse("2006-01-02", "2024-03-15")
	if !news.Date.Equal(expectedDate) {
		t.Errorf("Fecha incorrecta")
	}
	if news.ID != "evento-fisica" {
		t.Errorf("ID debería ser el slug por defecto: %s", news.ID)
	}
}

func TestParseBlogPost(t *testing.T) {
	filename := "mi-post.md"
	rawContent := []byte(`---
title: "Mi Post"
date: "2024-04-20"
author_id: "0000-0002-1825-0097"
tags: ["fisica", "educacion"]
---
Contenido del post.`)

	post, err := ParseBlogPost(filename, rawContent)
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}

	if post.Title != "Mi Post" {
		t.Errorf("Título incorrecto: %s", post.Title)
	}
	if post.AuthorID != "0000-0002-1825-0097" {
		t.Errorf("AuthorID incorrecto: %s", post.AuthorID)
	}
	if len(post.Tags) != 2 {
		t.Errorf("Esperaba 2 tags")
	}
}