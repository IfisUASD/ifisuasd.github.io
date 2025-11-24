package parsers

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/IfisUASD/ifisuasd.github.io/internal/types"
	"github.com/nickng/bibtex"
)
func ParseBibTeX(filename string, content []byte) ([]*types.Publication, error) {
	bib, err := bibtex.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("error parsing bibtex: %w", err)
	}

	var pubs []*types.Publication

	for _, entry := range bib.Entries {
		p := &types.Publication{
			ID:        entry.CiteName,
			Slug:      entry.CiteName,
			Type:      entry.Type,
			Title:     DecodeLaTeX(getField(entry, "title")),
			DOI:       getField(entry, "doi"),
			URL:       getField(entry, "url"),
			Journal:   DecodeLaTeX(getField(entry, "journal")),
			Volume:    getField(entry, "volume"),
			Number:    getField(entry, "number"),
			Pages:     getField(entry, "pages"),
			Publisher: DecodeLaTeX(getField(entry, "publisher")),
			School:    DecodeLaTeX(getField(entry, "school")),
			Booktitle: DecodeLaTeX(getField(entry, "booktitle")),
			Abstract:  DecodeLaTeX(getField(entry, "abstract")),
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

		pubs = append(pubs, p)
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
