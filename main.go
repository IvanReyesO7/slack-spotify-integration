package main

import (
	"net/http"

	Spotify "spotify-in-the-office/spotify"

	"github.com/gin-gonic/gin"
)

type Event struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type JsonRequest struct {
	Event Event `json:"event"`
}

func main() {
	r := gin.Default()
	r.POST("/test", func(c *gin.Context) {
		var jsonRequest JsonRequest
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if jsonRequest.Event.Text != "list" {
			Spotify.GetSongs(jsonRequest.Event.Text)
		}

		c.JSON(http.StatusOK, gin.H{"response": jsonRequest})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
