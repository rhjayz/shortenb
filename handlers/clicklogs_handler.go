package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shortb/config"
	"shortb/models"
	"strings"

	"time"

	"github.com/gorilla/mux"
)

type IPApiResponse struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

func RedirectClickLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]
	println("CODE DITERIMA:", code)

	var shortlink models.Shortlink
	if err := config.DB.Where("short_code = ?", code).First(&shortlink).Error; err != nil {
		println("SHORTLINK TIDAK DITEMUKAN:", err.Error())
		http.NotFound(w, r)
		return
	}

	ua := r.UserAgent()
	ip := getIP(r)
	loc := getLocationFromIP(ip)

	log := models.Clicklogs{
		ShortlinkID: shortlink.ID,
		ClickedAt:   time.Now(),
		IPAddress:   ip,
		Useragent:   &ua,
		Location:    &loc,
	}

	config.DB.Create(&log)
	println("REDIRECT KE:", *shortlink.OriginalUrl)
	println("Lokasi KE:", loc)

	http.Redirect(w, r, *shortlink.OriginalUrl, http.StatusFound)
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	// Handle format IP:port
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}
	return ip
}

func getLocationFromIP(ip string) string {
	type IPApiResponse struct {
		City    string `json:"city"`
		Country string `json:"country"`
	}

	if ip == "127.0.0.1" || ip == "::1" || ip == "" || strings.HasPrefix(ip, "172.") || strings.HasPrefix(ip, "192.") || strings.HasPrefix(ip, "10.") {
		return "Bogor, Indonesia"
	}

	url := fmt.Sprintf("https://ipapi.co/%s/json/", ip)
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var result IPApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	location := result.City
	if result.Country != "" {
		location += ", " + result.Country
	}
	return location
}
