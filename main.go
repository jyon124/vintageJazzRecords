package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

// album represents data about a record album.
type album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var db *sqlx.DB

func main() {
	var err error
	// Retrieve the DATABASE_URL from environment variables
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Print("Using Local PostgreSQL DB...")
		db, err = sqlx.Connect(
			"postgres",
			"user=adminjazz dbname=vintagejazzrecord sslmode=disable password=mys3cret",
		)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Connect to the database using the DATABASE_URL
		db, err = sqlx.Connect("postgres", databaseURL)
		if err != nil {
			log.Fatal(err)
		}
	}

	createTable()

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// createTable creates the albums table if it does not exist
func createTable() {
	// Drop the table if it exists
	// dropTableQuery := `DROP TABLE IF EXISTS albums;`
	// _, err := db.Exec(dropTableQuery)
	// if err != nil {
	// 	log.Fatalf("Failed to drop table: %v", err)
	// }

	createTableQuery := `CREATE TABLE IF NOT EXISTS albums (
		id SERIAL PRIMARY KEY,
		title varchar(50) NOT NULL,
		artist varchar(50) NOT NULL,
		price NUMERIC(10,2) NOT NULL
	);`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	} else {
		log.Println("Tables are ready...")
	}
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	var albums []album
	err := db.Select(&albums, "SELECT * FROM albums")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not retrieve albums"})
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	var a album
	err := db.Get(&a, "SELECT * FROM albums WHERE id=$1", id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, a)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	query := `INSERT INTO albums (title, artist, price) VALUES ($1, $2, $3) RETURNING id`
	err := db.Get(&newAlbum.ID, query, newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		// Print the actual error
		log.Printf("Error inserting album: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not add album"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newAlbum)
}
