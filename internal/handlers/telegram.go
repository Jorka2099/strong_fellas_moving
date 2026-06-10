package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"strong-fellas/internal/repository"
)

// sendTelegramNotification sends request to Telegram.
func sendTelegramNotification(lead repository.Lead) {

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatID == "" {
		log.Println("[Telegram] bot token or chat_id was not set. Skipping Telegram notification.")
		return
	}

	from := lead.MovingFrom
	if from == "" {
		from = "—"
	}
	to := lead.MovingTo
	if to == "" {
		to = "—"
	}
	details := lead.Details
	if details == "" {
		details = "—"
	}

	text := fmt.Sprintf(
		"🔔 *NEW LEAD RECEIVED!*\n\n"+
			"👤 *Name:* %s\n"+
			"📞 *Phone:* %s\n"+
			"📍 *From:* %s\n"+
			"🏁 *To:* %s\n"+
			"📅 *Date:* %s\n"+
			"🚚 *Movers:* %d Fellas\n"+
			"⏱ *Hours:* %d hrs\n"+
			"💰 *Estimated Cost:* $%d\n"+
			"📝 *Details:* %s",
		lead.Name,
		lead.Phone,
		from,
		to,
		lead.MovingDate,
		lead.FellasNumber,
		lead.Hours,
		lead.TotalPrice,
		details,
	)

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	resp, err := http.PostForm(apiURL, url.Values{
		"chat_id":    {chatID},
		"text":       {text},
		"parse_mode": {"Markdown"},
	})
	if err != nil {
		log.Printf("[Telegram] failed to send notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[Telegram] unexpected status: %s. Check your Token and ChatID", resp.Status)
	}
}
