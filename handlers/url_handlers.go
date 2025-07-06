package handlers

import (
	"net/http"
	"shortb/middleware"
	"shortb/models"
	"shortb/utils"
)

func IndexUrl(w http.ResponseWriter, r *http.Request) {

	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		user = models.User{}
	}

	data := map[string]interface{}{
		"User": user, // bisa kosong kalau belum login
	}

	utils.RenderTemplate(w, "views/index.html", data)

}
