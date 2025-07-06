package handlers

import (
	"net/http"
	"shortb/config"
	"shortb/models"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func RedirectClickLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	var shortlink models.Shortlink
	if err := config.DB.Where("short_code = ?", code).First(&shortlink).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	ua := r.UserAgent()
	ip := getIP(r)
	log := models.Clicklogs{
		ShortlinkID: shortlink.ID,
		ClickedAt:   time.Now(),
		IPAddress:   ip,
		Useragent:   &ua,
	}

	config.DB.Create(&log)

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
