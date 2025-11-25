package linker

import (
	"github.com/IfisUASD/ifisuasd.github.io/internal/diagnostics"
	"github.com/IfisUASD/ifisuasd.github.io/internal/types"
)

// LinkData recorre toda la base de datos y resuelve las referencias cruzadas (IDs -> Punteros).
// Modifica el objeto db directamente.
func LinkData(db *types.Database) {
	linkProjectsToPeople(db)
	linkPublicationsToPeople(db)
	linkPublicationsToProjects(db)
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
				diagnostics.LogWarning("⚠️  Warning: Proyecto %s referencia a PI inexistente %s", proj.ID, proj.PrincipalInvestigatorID)
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
				diagnostics.LogWarning("⚠️  Warning: Proyecto %s referencia a CoInvestigador inexistente %s", proj.ID, memberID)
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
				diagnostics.LogWarning("⚠️  Warning: Proyecto %s referencia a Asistente inexistente %s", proj.ID, memberID)
			}
		}
	}
}

func linkPublicationsToPeople(db *types.Database) {
	for _, pub := range db.Publications {
		// A. Rol de AUTOR
		for _, orcid := range pub.AuthorOrcids {
			if person, exists := db.People[orcid]; exists {
				pub.LinkedAuthors = append(pub.LinkedAuthors, person)
				person.Publications = append(person.Publications, pub)
			} else {
				// MEJORA: Advertencia de ORCID huérfano en autores
				diagnostics.LogWarning("⚠️  [Linker] Warning: La publicación '%s' referencia un Autor inexistente (ORCID: %s)", pub.ID, orcid)
			}
		}

		// B. Rol de ASESOR (x-advisors)
		for _, orcid := range pub.AdvisorOrcids {
			if person, exists := db.People[orcid]; exists {
				pub.LinkedAdvisors = append(pub.LinkedAdvisors, person)
				person.Mentored = append(person.Mentored, pub)
			} else {
				// MEJORA: Advertencia de ORCID huérfano en asesores
				diagnostics.LogWarning("⚠️  [Linker] Warning: La publicación '%s' referencia un Asesor inexistente (ORCID: %s)", pub.ID, orcid)
			}
		}
	}
}

func linkPublicationsToProjects(db *types.Database) {
	for _, pub := range db.Publications {
		if pub.ProjectID != "" {
			if project, exists := db.Projects[pub.ProjectID]; exists {
				pub.Project = project
				project.Publications = append(project.Publications, pub)
			} else {
				// MEJORA: Advertencia de Proyecto huérfano
				diagnostics.LogWarning("⚠️  [Linker] Warning: La publicación '%s' referencia un Proyecto inexistente (ID: %s)", pub.ID, pub.ProjectID)
			}
		}
	}
}