package parsers

import (
	"testing"
)

func TestParseBibTeX(t *testing.T) {
	// 1. Datos de entrada simulados
	filename := "institute.bib"
	rawContent := []byte(`
@article{key1,
  title = {Quantum Mechanics},
  author = {Pérez, V. and Montero, E.},
  journal = {Physical Review B},
  year = {2024},
  doi = {10.1103/PhysRevB.100.010101},
  x-orcids = {0000-0002-1825-0097, 0000-0001-COINV-1111},
  x-project = {FONDOCYT-2024}
}
`)

	// 2. Ejecución
	papers, err := ParseBibTeX(filename, rawContent)

	// 3. Aserciones
	if err != nil {
		t.Fatalf("ParseBibTeX devolvió error inesperado: %v", err)
	}
	if len(papers) != 1 {
		t.Fatalf("Esperaba 1 paper, obtuvo %d", len(papers))
	}

	p := papers[0]
	if p.Title != "Quantum Mechanics" {
		t.Errorf("Título incorrecto: %s", p.Title)
	}
	if p.DOI != "10.1103/PhysRevB.100.010101" {
		t.Errorf("DOI incorrecto: %s", p.DOI)
	}
	if p.Year != 2024 {
		t.Errorf("Año incorrecto: %d", p.Year)
	}

	// Validar campos custom
	if len(p.AuthorOrcids) != 2 {
		t.Errorf("Esperaba 2 ORCIDs, obtuvo %d", len(p.AuthorOrcids))
	}
	if p.AuthorOrcids[0] != "0000-0002-1825-0097" {
		t.Errorf("Primer ORCID incorrecto: %s", p.AuthorOrcids[0])
	}
	if p.ProjectID != "FONDOCYT-2024" {
		t.Errorf("ProjectID incorrecto: %s", p.ProjectID)
	}
}
