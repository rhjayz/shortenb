package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"shortb/config"
	"shortb/middleware"
	"shortb/models"
	"shortb/utils"
	"strconv"
	"time"
)

type FilterDate struct {
	SDate string `json:"start_date"`
	EDate string `json:"end_date"`
}

type FilterMonth struct {
	Month string `json:"month"`
	Year  string `json:"year"`
}

type ShortlinkMonthlySummary struct {
	ShortCode   string `json:"shortcode" gorm:"column:short_code"`
	OriginalURL string `json:"originalurl" gorm:"column:original_url"`
	TotalClicks int64  `json:"totalclicks" gorm:"column:totalclicks"`
}

func IndexReport(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		user = models.User{}
	}

	data := map[string]interface{}{
		"User": user,
	}

	utils.RenderTemplate(w, "views/pages/report/index.html", data)
}

func ReportFilterDate(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		user = models.User{}
	}

	var input FilterDate
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.SDate == "" || input.EDate == "" {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	if input.SDate > input.EDate {
		http.Error(w, "Invalid date range", http.StatusBadRequest)
		return
	}

	var data []models.Shortlink
	err = config.DB.Preload("Clicklogs", "created_at BETWEEN ? AND ?", input.SDate, input.EDate).Where("user_id", user.ID).Find(&data).Error
	if err != nil {
		http.Error(w, "Failed to get data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"type": "detail",
		"data": data,
	})
}

func ReportFilterMonth(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		user = models.User{}
	}

	var input FilterMonth
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Year == "" || input.Month == "" {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	yearInt, err := strconv.Atoi(input.Year)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	monthInt, err := strconv.Atoi(input.Month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}

	startDate := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	var data []ShortlinkMonthlySummary

	err = config.DB.
		Debug().
		Model(&models.Clicklogs{}).
		Joins("JOIN shortlinks ON shortlinks.id = clicklogs.shortlink_id").
		Where("shortlinks.user_id = ? AND clicklogs.created_at BETWEEN ? AND ?", user.ID, startDate, endDate).
		Select("shortlinks.short_code, shortlinks.original_url,  COUNT(*) as totalclicks, MAX(clicklogs.created_at) as created_at").
		Group("shortlinks.short_code, shortlinks.original_url").
		Order("created_at DESC").
		Scan(&data).Error

	if err != nil {
		http.Error(w, "Data undefined", http.StatusBadRequest)
	}

	log.Printf("Data hasil query: %+v\n", data)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"type": "summary",
		"data": data,
	})
}
