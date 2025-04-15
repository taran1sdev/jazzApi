package auth

import (
	"time"
	"crypto/md5"
	"errors"
	"net/http"
	"encoding/json"
	"encoding/hex"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)

// Create a data-type for the user

type User struct {
	ID 	 int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role 	 string	`json:"role"`
}

var users []User

var secretKey = []byte("secretkey") // Change this in prod

// Creates a new JWT Token valid for one hour

func createToken(username string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role": role,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Opens the user json file (will refactor to read from a database)

func getUsers() (err error) {
	var f *os.File

	f, err = os.Open("data/userExample.json")
	if err != nil {
		return errors.New("Error opening user file")
	}

	fInfo, _ := f.Stat()
	b := make([]byte, fInfo.Size())

	_, err = f.Read(b)

	if err != nil {
		return errors.New("Error reading user file")
	}

	err = json.Unmarshal(b, &users)

	if err != nil {
		return
	}
	return nil
}

// Function to return md5 hash 

func getHash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// Function handles login with username/password credentials - refactor to accept json requests also

func HandleLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	
	err := getUsers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading user data"})
		return
	}

	var role string

	h := getHash(password)
	
	for _, u := range users {
		if u.Username == username && u.Password == h {
			role = u.Role 
		}
	}
	if role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
		return
	}
	
	// Creates JWT

	token, err := createToken(username, role)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func verifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token)(interface{},error){
		return secretKey, nil
	})

	if !token.Valid {
		return nil, errors.New("Invalid JWT Token")
	}
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
	
		return claims, nil
	}

	return nil, errors.New("Error Mapping claims")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		
		// Checks if Header is missing
		if len(tokenString) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Verify JWT
		claims, valid := verifyToken(tokenString)

		if valid != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": valid.Error()})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("claims").(jwt.MapClaims)
		role := claims["role"].(string)

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
