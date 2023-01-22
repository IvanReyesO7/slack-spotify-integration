package slack

import (
	"fmt"
	"os"

	Spotify "slack-spotify-integration/spotify"

	"github.com/slack-go/slack"
)

func SendTracks(channel string, thread_ts string, tracks []Spotify.Song) {
	api := slack.New(os.Getenv("SLACK_TOKEN"), slack.OptionDebug(true))

	for _, track := range tracks {
		json := formatMessage(track)
		fmt.Println(json)
		_, _, err := api.PostMessage(channel, slack.MsgOptionTS(thread_ts), slack.MsgOptionText(json, true))
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

	}
}

func formatMessage(track Spotify.Song) string {
	json := fmt.Sprintf(`{
		"blocks": [
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": ">*Track Name*\n>%s\n>*Album Name*\n>%s\n>*Artist*\n>%s\n"
				}
			},
			{
				"type": "actions",
				"elements": [
					{
						"type": "button",
						"text": {
							"type": "plain_text",
							"emoji": true,
							"text": "Add to Playlist"
						},
						"style": "primary",
						"value": "%d"
					}
				]
			}
		]
	}`, track.Title, track.Album, track.Artist, track.Id)
	return json
}
