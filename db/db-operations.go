package db

import (
	"fmt"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func CheckTable(database *sql.DB) {
	
	createTableSqlQuery := "CREATE TABLE IF NOT EXISTS urlshortener (url TEXT NOT NULL, adler32checksum TEXT NOT NULL)";
	createTableStatement, err := database.Prepare(createTableSqlQuery)
	if err != nil {
		log.Println("[ERROR] ",err)
	}
	log.Println("[INFO] Executing query for checking table...")
	createTableStatement.Exec()

	log.Println("[INFO] Table checked successfully.")

}

func SearchObj(database *sql.DB, adler32 string) string {

	searchObjSqlQuery := fmt.Sprintf("SELECT url FROM urlshortener WHERE adler32checksum=('%s')", adler32)
	log.Println("[INFO] Executing query for searching object...")
	rows, err := database.Query(searchObjSqlQuery)

	if err != nil {
		log.Println("[ERROR] ",err)
		return "err"
	}
	var url string

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&url)
		if err != nil {
			log.Fatal(err)
		}
		return url

    }

	return "nil"
} 

func InsertUrl(database *sql.DB, url, adler32checksum string) {


	insertUrlSqlQuery := "INSERT INTO urlshortener (url, adler32checksum) VALUES (?, ?)";
	inserUrlStatement, err := database.Prepare(insertUrlSqlQuery)
	if err != nil {
		log.Println("[ERROR] ",err)
	}
	log.Println("[INFO] Executing query for inserting data...")
	inserUrlStatement.Exec(url, adler32checksum)

	log.Println("[INFO] Data inserted successfully.")
}
	
