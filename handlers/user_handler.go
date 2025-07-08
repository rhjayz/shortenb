package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shortb/config"
	"shortb/middleware"
	"shortb/models"
	"shortb/utils"

	"golang.org/x/crypto/bcrypt"
)

func IndexUser(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		user = models.User{}
	}

	var Users models.User
	err := config.DB.First(&Users, user.ID).Error
	if err != nil {
		http.Error(w, "User Undefined", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"User": user,
		"Data": Users,
	}

	utils.RenderTemplate(w, "views/pages/user/index.html", data)

}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">> UpdateUser hit")
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		user = models.User{}
	}

	fmt.Println(">> user:", user)

	err := r.ParseMultipartForm(10 << 20) // max 10MB
	if err != nil {
		http.Error(w, "Can't parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if name == "" || email == "" {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	fmt.Println(">> name:", name, "email:", email, "password:", password)

	hashedPassword := user.Password
	if password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		hashedPassword = string(hashed)
	}

	oldImage := ""
	if user.Image != nil {
		oldImage = *user.Image
	}

	fmt.Println(">> oldImage:", oldImage)

	newImagePath, err := utils.UploadImage(r, "image", oldImage)
	if err != nil {
		http.Error(w, "Image upload failed", http.StatusInternalServerError)
		return
	}

	fmt.Println(">> newImagePath:", newImagePath)

	data := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Image:    &newImagePath,
	}

	state := config.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(&data)

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
