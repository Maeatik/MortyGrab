package models

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)
type GrabSite struct {
	Id       string `json:"id"`
	NameSite string `json:"nameSite"`
	URL      string `json:"url"`
}

var grabsite1 []GrabSite

func GrabSitesGETHandler	(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	grabsite1 = nil

	keyName, okName := r.URL.Query()["nameSite"]
	keyUrl, okUrl := r.URL.Query()["url"]

	if (!okName || len(keyName[0]) < 1 )&&(!okUrl || len(keyUrl[0]) < 1 ) {
		rows, err := db.Query("SELECT * FROM grabsite ")
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var site GrabSite
			err = rows.Scan(&site.Id, &site.NameSite, &site.URL)
			if err != nil {
				panic(err)
			}
			grabsite1 = append(grabsite1, GrabSite{site.Id, site.NameSite, site.URL})
		}
		siteBytes, _ := json.MarshalIndent(grabsite1, "", "\t")

		w.Header().Set("Content-Type", "application/json")
		w.Write(siteBytes)

		defer rows.Close()
		defer db.Close()
	} else {
		if !okUrl|| len(keyUrl[0]) < 1{
			key := keyName[0]
			row, err := db.Query("SELECT * FROM grabsite where namesite=$1", key)
			if err != nil {
				log.Fatal(err)
			}

			var site GrabSite
			if row.Next(){
				err = row.Scan(&site.Id, &site.NameSite, &site.URL)
				if err != nil{
					panic(err)
				}
				siteName := GrabSite{site.Id, site.NameSite, site.URL}

				siteBytes, _ := json.Marshal(siteName)

				w.Header().Set("Content-Type", "application/json")
				w.Write(siteBytes)
			}

			defer row.Close()
			defer db.Close()
		} else {
			key := keyUrl[0]
			row, err := db.Query("SELECT * FROM grabsite where url=$1", key)
			if err != nil {
				log.Fatal(err)
			}

			var site GrabSite
			if row.Next() {
				err = row.Scan(&site.Id, &site.NameSite, &site.URL)
				if err != nil {
					panic(err)
				}
				siteName := GrabSite{site.Id, site.NameSite, site.URL}

				siteBytes, _ := json.Marshal(siteName)

				w.Header().Set("Content-Type", "application/json")
				w.Write(siteBytes)
			}

			defer row.Close()
			defer db.Close()
		}
	}
}

func GrabSitePOSTHandler	(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var p GrabSite
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO grabsite (namesite, url) VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, p.NameSite, p.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func GrabSitesDELETEHandler	(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	_, err := db.Exec("DELETE FROM grabsite")

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func GrabSiteGETHandler		(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	params := mux.Vars(r)
	idGet := params["id"]
	row, err := db.Query("SELECT * FROM grabsite where id=$1", idGet)
	if err != nil {
		log.Fatal(err)
	}

	var site GrabSite
	if row.Next(){
		err = row.Scan(&site.Id, &site.NameSite, &site.URL)
		if err != nil{
			panic(err)
		}
		userId := GrabSite{site.Id, site.NameSite, site.URL}

		siteBytes, _ := json.Marshal(userId)

		w.Header().Set("Content-Type", "application/json")
		w.Write(siteBytes)
	}


	defer row.Close()
	defer db.Close()
}

func GrabSiteDELETEHandler	(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	params := mux.Vars(r)
	idGet := params["id"]
	_, err := db.Exec("DELETE FROM grabsite where id=$1", idGet)

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func GrabSitePUTHandler		(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	params := mux.Vars(r)
	idGet := params["id"]

	var p GrabSite

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `UPDATE grabsite SET namesite = $1, url = $2 WHERE id=$3`
	_, err = db.Exec(sqlStatement, p.NameSite, p.URL, idGet)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)

	defer db.Close()
}
