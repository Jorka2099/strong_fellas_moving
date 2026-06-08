package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/http"
	"sync"
	"time"
	"os"
)

var (
	sessions      = make(map[string]time.Time)
	sessionsMutex sync.Mutex
	tmpl          *template.Template
)

const sessionCookieName = "session_token"

func InitTemplate(t *template.Template) {
	tmpl = t
}

func getAdminPassword() string {
	pass := os.Getenv("ADMIN_PASSWORD")
	if pass == "" {
		return ""
	}
	return pass
}	

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		password := r.FormValue("password")
		actualPassword := getAdminPassword()

		if password != actualPassword || actualPassword == "" {
			tmpl.ExecuteTemplate(w, "login.html", "Invalid password")
			return
		}

		b := make([]byte, 32)
		rand.Read(b)
		sessionToken := base64.URLEncoding.EncodeToString(b)

		expiresAt := time.Now().Add(12 * time.Hour)

		sessionsMutex.Lock()
		sessions[sessionToken] = expiresAt
		sessionsMutex.Unlock()

		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    sessionToken,
			Expires:  expiresAt,
			HttpOnly: true,
			Path:     "/",
		})

		http.Redirect(w, r, "/admin/leads", http.StatusSeeOther)

	}
}

func isValidSession(token string) bool {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	expiresAt, exists := sessions[token]
	if !exists {
		return false
	}

	if time.Now().After(expiresAt) {
		delete(sessions, token)
		return false
	}

	return true
}

func deleteSession(token string) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	delete(sessions, token)
}

func AdminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		deleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:    "/",
		MaxAge:  -1,
	})
	
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}