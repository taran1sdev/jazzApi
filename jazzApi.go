package main

import (
	"github.com/gin-gonic/gin"
	
	"jazzApi/auth"

	"net/http"
	"encoding/json"
	"os"
	"log"
	"strconv"
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
	router.POST("albums/create", auth.AdminMiddleware(), createAlbum)

	router.GET("/users", auth.AdminMiddleware(), getUsers)
	router.POST("/users/create", auth.AdminMiddleware(), createUser)

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

func getUsers(c *gin.Context) {
	if auth.Users == nil {
		err := auth.GetUsers()
	
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.IndentedJSON(http.StatusOK, auth.Users)
}

func createUser(c *gin.Context) {
	var newUser auth.User
	

	
	if c.GetHeader("Content-Type") == "application/json" {
		if err := c.BindJSON(&newUser); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Unable to create new user"})
			return
		}
	} else {
		if c.PostForm("username") != ""  && c.PostForm("password") != "" && c.PostForm("role") != "" {	
			newUser.Username = c.PostForm("username")
			newUser.Password = auth.GetHash(c.PostForm("password"))
			newUser.Role = c.PostForm("role")
		} else {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error":"Missing Necessary Fields"})
			return
		}
	}
	

	newUser.ID = auth.Users[len(auth.Users)-1].ID + 1
	auth.Users = append(auth.Users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}
