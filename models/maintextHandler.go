package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)
type MainText struct {
	Id       	string 	`json:"id"`
	User_id    	string 	`json:"user_id"`
	Page_id 	string 	`json:"page_id"`
	PageDate    string 	`json:"page_date"`
}

type MainTextExpand struct {
	Id       	string 			`json:"id"`
	User_id    	Users 			`json:"user_id"`
	Page_id 	PageSiteExpand 	`json:"page_id"`
	PageDate    string 			`json:"page_date"`
}

var maintext1 []MainText
var maintextExpand1 []MainTextExpand

func MainTextsGETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	maintext1 = nil
	maintextExpand1 = nil
	_, ok := r.URL.Query()["expand"]

	if !ok  {

		rows, err := db.Query("SELECT * FROM maintext ")
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var text MainText
			err = rows.Scan(&text.Id, &text.User_id, &text.Page_id, &text.PageDate)
			if err != nil {
				panic(err)
			}
			maintext1 = append(maintext1, MainText{text.Id, text.User_id, text.Page_id, text.PageDate})
		}
		textBytes, _ := json.MarshalIndent(maintext1, "", "\t")

		w.Header().Set("Content-Type", "application/json")
		w.Write(textBytes)

		defer rows.Close()
		defer db.Close()
	}else{
		rows, err := db.Query("SELECT * FROM maintext")
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var text MainText
			err = rows.Scan(&text.Id, &text.User_id, &text.Page_id, &text.PageDate)
			if err != nil {
				panic(err)
			}

			rowsUserExpand, err := db.Query("SELECT * FROM users where id=$1", text.User_id)
			if err != nil {
				log.Fatal(err)
			}

			rowsPageExpand, err := db.Query("SELECT * FROM pagesite where id=$1", text.Page_id)
			if err != nil {
				log.Fatal(err)
			}

			if rowsUserExpand.Next() {
				var userT Users
				err = rowsUserExpand.Scan(&userT.Id, &userT.Login, &userT.Password)
				if err != nil {
					panic(err)
				}
				if rowsPageExpand.Next() {
					var page PageSite
					err = rowsPageExpand.Scan(&page.Id, &page.Site_id, &page.Page_link, &page.Namepage, &page.Grabtext)
					if err != nil {
						panic(err)
					}
					rowsSiteExpand, err := db.Query("SELECT * FROM grabsite where id=$1", page.Site_id)
					if err != nil {
						log.Fatal(err)
					}
					if rowsSiteExpand.Next() {
						var site GrabSite
						err = rowsSiteExpand.Scan(&site.Id, &site.NameSite, &site.URL)
						if err != nil {
							panic(err)
						}

						maintextExpand1 = append(maintextExpand1, MainTextExpand{text.Id,
							Users{userT.Id, userT.Login, userT.Password},
							PageSiteExpand{page.Id, GrabSite{site.Id, site.NameSite, site.URL}, page.Page_link,
								page.Namepage, page.Grabtext}, text.PageDate})
						textBytes, _ := json.MarshalIndent(maintextExpand1, "", "\t")

						w.Header().Set("Content-Type", "application/json")
						w.Write(textBytes)
						defer rowsSiteExpand.Close()
					}
					defer rowsPageExpand.Close()
				}
				defer rowsUserExpand.Close()
			}
		}
		defer rows.Close()
		defer db.Close()
	}
}

func MainTextPOSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var p MainText
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO maintext ("user_id(FK)", "page_id(FK)", pagedate) VALUES ($1, $2, $3)`
	_, err = db.Exec(sqlStatement, p.User_id, p.Page_id, p.PageDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func MainTextsDELETEHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	_, err := db.Exec("DELETE FROM maintext")

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func MainTextGETHandler(w http.ResponseWriter, r *http.Request)  {
	db := OpenConnection()
	maintext1 = nil
	maintextExpand1 = nil
	params := mux.Vars(r)
	idGet := params["id"]

	_, ok := r.URL.Query()["expand"]
	fmt.Println(ok)
	if !ok {
		row, err := db.Query("SELECT * FROM maintext where id=$1", idGet)
		if err != nil {
			log.Fatal(err)
		}

		var text MainText
		if row.Next(){
			err = row.Scan(&text.Id, &text.User_id, &text.Page_id, &text.PageDate)
			if err != nil{
				panic(err)
			}
			userId := MainText{text.Id, text.User_id, text.Page_id, text.PageDate}

			textBytes, _ := json.Marshal(userId)

			w.Header().Set("Content-Type", "application/json")
			w.Write(textBytes)
		}

		defer row.Close()
		defer db.Close()
	} else {
		row, err := db.Query("SELECT * FROM maintext where id=$1", idGet)
		if err != nil {
			log.Fatal(err)
		}
		if row.Next() {
			var text MainText
			err = row.Scan(&text.Id, &text.User_id, &text.Page_id, &text.PageDate)
			if err != nil {
				panic(err)
			}

			rowsUserExpand, err := db.Query("SELECT * FROM users where id=$1", text.User_id)
			if err != nil {
				log.Fatal(err)
			}

			rowsPageExpand, err := db.Query("SELECT * FROM pagesite where id=$1", text.Page_id)
			if err != nil {
				log.Fatal(err)
			}

			if rowsUserExpand.Next() {
				var userT Users
				err = rowsUserExpand.Scan(&userT.Id, &userT.Login, &userT.Password)
				if err != nil {
					panic(err)
				}
				if rowsPageExpand.Next() {
					var page PageSite
					err = rowsPageExpand.Scan(&page.Id, &page.Site_id, &page.Page_link, &page.Namepage, &page.Grabtext)
					if err != nil {
						panic(err)
					}
					rowsSiteExpand, err := db.Query("SELECT * FROM grabsite where id=$1", page.Site_id)
					if err != nil {
						log.Fatal(err)
					}
					if rowsSiteExpand.Next() {
						var site GrabSite
						err = rowsSiteExpand.Scan(&site.Id, &site.NameSite, &site.URL)
						if err != nil {
							panic(err)
						}

						maintextExpand1 = append(maintextExpand1, MainTextExpand{text.Id,
							Users{userT.Id, userT.Login, userT.Password},
							PageSiteExpand{page.Id, GrabSite{site.Id, site.NameSite, site.URL}, page.Page_link,
								page.Namepage, page.Grabtext}, text.PageDate})
						textBytes, _ := json.MarshalIndent(maintextExpand1, "", "\t")

						w.Header().Set("Content-Type", "application/json")
						w.Write(textBytes)
						defer rowsSiteExpand.Close()
					}
					defer rowsPageExpand.Close()
				}
				defer rowsUserExpand.Close()
			}
		}
	defer row.Close()
	defer db.Close()
	}

}

func MainTextDELETEHandler(w http.ResponseWriter, r *http.Request)  {
	db := OpenConnection()
	params := mux.Vars(r)
	idGet := params["id"]
	_, err := db.Exec("DELETE FROM maintext where id=$1", idGet)

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func MainTextPUTHandler(w http.ResponseWriter, r *http.Request)  {
	db := OpenConnection()
	params := mux.Vars(r)
	idGet := params["id"]

	var p MainText

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `UPDATE maintext SET "user_id(FK)" = $1, "page_id(FK)" = $2, pagedate = $3 WHERE id = $4`
	_, err = db.Exec(sqlStatement, p.User_id, p.Page_id, p.PageDate, idGet)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)

	defer db.Close()
}
