package slack

import (
	"fmt"
	"os"

	Spotify "slack-spotify-integration/spotify"

	"github.com/slack-go/slack"
)

func SendTracks(channel string, thread_ts string, tracks []Spotify.Song) {
	fmt.Println(tracks)
	api := slack.New(os.Getenv("SLACK_TOKEN"), slack.OptionDebug(true))
	_, _, err := api.PostMessage(channel, slack.MsgOptionTS(thread_ts), slack.MsgOptionText("Got it!", true))

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

func formatMessage(tracks []Spotify.Song) {

}
