package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"

	Infrastructure "slack-spotify-integration/infrastructure"
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
	gin.SetMode(gin.ReleaseMode)
	Infrastructure.NewConfig()
	r := gin.Default()
	r.POST("/endpoint", func(c *gin.Context) {
		challenge, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Printf("%s", string(challenge))
		c.JSON(http.StatusOK, gin.H{
			"challenge": string(challenge),
		})
	})
	r.POST("/", func(c *gin.Context) {
		var jsonRequest JsonRequest
		if err := c.ShouldBindJSON(&jsonRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var event Event = jsonRequest.Event

		if event.BotId == "" && event.Text != "" {
			if strings.ToLower(event.Text) != "list" {
				tracks, err := Spotify.GetSongs(event.Text)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				if len(tracks) < 1 {
					c.JSON(http.StatusNotFound, gin.H{"error": "No track matched keywork passed"})
					return
				}

				Slack.SendTracks(event.Channel, event.Ts, tracks, true)
				return
			}
			if strings.ToLower(event.Text) == "list" {
				tracks, _ := Spotify.GetPlaylistQueue()
				if tracks == nil {
					c.JSON(418, gin.H{"error": "No tracks in the playlist"})
					return
				}

				Slack.SendTracks(event.Channel, event.Ts, tracks, false)
				return
			}
		}
		c.JSON(http.StatusOK, nil)
		return
	})
	r.POST("/add-to-playlist", func(c *gin.Context) {
		request, _ := ioutil.ReadAll(c.Request.Body)
		encodedValue := string(request)
		decodedValue, err := url.QueryUnescape(encodedValue)

		if err != nil {
			return
		}

		json := decodedValue[8:]

		trackValue := (gjson.Get(json, "actions.0.value")).String()
		trackId := (gjson.Get(trackValue, "id")).String()
		responseUrl := (gjson.Get(json, "response_url")).String()
		channelId := (gjson.Get(json, "channel.id")).String()
		messageTs := (gjson.Get(json, "container.message_ts")).String()
		snapshot, err := Spotify.AddTrackToPlaylist(trackId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"replace_original": "true",
				"text":             "⛔️ Sorry, Something went wrong",
			})
		} else if snapshot != nil {
			Slack.UpdateOriginalMessage(trackValue, channelId, messageTs, responseUrl)

		}

	})
	r.GET("/callback", func(c *gin.Context) {
		code := c.Request.URL.Query().Get("code")
		fmt.Printf("%s", string(code))
		c.JSON(http.StatusOK, gin.H{
			"code": string(code),
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
