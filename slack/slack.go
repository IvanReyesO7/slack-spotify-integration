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

	for _, track := range tracks {
		json := formatMessage(track)
		fmt.Println(json)
		_, _, err := api.PostMessage(channel, slack.MsgOptionTS(thread_ts), slack.MsgOptionText(track.Title, true))
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

	}
}

func formatMessage(track Spotify.Song) string {
	return ""
}
