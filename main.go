package main

import (
	"os"

	"github.com/gin-contrib/location"

	"github.com/gin-gonic/gin"

	"github.com/subosito/gotenv"

	"github.com/TIG/api-sqlite/routers"
)

func init() {
	gotenv.Load()
}

func getPort() string {
	p := os.Getenv("HOST_PORT")
	if p != "" {
		return ":" + p
	}
	return ":3000"
}

// CORSMiddleware func
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, auth-session, auth-session-v9, content-type, auth-nonce, auth-nonce-response, userid, token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	port := getPort()
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.Default()
	r.Use(location.Default())
	r.Use(CORSMiddleware())
	rg := r.Group("/api/v1")
	rg.Use(CORSMiddleware())
	{
		routers.UserRoute(rg)
		routers.BatchRoute(rg)
	}
	r.Run(port)
}
