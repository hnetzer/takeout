package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

//MenuItem is something you would see on menu
type MenuItem struct {
	MenuItemID   string
	RestaurantID string
	Name         string
	Description  *string
	ImageURL     *string
	Price        *int64
}

func main() {

	http.HandleFunc("/", pageHandler)
	http.HandleFunc("/scripts/", scriptHandler)

	http.HandleFunc("/menu/items", menuItemHandler)

	http.ListenAndServe(":8080", nil)
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("INFO: serving index.html")
	http.ServeFile(w, r, "index.html")
}

func scriptHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving script %s", r.URL.Path[1:])
	w.Header().Set("Content-Type", "text/jsx")
	http.ServeFile(w, r, r.URL.Path[1:])
}

func menuItemHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		restaurantID := r.URL.Query().Get("restaurantID")
		conn := getPostgresConnection()

		items, err := getAllMenuItems(restaurantID, conn)
		if err != nil {
			log.Printf("ERROR: getting menu items: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		encoder := json.NewEncoder(w)
		encoder.Encode(items)
		w.Header().Set("Content-Type", "application/json")
		return

	case "POST":
		decoder := json.NewDecoder(r.Body)
		var item MenuItem
		err := decoder.Decode(&item)
		if err != nil {
			log.Printf("ERROR: decoding menu item: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		conn := getPostgresConnection()

		//we might want to return the menuitemid here
		err = createMenuItem(&item, conn)
		if err != nil {
			log.Printf("ERROR: creating menu item - %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusNotFound)
	}

}

func createMenuItem(item *MenuItem, pgconn *sql.DB) error {

	log.Printf("INFO: creating menu item %+v", item)

	_, err := pgconn.Exec("INSERT INTO menuitems (name, description, price, restaurantid) VALUES ($1, $2, $3, $4)",
		item.Name, item.Description, item.Price, item.RestaurantID)
	if err != nil {
		return err
	}

	return nil
}

func getAllMenuItems(restaurantID string, pgconn *sql.DB) ([]MenuItem, error) {

	log.Printf("INFO: getting menu items for restaurantID %s", restaurantID)

	rows, err := pgconn.Query("SELECT menuitemid, name, description, imageurl, price FROM menuitems WHERE restaurantid = $1", restaurantID)
	if err != nil {
		return nil, err
	}

	menuitems := []MenuItem{}
	for rows.Next() {

		var item MenuItem
		var desc, image sql.NullString
		var price sql.NullInt64
		err = rows.Scan(&item.MenuItemID, &item.Name, &desc, &image, &price)
		if err != nil {
			log.Printf("ERROR: reading menu item - %s", err.Error())
			continue
		}

		if desc.Valid {
			item.Description = &desc.String
		}
		if image.Valid {
			item.ImageURL = &image.String
		}
		if price.Valid {
			item.Price = &price.Int64
		}

		item.RestaurantID = restaurantID

		menuitems = append(menuitems, item)
	}

	return menuitems, nil
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
