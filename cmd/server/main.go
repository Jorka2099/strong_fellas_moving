package main

import (
	"html/template"
	"log"
	"net/http"
	"strong-fellas/internal/handlers"
	"strong-fellas/internal/repository"
)

var tmpl *template.Template

func main() {
	dbPool, err := repository.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}
	defer dbPool.Close()

	tmpl = template.Must(template.ParseGlob("templates/*.html"))

	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/", handlers.HandleHome)
	http.HandleFunc("/quote", handlers.HandleQuote)

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}
