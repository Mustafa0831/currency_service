package handlers

import (
	"currency_service/models"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitHandlers(dbConn *sqlx.DB) {
	db = dbConn
}

type Rates struct {
	XMLName xml.Name `xml:"rates"`
	Items   []Item   `xml:"item"`
	Date    string   `xml:"date"`
}

type Item struct {
	FullName    string `xml:"fullname"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
}

// SaveCurrency saves currency rates from the National Bank of Kazakhstan to the database.
//
// @Summary Save currency rates
// @Description Saves currency rates from the National Bank of Kazakhstan for a specific date.
// @Tags currency
// @Accept json
// @Produce json
// @Param date query string true "Date in DD.MM.YYYY format"
// @Success 200 {object} map[string]bool "Currency rates successfully saved"
// @Failure 400 {object} map[string]bool "Invalid date format"
// @Failure 500 {object} map[string]bool "Internal server error"
func SaveCurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]

	parsedDate, err := time.Parse("02.01.2006", date)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}
	formattedDate := parsedDate.Format("2006-01-02")

	go func() {
		url := fmt.Sprintf("https://nationalbank.kz/rss/get_rates.cfm?fdate=%s", date)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error fetching data: %v", err)
			return
		}
		defer resp.Body.Close()

		var rates Rates
		if err := xml.NewDecoder(resp.Body).Decode(&rates); err != nil {
			log.Printf("Error decoding XML: %v", err)
			return
		}

		tx := db.MustBegin()
		for _, item := range rates.Items {
			value, err := strconv.ParseFloat(item.Description, 64)
			if err != nil {
				log.Printf("Error parsing value: %v", err)
				continue
			}

			currency := models.Currency{
				Title: item.FullName,
				Code:  item.Title,
				Value: value,
				ADate: formattedDate,
			}
			tx.NamedExec(`INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (:TITLE, :CODE, :VALUE, :A_DATE)`, &currency)
		}
		tx.Commit()
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// GetCurrency retrieves currency rates for a specific date and code.
//
// @Summary Get currency rates
// @Description Retrieves currency rates for a specific date and code.
// @Tags currency
// @Accept json
// @Produce json
// @Param date query string true "Date in DD.MM.YYYY format"
// @Param code query string false "Currency code"
// @Success 200 {array} models.Currency "Currency rates successfully retrieved"
// @Failure 400 {object} map[string]bool "Invalid date format"
// @Failure 500 {object} map[string]bool "Internal server error"
func GetCurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]
	code := vars["code"]

	parsedDate, err := time.Parse("02.01.2006", date)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}
	formattedDate := parsedDate.Format("2006-01-02")

	var currencies []models.Currency
	query := "SELECT * FROM R_CURRENCY WHERE A_DATE = :A_DATE"
	namedArgs := map[string]interface{}{
		"A_DATE": formattedDate,
	}
	if code != "" {
		query += " AND CODE = :CODE"
		namedArgs["CODE"] = code
	}

	rows, err := db.NamedQuery(query, namedArgs)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Database query error: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var currency models.Currency
		if err := rows.StructScan(&currency); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Database scan error: %v", err)
			return
		}
		currencies = append(currencies, currency)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Database rows error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currencies)
}
