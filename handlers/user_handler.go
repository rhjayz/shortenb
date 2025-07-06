package handlers

import (
	"net/http"
	"shortb/config"
	"shortb/middleware"
	"shortb/models"
	"shortb/utils"
)

func IndexUser(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		user = models.User{}
	}

	var Users models.User
	err := config.DB.First(&Users, user.ID).Error
	if err == nil {
		http.Error(w, "User Undefined", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"User": user,
		"Data": Users,
	}

	utils.RenderTemplate(w, "views/pages/user/index.html", data)

}
