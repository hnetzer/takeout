package main

import (
	"encoding/json"
	"github.com/hnetzer/takeout/menu"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func menuItemHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		getMenuItemsHandler(w, r)
	case "POST":
		postMenuItemHandler(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}

}

func getMenuItemsHandler(w http.ResponseWriter, r *http.Request) {
	restaurantID := r.URL.Query().Get("restaurantID")
	conn := getPostgresConnection()

	items, err := menu.ReadAllItems(restaurantID, conn)
	if err != nil {
		log.Printf("ERROR: getting menu items: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	encoder.Encode(items)
}

func postMenuItemHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var item menu.Item
	err := decoder.Decode(&item)
	if err != nil {
		log.Printf("ERROR: decoding menu item: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	conn := getPostgresConnection()

	//we might want to return the menuitemid here
	err = menu.CreateItem(&item, conn)
	if err != nil {
		log.Printf("ERROR: creating menu item - %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
