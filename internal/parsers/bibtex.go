package parsers

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/IfisUASD/ifisuasd.github.io/internal/types"
	"github.com/nickng/bibtex"
)

// ParseBibTeX procesa el contenido de un archivo .bib y devuelve una lista de Papers.
func ParseBibTeX(filename string, content []byte) ([]*types.Paper, error) {
	bib, err := bibtex.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("error parsing bibtex: %w", err)
	}

	var papers []*types.Paper

	for _, entry := range bib.Entries {
		p := &types.Paper{
			ID:    entry.CiteName,
			Type:  entry.Type,
			Title: getField(entry, "title"),
			DOI:   getField(entry, "doi"),
			URL:   getField(entry, "url"),
		}

		// Parse Year
		if yearStr := getField(entry, "year"); yearStr != "" {
			if y, err := strconv.Atoi(yearStr); err == nil {
				p.Year = y
			}
		}

		// Parse Authors
		if authors := getField(entry, "author"); authors != "" {
			p.Authors = []string{authors} 
		}

		// Parse Custom Fields
		if orcids := getField(entry, "x-orcids"); orcids != "" {
			orcids = strings.Trim(orcids, "{}")
			parts := strings.Split(orcids, ",")
			for _, part := range parts {
				p.AuthorOrcids = append(p.AuthorOrcids, strings.TrimSpace(part))
			}
		}

		if project := getField(entry, "x-project"); project != "" {
			p.ProjectID = strings.Trim(strings.Trim(project, "{}"), " ")
		}

		papers = append(papers, p)
	}

	return papers, nil
}

func getField(entry *bibtex.BibEntry, key string) string {
	if val, ok := entry.Fields[key]; ok {
		return val.String()
	}
	return ""
}
