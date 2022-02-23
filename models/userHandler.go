package models

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)
type Users struct {
	Id       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

var user1 []Users

func UsersGETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	user1 = nil

	keyLogin, okL := r.URL.Query()["login"]
	keyPassword, okP := r.URL.Query()["password"]

	if (!okL || len(keyLogin[0]) < 1 )&&(!okP || len(keyPassword[0]) < 1 ){
		rows, err := db.Query("SELECT * FROM users ORDER BY id")
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var person Users
			err = rows.Scan(&person.Id, &person.Login, &person.Password)
			if err!=nil{
				panic(err)
			}
			user1 = append(user1, Users{person.Id, person.Login, person.Password})
		}
		peopleBytes, _ := json.MarshalIndent(user1,"", "\t")

		w.Header().Set("Content-Type", "application/json")
		w.Write(peopleBytes)

		defer rows.Close()
		defer db.Close()
	} else {
		if !okL || len(keyLogin[0]) < 1{
			key := keyPassword[0]
			row, err := db.Query("SELECT * FROM users where password=$1", key)
			if err != nil {
				log.Fatal(err)
			}

			var person Users
			if row.Next(){
				err = row.Scan(&person.Id, &person.Login, &person.Password)
				if err != nil{
					panic(err)
				}
				userId := Users{person.Id, person.Login, person.Password}

				peopleBytes, _ := json.Marshal(userId)

				w.Header().Set("Content-Type", "application/json")
				w.Write(peopleBytes)
			}

			defer row.Close()
			defer db.Close()
		} else {
			key := keyLogin[0]
			row, err := db.Query("SELECT * FROM users where login=$1", key)
			if err != nil {
				log.Fatal(err)
			}

			var person Users
			if row.Next(){
				err = row.Scan(&person.Id, &person.Login, &person.Password)
				if err != nil{
					panic(err)
				}
				userId := Users{person.Id, person.Login, person.Password}

				peopleBytes, _ := json.Marshal(userId)

				w.Header().Set("Content-Type", "application/json")
				w.Write(peopleBytes)
			}


			defer row.Close()
			defer db.Close()
		}
	}
	// Query()["key"] will return an array of items,
	// we only want the single item.
}

func UserPOSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var p Users
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO users (login, password) VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, p.Login, p.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func UsersDELETEHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	_, err := db.Exec("DELETE FROM users")

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func UserGETHandler(w http.ResponseWriter, r *http.Request)  {
	db := OpenConnection()
	params := mux.Vars(r)
	idGet := params["id"]
	row, err := db.Query("SELECT * FROM users where id=$1", idGet)
	if err != nil {
		log.Fatal(err)
	}

	var person Users
	if row.Next(){
		err = row.Scan(&person.Id, &person.Login, &person.Password)
		if err != nil{
			panic(err)
		}
		userId := Users{person.Id, person.Login, person.Password}

		peopleBytes, _ := json.Marshal(userId)

		w.Header().Set("Content-Type", "application/json")
		w.Write(peopleBytes)
	}


	defer row.Close()
	defer db.Close()
}

func UserDELETEHandler(w http.ResponseWriter, r *http.Request)  {
	db := OpenConnection()
	params := mux.Vars(r)
	idGet := params["id"]
	_, err := db.Exec("DELETE FROM users where id=$1", idGet)

	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func UserPUTHandler(w http.ResponseWriter, r *http.Request)  {
	db := OpenConnection()
	params := mux.Vars(r)
	idGet := params["id"]

	var p Users

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `UPDATE users SET login = $1, password = $2 WHERE id=$3`
	_, err = db.Exec(sqlStatement, p.Login, p.Password, idGet)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)


	defer db.Close()
}
