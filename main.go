package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const DEFAULT_PORT = "8000"

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	err := godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT
	}
	dbUrl := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	log.Printf("Connecting to db %v...\n", dbUrl)
	db, err = sql.Open("mysql", dbUrl)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	log.Printf("Server listening on port %v...\n", port)
	http.HandleFunc("/albums", handleList)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func handleList(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling request %v...\n", r)
	albums, err := getAlbums()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(albums)
}

func getAlbums() ([]Album, error) {
	log.Printf("Reading albums from db...\n")
	var albums []Album
	log.Printf("%v\n", db)
	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist: %v", err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist: %v", err)
	}
	return albums, nil
}
