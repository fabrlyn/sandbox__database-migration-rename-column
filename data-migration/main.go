package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "migration-user"
	password = "password"
	dbname   = "migration-test"
)

type Beer struct {
	ID    uuid.UUID `json:"id"`
	Label string    `json:"label"`
}

type CreateBeer struct {
	Label string `json:"label"`
}

func connectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func getBeers() ([]Beer, error) {
	beers := []Beer{}
	db, err := sql.Open("postgres", connectionString())
	if err != nil {
		return beers, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, label FROM beer WHERE name = NULL OR brewery_name = NULL LIMIT 1")
	if err != nil {
		return beers, err
	}

	for rows.Next() {
		var beer Beer
		if err := rows.Scan(&beer.ID, &beer.Label); err != nil {
			return beers, err
		}
		beers = append(beers, beer)
	}

	return beers, nil
}

func createBeer(toCreate CreateBeer) (Beer, error) {
	beer := Beer{ID: uuid.NewV4(), Label: toCreate.Label}

	db, err := sql.Open("postgres", connectionString())
	if err != nil {
		return beer, err
	}
	defer db.Close()

	if _, err := db.Query("INSERT INTO beer (id, label) VALUES($1, $2)", beer.ID, beer.Label); err != nil {
		return beer, err
	}

	return beer, nil
}

func main() {
}
