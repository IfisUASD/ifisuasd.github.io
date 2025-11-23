package linker

import (
	"log"

	"github.com/IfisUASD/ifisuasd.github.io/internal/types"
)

// LinkData recorre toda la base de datos y resuelve las referencias cruzadas (IDs -> Punteros).
// Modifica el objeto db directamente.
func LinkData(db *types.Database) {
	linkProjectsToPeople(db)
	linkPapersToPeople(db)
	linkPapersToProjects(db)
}

// linkProjectsToPeople conecta Proyectos con PI, CoInvestigadores y Asistentes.
func linkProjectsToPeople(db *types.Database) {
	for _, proj := range db.Projects {
		// 1. Vincular Investigador Principal (PI)
		if pi, exists := db.People[proj.PrincipalInvestigatorID]; exists {
			// Enlace Proy -> Persona
			proj.PrincipalInvestigator = pi
			// Enlace Persona -> Proy (Bidireccional)
			pi.Projects = append(pi.Projects, proj)
		} else {
			if proj.PrincipalInvestigatorID != "" {
				log.Printf("⚠️  Warning: Proyecto %s referencia a PI inexistente %s", proj.ID, proj.PrincipalInvestigatorID)
			}
		}

		// 2. Vincular Co-Investigadores
		for _, memberID := range proj.CoinvestigatorIDs {
			if member, exists := db.People[memberID]; exists {
				// Enlace Proy -> Persona
				proj.Coinvestigators = append(proj.Coinvestigators, member)
				// Enlace Persona -> Proy (Bidireccional)
				member.Projects = append(member.Projects, proj)
			} else {
				log.Printf("⚠️  Warning: Proyecto %s referencia a CoInvestigador inexistente %s", proj.ID, memberID)
			}
		}

		// 3. Vincular Asistentes de Investigación
		for _, memberID := range proj.ResearchAssistantIDs {
			if member, exists := db.People[memberID]; exists {
				// Enlace Proy -> Persona
				proj.ResearchAssistants = append(proj.ResearchAssistants, member)
				// Enlace Persona -> Proy (Bidireccional)
				member.Projects = append(member.Projects, proj)
			} else {
				log.Printf("⚠️  Warning: Proyecto %s referencia a Asistente inexistente %s", proj.ID, memberID)
			}
		}
	}
}

// linkPapersToPeople conecta Papers con sus Autores (usando x-orcids).
func linkPapersToPeople(db *types.Database) {
	for _, paper := range db.Papers {
		for _, orcid := range paper.AuthorOrcids {
			if person, exists := db.People[orcid]; exists {
				// Enlace Persona -> Paper
				person.Publications = append(person.Publications, paper)
			} else {
				log.Printf("⚠️  Warning: Paper %s referencia a Autor inexistente %s", paper.ID, orcid)
			}
		}
	}
}

// linkPapersToProjects conecta Papers con Proyectos (usando x-project).
func linkPapersToProjects(db *types.Database) {
	for _, paper := range db.Papers {
		if paper.ProjectID != "" {
			if project, exists := db.Projects[paper.ProjectID]; exists {
				// Enlace Proyecto -> Paper
				project.Publications = append(project.Publications, paper)
			} else {
				log.Printf("⚠️  Warning: Paper %s referencia a Proyecto inexistente %s", paper.ID, paper.ProjectID)
			}
		}
	}
}