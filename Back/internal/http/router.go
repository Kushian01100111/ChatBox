package http

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewHandler() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Content-type", "Accept", "Authorization", "Origin"},
		ExposeHeaders:    []string{"Content-length"},
		AllowCredentials: true,
		MaxAge:           25 * time.Minute,
	}))

	return r
}
