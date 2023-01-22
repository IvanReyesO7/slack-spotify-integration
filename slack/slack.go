package slack

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

// func SendTracks(channel string, tracks []Spotify.Song) {

// }

func Ping(channel string, thread_ts string) {
	// api := slack.New("")
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	api := slack.New(os.Getenv("SLACK_TOKEN"), slack.OptionDebug(true))
	_, _, err := api.PostMessage(channel, slack.MsgOptionTS(thread_ts), slack.MsgOptionText("Got it!", true))

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}
