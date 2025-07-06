package routes

import (
	"net/http"
	"shortb/handlers"

	"shortb/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.UsertoContextMiddleware)

	fs := http.FileServer(http.Dir("./public/assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Public routes
	r.HandleFunc("/", handlers.IndexUrl)

	// Auth routes
	r.Handle("/login", middleware.BeforeLoginMiddleware(http.HandlerFunc(handlers.LoginUrl))).Methods("GET")
	r.Handle("/register", middleware.BeforeLoginMiddleware(http.HandlerFunc(handlers.RegisterUrl))).Methods("GET")
	r.Handle("/register/store", middleware.BeforeLoginMiddleware(http.HandlerFunc(handlers.RegisterStore))).Methods("POST")
	r.Handle("/login/store", middleware.BeforeLoginMiddleware(http.HandlerFunc(handlers.LoginStore))).Methods("POST")
	r.Handle("/logout", middleware.AfterLoginMiddleware(http.HandlerFunc(handlers.Logout))).Methods("POST")
	//Home
	r.Handle("/home", middleware.AfterLoginMiddleware(http.HandlerFunc(handlers.Index))).Methods("GET")
	//Shorlink
	r.Handle("/shortlink", middleware.AfterLoginMiddleware(http.HandlerFunc(handlers.IndexShortLink))).Methods("GET")
	r.Handle("/shortlink/store", middleware.AfterLoginMiddleware(http.HandlerFunc(handlers.ShortLinkStore))).Methods("POST")
	r.HandleFunc("/{code}", handlers.RedirectClickLogs).Methods("GET")
	r.Handle("/shortlink/destroy/{id}", middleware.AfterLoginMiddleware(http.HandlerFunc(handlers.ShortLinkDestroy))).Methods("POST")
	r.Handle("/shortlink/detail/{id}", middleware.AfterLoginMiddleware(http.HandlerFunc(handlers.ShortLinkDetail))).Methods("GET")
	r.Handle("/shortlink/update/{id}", middleware.AfterLoginMiddleware(http.HandlerFunc(handlers.ShortLinkUpdate))).Methods("POST")
	//User
	r.Handle("/user", middleware.AfterLoginMiddleware(http.HandlerFunc(handlers.IndexUser))).Methods("POST")
	return r
}
