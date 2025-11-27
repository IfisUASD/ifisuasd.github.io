package parsers

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/IfisUASD/ifisuasd.github.io/internal/types"
	"github.com/nickng/bibtex"
)
func ParseBibTeX(filename string, content []byte) ([]*types.Publication, error) {
	bib, err := bibtex.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("❌ ERROR FATAL DE SINTAXIS en archivo '%s': %w (Revisa llaves de cierre o comas faltantes)", filename, err)
	}

	var pubs []*types.Publication

	for _, entry := range bib.Entries {
		// Variable para contexto de error
		currentField := ""
		
		// Función auxiliar para capturar pánicos o errores en campos específicos
		safeDecode := func(fieldKey string) string {
			currentField = fieldKey
			raw := getField(entry, fieldKey)
			// Aquí podrías validar paréntesis balanceados si fuera necesario antes de DecodeLaTeX
			return DecodeLaTeX(raw)
		}

		// Usamos un bloque anónimo para recuperar pánicos si DecodeLaTeX es inestable
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("❌ [BibTeX Error] Pánico procesando entrada '%s' en el campo '%s': %v", entry.CiteName, currentField, r)
				}
			}()

			p := &types.Publication{
				ID:        entry.CiteName,
				Slug:      entry.CiteName,
				Type:      entry.Type,
				Title:     safeDecode("title"),
				DOI:       getField(entry, "doi"),
				URL:       getField(entry, "url"),
				Journal:   safeDecode("journal"),
				Volume:    getField(entry, "volume"),
				Number:    getField(entry, "number"),
				Pages:     getField(entry, "pages"),
				Publisher: safeDecode("publisher"),
				School:    safeDecode("school"),
				Booktitle: safeDecode("booktitle"),
				Abstract:  safeDecode("abstract"),
			}

			if yearStr := getField(entry, "year"); yearStr != "" {
				if y, err := strconv.Atoi(yearStr); err == nil {
					p.Year = y
				}
			}

			if authors := getField(entry, "author"); authors != "" {
				// Limpieza básica de " and " que usa BibTeX
				rawAuthors := strings.Split(authors, " and ")
				for _, auth := range rawAuthors {
					p.Authors = append(p.Authors, DecodeLaTeX(auth))
				}
			}

			// --- CAMPOS PERSONALIZADOS ---

			// 1. Autores (x-orcids)
			if val := getField(entry, "x-orcids"); val != "" {
				p.AuthorOrcids = parseList(val)
			}

			// 2. Asesores (x-advisors) - NUEVO
			if val := getField(entry, "x-advisors"); val != "" {
				p.AdvisorOrcids = parseList(val)
			}

			// 3. Proyecto (x-project)
			if val := getField(entry, "x-project"); val != "" {
				p.ProjectID = strings.Trim(strings.Trim(val, "{}"), " ")
			}

			if p.Title == "" {
				log.Printf("⚠️  [BibTeX Warning] La entrada '%s' tiene un título vacío o inválido.", entry.CiteName)
			}

			// Validar campos desconocidos para evitar errores de tipeo (ej: x-orcid vs x-orcids)
			allowedFields := map[string]bool{
				// Standard BibTeX
				"title": true, "author": true, "year": true, "journal": true,
				"volume": true, "number": true, "pages": true, "publisher": true,
				"school": true, "booktitle": true, "abstract": true, "doi": true,
				"url": true, "month": true, "editor": true, "series": true,
				"address": true, "edition": true, "howpublished": true,
				"institution": true, "note": true, "key": true, "crossref": true,
				"type": true, "isbn": true, "issn": true, "copyright": true,
				"language": true, "location": true, "keywords": true,
				
				// Custom Fields
				"x-orcids": true, "x-advisors": true, "x-project": true,
				"x-fetchedfrom": true, // Metadata tool
			}

			for key := range entry.Fields {
				keyLower := strings.ToLower(key)
				if !allowedFields[keyLower] {
					log.Printf("⚠️  [BibTeX Warning] Campo desconocido '%s' en entrada '%s'. ¿Es un error de tipeo?", key, entry.CiteName)
				}
			}

			pubs = append(pubs, p)
		}()
	}
	return pubs, nil
}

func parseList(raw string) []string {
	raw = strings.Trim(raw, "{}")
	parts := strings.Split(raw, ",")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func getField(entry *bibtex.BibEntry, key string) string {
	if val, ok := entry.Fields[key]; ok {
		return val.String()
	}
	return ""
}
