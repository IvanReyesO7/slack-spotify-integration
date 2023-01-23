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
		attachments := buildAttachment(track)

		_, _, err := api.PostMessage(channel, slack.MsgOptionTS(thread_ts), slack.MsgOptionAttachments(attachments))
		if err != nil {
			fmt.Printf("%s\n", err)
		}

	}
}

func buildAttachment(track Spotify.Song) slack.Attachment {
	var text string = fmt.Sprintf("*Track Name*\n%s\n*Album Name*\n%s\n*Artist*\n%s\n", track.Title, track.Album, track.Artist)
	actions := []slack.AttachmentAction{slack.AttachmentAction{Name: "Add", Text: "Add to Playlist", Type: "button", Value: track.Id, Style: "primary"}}
	attachment := slack.Attachment{Color: "#1CDF63", Text: text, Actions: actions, CallbackID: track.Id, Fallback: "Done!"}
	return attachment
}
