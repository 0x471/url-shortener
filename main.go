package main

import (
	"fmt"
	"log"
	"net/http"
	"database/sql"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	generators "github.com/0x471/url-shortener/generators"
	dbOp "github.com/0x471/url-shortener/db"

)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
)

type dbHandler struct {
	Conn *sql.DB
}

func shortenHandler(db *dbHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlArr, ok := r.URL.Query()["url"]

		if !ok || len(urlArr[0]) < 1 {
			w.WriteHeader(500)
			fmt.Fprintln(w, "URL parameter is missing")
			return
		}
		url := urlArr[0]
		isValidUrl := govalidator.IsURL(url)

		if isValidUrl != true {
			w.WriteHeader(500)
			fmt.Fprintln(w, "Url format is not OK")
			return
		}

		id := generators.GenerateUniqueID(url)
		log.Printf("[INFO] URL: %s Adler-32 Checksum: %s", url, id)	
		
		if dbOp.SearchObj(db.Conn, id) == "nil" {
			dbOp.InsertUrl(db.Conn, url, id)
			fmt.Fprintf(w, "/s/%s", id)
			return
		}

		
		fmt.Fprintf(w, "/s/%s", id)
		

	}
}
	    

func searchHandler(db *dbHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["adler32"]

		searchResult := dbOp.SearchObj(db.Conn, id)

		if searchResult == "err" {
			fmt.Fprintf(w, "url not found.")
			return
		}
		url := searchResult
		http.Redirect(w, r, url, 302)
	}
}

func main() {
	log.Println("[INFO] Checking database...")

	database, err := sql.Open("sqlite3", "url-shortener.db")
	log.Println("[INFO] Opening database...")
	if err != nil {
		log.Println("[ERROR] An error occured while opening database:",err)
	}
	defer database.Close()

	dbOp.CheckTable(database)
	db := &dbHandler{Conn: database}

	log.Println("[INFO] Starting HTTP server...")
	router := mux.NewRouter()
	router.Handle("/", http.FileServer(http.Dir("./static")))
	router.Handle("/s", http.RedirectHandler("/", 302))
	router.Handle("/s/{adler32}", searchHandler(db))
	router.Handle("/shorten", shortenHandler(db))
	
	err = http.ListenAndServe(CONN_HOST+":"+CONN_PORT, router)
	if err != nil {
		log.Fatalf("[ERROR] An error occured while starting HTTP server: %s", err)
		return 
	}


	
}
