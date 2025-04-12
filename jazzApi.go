package main

import (
	"github.com/gin-gonic/gin"

	"net/http"
	"encoding/json"
	"os"
	"log"
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
	router.GET("/albums", getAlbums)
	router.POST("/create", createAlbum)

	router.Run("localhost:80")
}

func readFile() (err error) {
	var f *os.File

	f, err = os.Open("example.json")
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
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}
