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
	Infrastructure.NewConfig()
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"response": "Server up"})
	})
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

		if event.BotId == "" && event.Text != "" {
			switch strings.ToLower(event.Text) {
			case "list":
				tracks, _ := Spotify.GetPlaylistQueue()
				if tracks == nil {
					c.JSON(418, gin.H{"error": "No tracks in the playlist"})
					return
				}

				Slack.SendTracks(event.Channel, event.Ts, tracks, false)
				return
			case "commands", "help":
				err := Slack.SendCommands(event.Channel, event.Ts)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				} else {
					c.JSON(http.StatusOK, nil)
					return
				}
			default:
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
		action := (gjson.Get(trackValue, "action")).String()
		if action == "add" {
			snapshot, err := Spotify.AddTrackToPlaylist(trackId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"replace_original": "true",
					"text":             "⛔️ Sorry, Something went wrong",
				})
			} else if snapshot != nil {
				Slack.UpdateOriginalMessage(trackValue, channelId, messageTs, responseUrl, action)
			}
		} else if action == "remove" {
			snapshot, err := Spotify.RemoveFromPlaylist(trackId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"replace_original": "true",
					"text":             "⛔️ Sorry, Something went wrong",
				})
			} else if snapshot != nil {
				Slack.UpdateOriginalMessage(trackValue, channelId, messageTs, responseUrl, action)
			}
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
