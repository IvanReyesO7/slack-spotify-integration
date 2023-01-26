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
	var text string = fmt.Sprintf("*Track Name*\n%s\n*Artist*\n%s\n*Album Name*\n%s\n", track.Title, track.Artist, track.Album)
	actions := []slack.AttachmentAction{slack.AttachmentAction{Name: "Add", Text: "Add to Playlist", Type: "button", Value: track.Id, Style: "primary"}}
	attachment := slack.Attachment{Color: "#1CDF63", Text: text, Actions: actions, CallbackID: track.Id, Fallback: "Done!"}
	return attachment
}

func UpdateOriginalMessage(channel_id string, responseUrl string, text string) {
	api := slack.New(os.Getenv("SLACK_TOKEN"), slack.OptionDebug(true))
	attachment := slack.Attachment{Color: "#1CDF63", Text: text}
	_, _, _, err := api.UpdateMessage(channel_id, "1405894322.002768", slack.MsgOptionResponseURL(responseUrl, "in_channel"), slack.MsgOptionReplaceOriginal(responseUrl), slack.MsgOptionAttachments(attachment))
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
