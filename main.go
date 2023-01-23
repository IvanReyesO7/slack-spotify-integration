package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	Slack "slack-spotify-integration/slack"
	Spotify "slack-spotify-integration/spotify"

	"github.com/gin-gonic/gin"
)

type Event struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	Channel string `json:"channel"`
	Ts      string `json:"ts"`
	BotId   string `json:"bot_id"`
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
		var event Event = jsonRequest.Event

		if (event.Text != "list") && (event.BotId == "") && (event.Text != "") {
			tracks, err := Spotify.GetSongs(event.Text)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if len(tracks) < 1 {
				c.JSON(http.StatusNotFound, gin.H{"error": "No track matched keywork passed"})
				return
			}

			Slack.SendTracks(event.Channel, event.Ts, tracks)
			return
		}
		c.JSON(http.StatusOK, gin.H{"response": jsonRequest})
	})
	r.POST("/add-to-playlist", func(c *gin.Context) {
		request, _ := ioutil.ReadAll(c.Request.Body)
		encodedValue := string(request)
		decodedValue, err := url.QueryUnescape(encodedValue)

		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(decodedValue)

		c.JSON(http.StatusOK, gin.H{
			"replace_original": "true",
			"text":             "✅ Added to the playlist",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
