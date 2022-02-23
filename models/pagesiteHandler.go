package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type PageSite struct {
	Id       	string 	`json:"id"`
	Site_id    	string	`json:"site_id"`
	Page_link 	string	`json:"page_link"`
	Namepage    string 	`json:"namepage"`
	Grabtext	string 	`json:"grabtext""`
}

type PageSiteExpand struct {
	Id       	string 		`json:"id"`
	Site_id    	string		`json:"site_id"`
	Page_link 	string		`json:"page_link"`
	Namepage    string 		`json:"namepage"`
	Grabtext	string 		`json:"grabtext"`
	Site_expand GrabSite	`json:"site_exp"`
}

func ToExpand (w http.ResponseWriter, db *sql.DB) {
	rows, err := db.Query("SELECT * FROM pagesite ORDER BY id")
	if err != nil {
		log.Fatal(err)
	}

	var pageBytes []byte

	for rows.Next() {
		var page PageSite
		err = rows.Scan(&page.Id, &page.Site_id, &page.Page_link, &page.Namepage, &page.Grabtext)
		if err != nil {
			panic(err)
		}

		rowsExpand, err := db.Query("SELECT * FROM grabsite where id=$1", page.Site_id)
		if err != nil {
			log.Fatal(err)
		}

		if rowsExpand.Next() {
			var site GrabSite
			err = rowsExpand.Scan(&site.Id, &site.NameSite, &site.URL)
			if err != nil {
				panic(err)
			}
			pagesiteExpand1 = append(pagesiteExpand1, PageSiteExpand{page.Id,
				page.Site_id, page.Page_link, page.Namepage, page.Grabtext, GrabSite{site.Id, site.NameSite, site.URL}})
			defer rowsExpand.Close()
		}

		pageBytes, _ = json.MarshalIndent(pagesiteExpand1, "", "\t")

		w.Header().Set("Content-Type", "application/json")

		defer rows.Close()
	}
	w.Write(pageBytes)
	defer db.Close()
}

var pagesite1 []PageSite
var pagesiteExpand1 []PageSiteExpand

func PageSitesGETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	pagesite1 = nil
	pagesiteExpand1 = nil
	_, ok := r.URL.Query()["expand"]
	if !ok {
		rows, err := db.Query("SELECT * FROM pagesite ORDER BY id")
		if err != nil {
			log.Fatal(err)
		}
		var pageBytes []byte
		for rows.Next() {
			var page PageSite
			err = rows.Scan(&page.Id, &page.Site_id, &page.Page_link, &page.Namepage, &page.Grabtext)
			if err != nil {
				panic(err)
			}
			pagesite1 = append(pagesite1, PageSite{page.Id, page.Site_id, page.Page_link, page.Namepage, page.Grabtext})
		}
		pageBytes, _ = json.MarshalIndent(pagesite1, "", "\t")

		w.Header().Set("Content-Type", "application/json")
		w.Write(pageBytes)

		defer rows.Close()
		defer db.Close()
	} else {
		ToExpand(w, db)
	}
}

func PageSitePOSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var p PageSite
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO pagesite ("site_id(FK)", "page_link", namepage, grabtext) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, p.Site_id, p.Page_link, p.Namepage, p.Grabtext)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func PageSitesDELETEHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	_, err := db.Exec("DELETE FROM pagesite")

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func PageSiteGETHandler(w http.ResponseWriter, r *http.Request)  {
	db := OpenConnection()
	maintext1 = nil
	maintextExpand1 = nil
	params := mux.Vars(r)
	idGet := params["id"]
	_, ok := r.URL.Query()["expand"]
	if !ok{
		row, err := db.Query("SELECT * FROM pagesite where id=$1", idGet)
		if err != nil {
			log.Fatal(err)
		}

		var page PageSite
		if row.Next(){
			err = row.Scan(&page.Id, &page.Site_id, &page.Page_link, &page.Namepage, &page.Grabtext)
			if err != nil{
				panic(err)
			}
			userId := PageSite{page.Id, page.Site_id, page.Page_link, page.Namepage, page.Grabtext}

			pageBytes, _ := json.Marshal(userId)
			w.Header().Set("Content-Type", "application/json")
			w.Write(pageBytes)
		}

		defer row.Close()
		defer db.Close()
	} else {
		row, err := db.Query("SELECT * FROM pagesite where id=$1", idGet)
		if err != nil {
			log.Fatal(err)
		}

		if row.Next() {
			var page PageSite
			err = row.Scan(&page.Id, &page.Site_id, &page.Page_link, &page.Namepage, &page.Grabtext)
			if err != nil {
				panic(err)
			}
				rowsExpand, err := db.Query("SELECT * FROM grabsite where id=$1", page.Site_id)
				if err != nil {
					log.Fatal(err)
				}
					if rowsExpand.Next() {
						var site GrabSite
						err = rowsExpand.Scan(&site.Id, &site.NameSite, &site.URL)
						if err != nil {
							panic(err)
						}
						pagesite1 = append(pagesite1, PageSite{page.Id, page.Site_id, page.Page_link, page.Namepage, page.Grabtext})


						defer rowsExpand.Close()
				}
			pageBytes, _ := json.MarshalIndent(pagesiteExpand1, "", "\t")
			w.Header().Set("Content-Type", "application/json")
			w.Write(pageBytes)
			defer row.Close()
			}
		defer db.Close()
	}
}

func PageSiteDELETEHandler(w http.ResponseWriter, r *http.Request)  {
	db := OpenConnection()
	params := mux.Vars(r)
	idGet := params["id"]
	fmt.Println(idGet)
	_, err := db.Exec("DELETE FROM pagesite where id=$1", idGet)

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func PageSitePUTHandler(w http.ResponseWriter, r *http.Request)  {
	fmt.Print(1)
	db := OpenConnection()
	params := mux.Vars(r)
	idGet := params["id"]

	var p PageSite

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `UPDATE pagesite SET "site_id(FK)" = $1, page_link = $2, namepage = $3, grabtext = $4 
	WHERE id = $5`
	_, err = db.Exec(sqlStatement, p.Site_id, p.Page_link, p.Namepage, p.Grabtext, idGet)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)

	defer db.Close()
}
