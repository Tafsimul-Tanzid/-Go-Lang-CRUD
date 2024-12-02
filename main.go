package main

import (
    "net/http"
    "os"
    "log"
    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
    ID     string  `json:"id" gorm:"primaryKey"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

var db *gorm.DB
var err error

func initDB() {
    // Load environment variables from .env file
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // Get the database connection URL from environment variables
    dsn := os.Getenv("DATABASE_URL")
    
    // Initialize the database connection
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Migrate the schema
    db.AutoMigrate(&album{})
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
    var albums []album
    db.Find(&albums)
    c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
    var newAlbum album

    // Call BindJSON to bind the received JSON to newAlbum.
    if err := c.BindJSON(&newAlbum); err != nil {
        return
    }

    // Add the new album to the database.
    db.Create(&newAlbum)
    c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
    id := c.Param("id")
    var album album

    // Find the album by ID.
    if err := db.First(&album, "id = ?", id).Error; err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
        return
    }
    c.IndentedJSON(http.StatusOK, album)
}

func main() {
    initDB()

    router := gin.Default()
    router.GET("/albums", getAlbums)
    router.GET("/albums/:id", getAlbumByID)
    router.POST("/albums", postAlbums)

    router.Run("localhost:8080")
}