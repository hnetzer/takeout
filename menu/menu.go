package menu

import (
	"database/sql"
	"log"
)

//Item is something you would see on menu
type Item struct {
	MenuItemID   string
	RestaurantID string
	Name         string
	Description  *string
	ImageURL     *string
	Price        *int64
}

//CreateItem creates a menu item for a given restaurant
func CreateItem(item *Item, pgconn *sql.DB) error {

	log.Printf("INFO: creating menu item %+v", item)

	_, err := pgconn.Exec("INSERT INTO menuitems (name, description, price, restaurantid) VALUES ($1, $2, $3, $4)",
		item.Name, item.Description, item.Price, item.RestaurantID)
	if err != nil {
		return err
	}

	return nil
}

//ReadAllItems reads all menu items for a given restaurantID
func ReadAllItems(restaurantID string, pgconn *sql.DB) ([]Item, error) {

	log.Printf("INFO: getting menu items for restaurantID %s", restaurantID)

	rows, err := pgconn.Query("SELECT menuitemid, name, description, imageurl, price FROM menuitems WHERE restaurantid = $1", restaurantID)
	if err != nil {
		return nil, err
	}

	menuitems := []Item{}
	for rows.Next() {
		var item Item
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
