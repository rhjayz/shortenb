package handlers

import (
	"encoding/json"
	"net/http"
	"shortb/config"
	"shortb/middleware"
	"shortb/models"
	"shortb/utils"
	"strconv"

	"github.com/gorilla/mux"
)

type ShortlinkProps struct {
	OriginalURL string `json:"originalurl"`
	CDomain     string `json:"customdomain"`
}

type ShortlinkUpdateProps struct {
	CDomain string `json:"customdomain"`
}

type ShortlinkWithClicks struct {
	Shortlink  models.Shortlink
	ClickCount int64
}

func IndexShortLink(w http.ResponseWriter, r *http.Request) {

	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		user = models.User{}
	}

	var links []models.Shortlink
	err := config.DB.Where("user_id = ?", user.ID).Find(&links).Error
	if err != nil {
		http.Error(w, "Failed to get shortlinks", http.StatusInternalServerError)
		return
	}

	var dataWithClicks []ShortlinkWithClicks

	for _, link := range links {
		var count int64
		err := config.DB.Model(&models.Clicklogs{}).Where("shortlink_id = ?", link.ID).Count(&count).Error
		if err != nil {
			http.Error(w, "Failed to count clicks", http.StatusInternalServerError)
			return
		}
		dataWithClicks = append(dataWithClicks, ShortlinkWithClicks{
			Shortlink:  link,
			ClickCount: count,
		})
	}

	data := map[string]interface{}{
		"User": user,
		"Data": dataWithClicks,
	}

	utils.RenderTemplate(w, "views/pages/shortlink/index.html", data)

}

func ShortLinkDetail(w http.ResponseWriter, r *http.Request) {
	idParams := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idParams, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		user = models.User{}
	}

	var shortlink models.Shortlink
	err = config.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&shortlink).Error
	if err != nil {
		http.Error(w, "Data Undefined", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"User": user,
		"Data": shortlink,
	}
	utils.RenderTemplate(w, "views/pages/shortlink/detail.html", data)

}

func ShortLinkStore(w http.ResponseWriter, r *http.Request) {

	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		user = models.User{}
	}

	var input ShortlinkProps
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.OriginalURL == "" || input.CDomain == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	data := models.Shortlink{
		UserID:       user.ID,
		OriginalUrl:  &input.OriginalURL,
		ShortCode:    utils.GenerateShortCode(6),
		CustomDomain: input.CDomain,
	}

	state := config.DB.Create(&data)
	if state.Error != nil {
		http.Error(w, "Something Wrong ", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Add successfull",
	})
}

func ShortLinkDestroy(w http.ResponseWriter, r *http.Request) {
	idParams := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(idParams, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		user = models.User{}
	}

	var shortlink models.Shortlink
	state := config.DB.Where("id = ? AND user_id = ?", id, user.ID).Delete(&shortlink)
	if state.RowsAffected == 0 {
		http.Error(w, "Shortlink not found or not yours", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Shortlink deleted successfully",
	})

}

func ShortLinkUpdate(w http.ResponseWriter, r *http.Request) {
	idParams := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idParams, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		user = models.User{}
	}

	var input ShortlinkUpdateProps
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.CDomain == "" {
		http.Error(w, "Something Wrong", http.StatusBadRequest)
		return
	}

	data := models.Shortlink{
		CustomDomain: input.CDomain,
	}
	state := config.DB.Model(&models.Shortlink{}).Where("id = ? AND user_id = ?", id, user.ID).Updates(&data)

	if state.Error != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Update successfull",
	})

}
