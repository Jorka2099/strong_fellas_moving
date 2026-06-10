package main

import (
	"html/template"
	"log"
	"net/http"
	"strong-fellas/internal/handlers"
	"strong-fellas/internal/repository"

	"github.com/joho/godotenv"
)

var tmpl *template.Template

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Caution, .env file not found, relying on environment variables")
	}

	dbPool, err := repository.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}
	defer dbPool.Close()

	handlers.DBPool = dbPool

	tmpl = template.Must(template.ParseGlob("templates/*.html"))

	handlers.InitTemplate(tmpl)

	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/submit-quote", handlers.RateLimitMiddleware(handlers.SubmitQuoteHandler))

	http.HandleFunc("/privacy.html", func(w http.ResponseWriter, r *http.Request) {
		// Рендерим напрямую сам файл privacy.html
		err := tmpl.ExecuteTemplate(w, "privacy.html", nil)
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/terms.html", func(w http.ResponseWriter, r *http.Request) {
		// Рендерим напрямую сам файл terms.html
		err := tmpl.ExecuteTemplate(w, "terms.html", nil)
		if err != nil {
			http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/admin/login", handlers.RateLimitMiddleware(handlers.AdminLoginHandler))
	http.HandleFunc("/admin/leads", handlers.AuthMiddleware(handlers.AdminLeadsHandler))
	http.HandleFunc("/admin/leads/delete", handlers.AuthMiddleware(handlers.AdminDeleteLeadHandler))
	http.HandleFunc("/admin/logout", handlers.AdminLogoutHandler)
	http.HandleFunc("/admin/leads/export", handlers.AuthMiddleware(handlers.ExportLeadsHandler(dbPool)))

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}

}
