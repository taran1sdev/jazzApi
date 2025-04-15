package main

import (
	"github.com/gin-gonic/gin"
	
	"jazzApi/auth"

	"net/http"
	"encoding/json"
	"os"
	"log"
	"strconv"
	"fmt"
)

type album struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums []album

func main() {
	err := readFile()
	
	if err != nil{
		log.Fatal(err)
	}

	router := gin.Default()
	
	router.POST("/login", auth.HandleLogin)
	
	// User endpoints
	router.Use(auth.AuthMiddleware())
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", searchAlbum)
	
	// Admin endpoint
	// Add user - user/create
	// List users - user
	router.POST("/create", auth.AdminMiddleware(), createAlbum)

	router.Run("localhost:80")
}

func readFile() (err error) {
	var f *os.File

	f, err = os.Open("data/example.json")
	if err != nil {
		return 
	}

	fInfo, _ := f.Stat()
	b := make([]byte, fInfo.Size())

	_, err = f.Read(b)
	
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &albums)

	if err != nil {
		return
	}
	return nil
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func createAlbum(c *gin.Context) {
	fmt.Println("Reached Here")
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func searchAlbum(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)

	for _, a := range albums {
		if a.ID == idInt {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "album not found"})
}
