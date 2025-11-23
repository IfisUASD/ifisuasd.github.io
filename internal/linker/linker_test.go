package linker

import (
	"testing"

	"github.com/IfisUASD/ifisuasd.github.io/internal/types"
)

func TestLinkData_Projects(t *testing.T) {
	// 1. Setup: Crear una DB falsa en memoria
	db := types.NewDatabase()

	// Crear Personas
	piID := "0000-PI"
	coInvID := "0000-COINV"
	
	pi := &types.Person{ID: piID, Name: "Dr. Principal"}
	coInv := &types.Person{ID: coInvID, Name: "Dr. CoInvestigador"}
	
	db.People[piID] = pi
	db.People[coInvID] = coInv

	// Crear un Proyecto que referencia a esas Personas
	projID := "PROJ-01"
	proj := &types.Project{
		ID:                      projID,
		Title:                   "Proyecto Test",
		PrincipalInvestigatorID: piID,
		CoinvestigatorIDs:        []string{coInvID}, // Usamos el nuevo campo
	}
	db.Projects[projID] = proj

	// 2. Ejecución: Correr el vinculador
	LinkData(db)

	// 3. Aserciones

	if len(pi.Projects) == 0 {
		t.Error("El Linker no añadió el proyecto al historial del PI")
	}

	if pi.Projects[0].Title != proj.Title {
		t.Error("El Linker no añadió correctamente el proyecto al historial del PI")
	}

	if len(coInv.Projects) == 0 {
		t.Error("El Linker no añadió el proyecto al historial del CoInvestigador")
	}

	if coInv.Projects[0].Title != proj.Title {
		t.Error("El Linker no añadió correctamente el proyecto al historial del PI")
	}
}

func TestLinkData_Papers(t *testing.T) {
	db := types.NewDatabase()

	// 1. Setup
	authorID := "0000-AUTHOR"
	author := &types.Person{ID: authorID, Name: "Dr. Author"}
	db.People[authorID] = author

	projID := "PROJ-PAPER"
	proj := &types.Project{ID: projID, Title: "Project Paper"}
	db.Projects[projID] = proj

	paper := &types.Paper{
		ID:           "DOI-1",
		Title:        "My Paper",
		AuthorOrcids: []string{authorID},
		ProjectID:    projID,
	}
	db.Papers = append(db.Papers, paper)

	// 2. Ejecución
	LinkData(db)

	// 3. Aserciones
	
	// Verificar enlace Paper -> Persona
	if len(author.Publications) != 1 {
		t.Errorf("El Linker no añadió el paper al autor. Tiene %d", len(author.Publications))
	}
	if author.Publications[0].Title != "My Paper" {
		t.Error("Paper incorrecto en autor")
	}

	// Verificar enlace Paper -> Proyecto
	if len(proj.Publications) != 1 {
		t.Errorf("El Linker no añadió el paper al proyecto. Tiene %d", len(proj.Publications))
	}
	if proj.Publications[0].Title != "My Paper" {
		t.Error("Paper incorrecto en proyecto")
	}
}