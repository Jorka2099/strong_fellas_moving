package models

type QuoteRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	From    string `json:"from"`
	To      string `json:"to"`
	Details string `json:"details"`
}
