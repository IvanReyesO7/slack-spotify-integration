package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	Spotify "slack-spotify-integration/spotify"

	"github.com/gin-gonic/gin"
)

type Event struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	Channel string `json:"channel"`
}

type JsonRequest struct {
	Event Event `json:"event"`
}

func main() {
	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		challenge, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Printf("%s", string(challenge))
		c.JSON(http.StatusOK, gin.H{
			"challenge": string(challenge),
		})
	})
	r.POST("/endpoint", func(c *gin.Context) {
		var jsonRequest JsonRequest
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if jsonRequest.Event.Text != "list" {
			tracks, err := Spotify.GetSongs(jsonRequest.Event.Text)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if len(tracks) < 1 {
				c.JSON(http.StatusNotFound, gin.H{"error": "No track matched keywork passed"})
				return
			}
			fmt.Println(tracks)
			// slack.SendTracks(channel, tracks)
		}

		c.JSON(http.StatusOK, gin.H{"response": jsonRequest})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
