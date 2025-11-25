package diagnostics

import (
	"fmt"
	"log"
)

var WarningCount int
var ErrorList []string

// LogWarning imprime el mensaje inmediatamente y lo guarda para el resumen final
func LogWarning(msg string, args ...interface{}) {
	formatted := fmt.Sprintf(msg, args...)
	// Imprimimos con prefijo amarillo (si la terminal lo soporta) o texto plano
	log.Printf("⚠️  %s", formatted) 
	ErrorList = append(ErrorList, formatted)
	WarningCount++
}