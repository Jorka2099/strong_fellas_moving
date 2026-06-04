package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strong-fellas/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type QuoteRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	From    string `json:"from"`
	To      string `json:"to"`
	Date    string `json:"date"`
	Fellas  string `json:"fellas"`
	Details string `json:"details"`
}

func (q QuoteRequest) ToLead() repository.Lead {
	fellasNum, err := strconv.Atoi(q.Fellas)
	if err != nil {
		fellasNum = 2 // default to 2 if conversion fails
	}
	return repository.Lead{
		Name:         q.Name,
		Phone:        q.Phone,
		MovingFrom:   q.From,
		MovingTo:     q.To,
		MovingDate:   q.Date,
		FellasNumber: fellasNum,
		Details:      q.Details,
	}
}

var DBPool *pgxpool.Pool

func AdminLeadsHandler(w http.ResponseWriter, r *http.Request) {
	leads, err := repository.GetAllLeads(DBPool)
	if err != nil {
		http.Error(w, "Unable to fetch leads", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/admin.html"))
	tmpl.Execute(w, leads)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

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

func SubmitQuoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	req := QuoteRequest{
		Name:    r.FormValue("name"),
		Phone:   r.FormValue("phone"),
		From:    r.FormValue("from"),
		To:      r.FormValue("to"),
		Date:    r.FormValue("date"),
		Fellas:  r.FormValue("fellas"),
		Details: r.FormValue("details"),
	}

	if req.Name == "" || req.Phone == "" || req.From == "" || req.To == "" {
		http.Error(w, "Please fill in all required fields", http.StatusBadRequest)
		return
	}
	reqLeads := req.ToLead()

	err := repository.SaveLead(DBPool, reqLeads)
	if err != nil {
		http.Error(w, "Unable to save your request. Please try again later.", http.StatusInternalServerError)
		return
	}

	fmt.Println("============== NEW REQUEST! ==============")
	fmt.Printf("Client: %s\nTel: %s\nFrom: %s\nWhere: %s\nDetails: %s\n", req.Name, req.Phone, req.From, req.To, req.Details)
	fmt.Println("===========================================")

	w.Write([]byte("Thank you for your request! We will contact you soon."))
}
