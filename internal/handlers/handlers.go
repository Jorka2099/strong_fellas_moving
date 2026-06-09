package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strong-fellas/internal/repository"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type QuoteRequest struct {
	Name       string
	Phone      string
	From       string
	To         string
	Date       string
	Fellas     string
	Hours      string
	TotalPrice string
	Details    string
}

func (q QuoteRequest) ToLead() repository.Lead {
	fellasNum, err := strconv.Atoi(q.Fellas)
	if err != nil {
		fellasNum = 2
	}
	hoursNum, _ := strconv.Atoi(q.Hours)
	priceNum, _ := strconv.Atoi(q.TotalPrice)

	return repository.Lead{
		Name:         q.Name,
		Phone:        q.Phone,
		MovingFrom:   q.From,
		MovingTo:     q.To,
		MovingDate:   q.Date,
		FellasNumber: fellasNum,
		Hours:        hoursNum,
		TotalPrice:   priceNum,
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

// AdminDeleteLeadHandler deletes lead by ID
func AdminDeleteLeadHandler(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Lead ID", http.StatusBadRequest)
		return
	}

	
	err = repository.DeleteLead(DBPool, id) 
	if err != nil {
		log.Printf("Error deleting lead: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// После успешного удаления перенаправляем админа обратно на список лидов
	http.Redirect(w, r, "/admin/leads", http.StatusSeeOther)
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
		"templates/about.html",
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
		Name:       r.FormValue("name"),
		Phone:      r.FormValue("phone"),
		From:       r.FormValue("from"),
		To:         r.FormValue("to"),
		Date:       r.FormValue("date"),
		Fellas:     r.FormValue("fellas"),
		Hours:      r.FormValue("hours"),
		TotalPrice: r.FormValue("total_price"),
		Details:    r.FormValue("details"),
	}

	if req.From == "" && req.To == "" {
		http.Error(w, "Please provide at least a 'Moving From' or 'Moving To' address", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Phone == "" || req.Date == "" {
		http.Error(w, "Please fill in all required fields", http.StatusBadRequest)
		return
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	today := time.Now().Truncate(24 * time.Hour)
	if parsedDate.Before(today) {
		http.Error(w, "Moving date cannot be in the past", http.StatusBadRequest)
		return
	}

	reqLeads := req.ToLead()

	err = repository.SaveLead(DBPool, reqLeads)
	if err != nil {
		log.Printf("DB error in SubmitQuoteHandler: %v\n", err)

		http.Error(w, "Unable to save your request. Please try again later", http.StatusInternalServerError)
		return
	}

	fmt.Println("============== NEW REQUEST! ==============")
	fmt.Printf("Client: %s\nTel: %s\nFrom: %s\nWhere: %s\nDate: %s\nFellas: %d\nHours: %d\nTotal: $%d\nDetails: %s\n",
		reqLeads.Name, reqLeads.Phone, reqLeads.MovingFrom, reqLeads.MovingTo, reqLeads.MovingDate, reqLeads.FellasNumber, reqLeads.Hours, reqLeads.TotalPrice, reqLeads.Details)
	fmt.Println("===========================================")

	w.Write([]byte("Thank you for your request! We will contact you soon."))
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil || !isValidSession(cookie.Value) {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
