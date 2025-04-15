package auth

import (
	"time"
	"crypto/md5"
	"fmt"
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
		return fmt.Errorf("Error opening user file")
	}

	fInfo, _ := f.Stat()
	b := make([]byte, fInfo.Size())

	_, err = f.Read(b)

	if err != nil {
		return fmt.Errorf("Error reading user file")
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

	if err != nil{
		fmt.Println(err)
	}

	var role string

	h := getHash(password)
	fmt.Println(h)
	
	for _, u := range users {
		fmt.Println(u.Password)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func verifyToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token)(interface{},error){
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid JWT Token")
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("Error Mapping claims")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		
		// Verify JWT
		claims, valid := verifyToken(tokenString)

		if valid != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": valid})
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


func HandleGeneral(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "General Resource Accessed",
	})
}

func HandleAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin Resource Accessed",
	})
}
