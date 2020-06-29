package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
	ID   uuid.UUID `json:"id"`
	Name string    `json:"title"`
}

type CreateBeer struct {
	Title string `json:"title"`
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
	rows, err := db.Query("SELECT id, name FROM beer")
	if err != nil {
		return beers, err
	}

	for rows.Next() {
		var beer Beer
		if err := rows.Scan(&beer.ID, &beer.Name); err != nil {
			return beers, err
		}
		beers = append(beers, beer)
	}

	return beers, nil
}

func createBeer(toCreate CreateBeer) (Beer, error) {
	beer := Beer{ID: uuid.NewV4(), Name: toCreate.Title}

	db, err := sql.Open("postgres", connectionString())
	if err != nil {
		return beer, err
	}
	defer db.Close()

	if _, err := db.Query("INSERT INTO beer (id, name) VALUES($1, $2)", beer.ID, beer.Name); err != nil {
		return beer, err
	}

	return beer, nil
}

func main() {
	r := gin.Default()

	r.GET("beer", func(c *gin.Context) {
		beers, err := getBeers()
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, beers)
	})

	r.POST("beer", func(c *gin.Context) {
		var toCreate CreateBeer
		if err := c.Bind(&toCreate); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		beer, err := createBeer(toCreate)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, beer)
	})

	r.Run(":5003")
}
