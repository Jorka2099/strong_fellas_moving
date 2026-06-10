package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"
)

// ExportLeadsHandler GET /admin/leads/export
func ExportLeadsHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if pool == nil {
			http.Error(w, "Database pool is nil", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// ── 1. Fetch leads через pgxpool ──────────────────────────────────
		rows, err := pool.Query(ctx, `
			SELECT id, created_at, name, phone,
			       moving_from, moving_to, moving_date,
			       fellas_number, hours, total_price, details
			FROM leads
			ORDER BY created_at DESC
		`)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// ── 2. Build workbook ─────────────────────────────────────────────
		f := excelize.NewFile()
		sheet := "Leads"
		f.SetSheetName("Sheet1", sheet)

		// Header row
		headers := []string{
			"ID", "Created At", "Name", "Phone",
			"Moving From", "Moving To", "Moving Date",
			"Movers", "Hours", "Total Price ($)", "Details",
		}

		headerStyle, _ := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Size: 11},
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"C25E52"}, Pattern: 1},
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
			Border: []excelize.Border{
				{Type: "bottom", Color: "FFFFFF", Style: 2},
			},
		})

		// Alternating row style
		rowStyle, _ := f.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{Vertical: "center", WrapText: true},
		})
		altRowStyle, _ := f.NewStyle(&excelize.Style{
			Fill:      excelize.Fill{Type: "pattern", Color: []string{"FDF4F3"}, Pattern: 1},
			Alignment: &excelize.Alignment{Vertical: "center", WrapText: true},
		})

		// Write headers
		cols := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}
		for i, h := range headers {
			cell := cols[i] + "1"
			f.SetCellValue(sheet, cell, h)
			f.SetCellStyle(sheet, cell, cell, headerStyle)
		}
		f.SetRowHeight(sheet, 1, 22)

		// Column widths
		widths := map[string]float64{
			"A": 6, "B": 18, "C": 20, "D": 16,
			"E": 22, "F": 22, "G": 14,
			"H": 8, "I": 8, "J": 14, "K": 40,
		}
		for col, w := range widths {
			f.SetColWidth(sheet, col, col, w)
		}

		// ── 3. Write data rows ────────────────────────────────────────────
		rowNum := 2
		for rows.Next() {
			var (
				id                               int
				createdAt                        time.Time
				name, phone                      string
				movingFrom, movingTo, movingDate *string
				fellas                           int
				hours                            int
				totalPrice                       int
				details                          *string
			)

			if err := rows.Scan(
				&id, &createdAt, &name, &phone,
				&movingFrom, &movingTo, &movingDate,
				&fellas, &hours, &totalPrice, &details,
			); err != nil {
				continue
			}

			values := []interface{}{
				id,
				createdAt.Format("2006-01-02 15:04"),
				name,
				phone,
				ptrToStr(movingFrom),
				ptrToStr(movingTo),
				ptrToStr(movingDate),
				fmt.Sprintf("%d Fellas", fellas),
				fmt.Sprintf("%d hrs", hours),
				fmt.Sprintf("$%d", totalPrice),
				ptrToStr(details),
			}

			style := rowStyle
			if rowNum%2 == 0 {
				style = altRowStyle
			}

			for i, val := range values {
				cell := cols[i] + fmt.Sprint(rowNum)
				f.SetCellValue(sheet, cell, val)
				f.SetCellStyle(sheet, cell, cell, style)
			}
			f.SetRowHeight(sheet, rowNum, 18)
			rowNum++
		}

		f.SetPanes(sheet, &excelize.Panes{
			Freeze:      true,
			Split:       false,
			YSplit:      1,
			TopLeftCell: "A2",
			ActivePane:  "bottomLeft",
		})

		// ── 4. Stream to browser ──────────────────────────────────────────
		filename := fmt.Sprintf("leads_%s.xlsx", time.Now().Format("2006-01-02"))
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		w.Header().Set("Cache-Control", "no-cache")

		if err := f.Write(w); err != nil {
			http.Error(w, "Failed to write Excel file: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func ptrToStr(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return "—"
}
