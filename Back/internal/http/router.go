package http

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kushian01100111/ChatBox/internal/app/chat"
)

func NewHandler(hub *chat.Hub) http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Content-type", "Accept", "Authorization", "Origin"},
		ExposeHeaders:    []string{"Content-length"},
		AllowCredentials: true,
		MaxAge:           25 * time.Minute,
	}))

	context := r.Group("/ws")
	{
		context.GET("/:roomId", func(g *gin.Context) {
			roomId := g.Param("roomId")
			chat.ServeWS(g, roomId, hub)
		})
	}
	return r
}
