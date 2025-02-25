package main

import (
	"log"
	"os"
	"strings"

	"github.com/electerm/electerm-sync-server-golang/src/store"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"status": "error", "message": "No authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"status": "error", "message": "Invalid token"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userId := claims["id"].(string)

		users := strings.Split(os.Getenv("JWT_USERS"), ",")
		authorized := false
		for _, user := range users {
			if user == userId {
				authorized = true
				break
			}
		}

		if !authorized {
			c.JSON(401, gin.H{"status": "error", "message": "Unauthorized!"})
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}

func handleSync(c *gin.Context) {
	userId := c.GetString("userId")

	if c.Request.Method == "GET" {
		data, err := store.FileStore.Read(userId)
		if err != nil {
			c.JSON(404, gin.H{"status": "error", "message": "File not found"})
			return
		}
		c.JSON(200, data)
		return
	}

	if c.Request.Method == "PUT" {
		var data map[string]interface{}
		if err := c.BindJSON(&data); err != nil {
			c.JSON(400, gin.H{"status": "error", "message": "Invalid JSON"})
			return
		}

		err := store.FileStore.Write(userId, data)
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "message": "Failed to write data"})
			return
		}

		c.String(200, "ok")
		return
	}

	c.JSON(405, gin.H{"status": "error", "message": "Method not allowed"})
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	authorized := r.Group("/api")
	authorized.Use(authMiddleware())
	{
		authorized.GET("/sync", handleSync)
		authorized.PUT("/sync", handleSync)
	}

	return r
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := setupRouter()

	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	log.Printf("Server running at http://%s:%s\n", host, port)
	r.Run(host + ":" + port)
}
