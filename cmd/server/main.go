package main

import (
	"log"
	"net/http"
)

func main() {
	// Servir el directorio ./output
	fs := http.FileServer(http.Dir("./output"))
	http.Handle("/", fs)

	port := ":8180"
	log.Printf("🚀 Servidor iniciado en http://localhost%s", port)
	log.Println("Presiona Ctrl+C para detener")
	
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
