package main

import (
	"checking/models"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/users", models.UsersGETHandler).Methods("GET")
	r.HandleFunc("/users", models.UserPOSTHandler).Methods("POST")
	r.HandleFunc("/users", models.UsersDELETEHandler).Methods("DELETE")
	r.HandleFunc("/users/{id}", models.UserGETHandler).Methods("GET")
	r.HandleFunc("/users/{id}", models.UserDELETEHandler).Methods("DELETE")
	r.HandleFunc("/users/{id}", models.UserPUTHandler).Methods("PUT")

	r.HandleFunc("/sites", models.GrabSitesGETHandler).Methods("GET")
	r.HandleFunc("/sites", models.GrabSitePOSTHandler).Methods("POST")
	r.HandleFunc("/sites", models.GrabSitesDELETEHandler).Methods("DELETE")
	r.HandleFunc("/sites/{id}", models.GrabSiteGETHandler).Methods("GET")
	r.HandleFunc("/sites/{id}", models.GrabSiteDELETEHandler).Methods("DELETE")
	r.HandleFunc("/sites/{id}", models.GrabSitePUTHandler).Methods("PUT")

	r.HandleFunc("/main_texts", models.MainTextsGETHandler).Methods("GET")
	r.HandleFunc("/main_texts", models.MainTextPOSTHandler).Methods("POST")
	r.HandleFunc("/main_texts", models.MainTextsDELETEHandler).Methods("DELETE")
	r.HandleFunc("/main_texts/{id}", models.MainTextGETHandler).Methods("GET")
	r.HandleFunc("/main_texts/{id}", models.MainTextDELETEHandler).Methods("DELETE")
	r.HandleFunc("/main_texts/{id}", models.MainTextPUTHandler).Methods("PUT")

	r.HandleFunc("/page_sites", models.PageSitesGETHandler).Methods("GET")
	r.HandleFunc("/page_sites", models.PageSitePOSTHandler).Methods("POST")
	r.HandleFunc("/page_sites", models.PageSitesDELETEHandler).Methods("DELETE")
	r.HandleFunc("/page_sites/{id}", models.PageSiteGETHandler).Methods("GET")
	r.HandleFunc("/page_sites/{id}", models.PageSiteDELETEHandler).Methods("DELETE")
	r.HandleFunc("/page_sites/{id}", models.PageSitePUTHandler).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", r))


}

