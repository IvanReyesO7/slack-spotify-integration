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
		attachments := buildAttachmentTwo(track)

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

func buildAttachmentTwo(track Spotify.Song) slack.Attachment {
	header := buildHeader(track)
	footer := buildFooter(track)
	action := slack.ActionBlock{Type: "actions", Elements: &slack.BlockElements{ElementSet: []slack.BlockElement{slack.ButtonBlockElement{Type: "button", Text: &slack.TextBlockObject{Type: "plain_text", Text: "Add To Playlist"}, Value: track.Id, Style: "primary"}}}}
	blocks := []slack.Block{header, footer, action}
	attachment := slack.Attachment{Color: "#1CDF63", Blocks: slack.Blocks{BlockSet: blocks}}
	return attachment
}

func buildHeader(track Spotify.Song) slack.SectionBlock {
	text := fmt.Sprintf("*Track Name*\n%s\n*Album Name*\n%s", track.Title, track.Album)
	accessory := slack.Accessory{ImageElement: &slack.ImageBlockElement{Type: "image", ImageURL: track.UrlImage}}
	return slack.SectionBlock{Type: "section", Text: &slack.TextBlockObject{Type: "mrkdwn", Text: text}, Accessory: &accessory}
}

func buildFooter(track Spotify.Song) slack.SectionBlock {
	text := fmt.Sprintf("*Artist*\n%s\n*ID*\n`%s`", track.Artist, track.Id)
	return slack.SectionBlock{Type: "section", Text: &slack.TextBlockObject{Type: "mrkdwn", Text: text}}
}

func UpdateOriginalMessage(channelId string, messageTs string, responseUrl string) {
	api := slack.New(os.Getenv("SLACK_TOKEN"), slack.OptionDebug(true))
	_, _, _, err := api.UpdateMessage(channelId, messageTs, slack.MsgOptionReplaceOriginal(responseUrl), slack.MsgOptionText("Track added to the Playlist", false))
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
