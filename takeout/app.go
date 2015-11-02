package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {

	http.HandleFunc("/", pageHandler)
	http.HandleFunc("/scripts/", scriptHandler)

	http.HandleFunc("/menu/items", menuItemHandler)

	http.ListenAndServe(":8080", nil)
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/index.html")
}

func scriptHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	http.ServeFile(w, r, r.URL.Path[1:])
}

var _cachedPostgresConnection *sql.DB

func getPostgresConnection() *sql.DB {
	if _cachedPostgresConnection == nil {
		conn, err := sql.Open("postgres", "user=hnetzer dbname=takeout sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}
		_cachedPostgresConnection = conn
	}
	return _cachedPostgresConnection
}
