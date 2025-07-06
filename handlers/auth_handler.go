package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"shortb/config"
	"shortb/models"
	"shortb/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RegisterProps struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginProps struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func generateToken() string {
	b := make([]byte, 32) // 256-bit token
	rand.Read(b)
	return hex.EncodeToString(b)
}

func LoginUrl(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "views/pages/login.html", nil)
}

func RegisterUrl(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "views/pages/register.html", nil)
}

func RegisterStore(w http.ResponseWriter, r *http.Request) {
	var input RegisterProps
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Email == "" || input.Password == "" || input.Name == "" {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	data := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	result := config.DB.Create(&data)
	if result.Error != nil {
		http.Error(w, "Email already used ", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Registered successfully",
	})

}

func LoginStore(w http.ResponseWriter, r *http.Request) {
	var input LoginProps
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Email == "" || input.Password == "" {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		http.Error(w, "Email not found", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
		return
	}

	token := generateToken()

	session := models.Session{
		UserID:    user.ID,
		Token:     token,
		UserAgent: r.UserAgent(),
		IPAddress: r.RemoteAddr,
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	state := config.DB.Create(&session)
	if state == nil {
		http.Error(w, "Something Wrong password", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		Expires:  session.ExpiredAt,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login Succesfull",
	})

}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Ini yang bikin cookie langsung expired
		HttpOnly: true,
	})

	token, err := r.Cookie("session_token")
	if err == nil {
		config.DB.Where("token = ?", token.Value).Delete(&models.Session{})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out",
	})

}
