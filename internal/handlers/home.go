package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strong-fellas/internal/models"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/home.html",
		"templates/header.html",
		"templates/hero.html",
		"templates/footer.html",
		"templates/quote.html",
	)
	if err != nil {
		http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		http.Error(w, "Rendering Error"+err.Error(), http.StatusInternalServerError)
	}
}

func HandleQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	req := models.QuoteRequest{
		Name:    r.FormValue("name"),
		Phone:   r.FormValue("phone"),
		From:    r.FormValue("from"),
		To:      r.FormValue("to"),
		Details: r.FormValue("details"),
	}

	// Выводим в консоль (позже заменим на запись в админку)
	fmt.Println("============== НОВАЯ ЗАЯВКА! ==============")
	fmt.Printf("Клиент: %s\nТелефон: %s\nОткуда: %s\nКуда: %s\nДетали: %s\n", req.Name, req.Phone, req.From, req.To, req.Details)
	fmt.Println("===========================================")

	w.Write([]byte("Thank you for your request! We will contact you soon."))
}
