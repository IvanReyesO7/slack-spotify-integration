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
		sectionBlock := buildSection(track)
		actionBlock := buildAction(track)

		_, _, err := api.PostMessage(channel, slack.MsgOptionTS(thread_ts), slack.MsgOptionBlocks(slack.SectionBlock(sectionBlock), slack.ActionBlock(actionBlock)))
		if err != nil {
			fmt.Printf("%s\n", err)
		}

	}
}

func buildSection(track Spotify.Song) slack.SectionBlock {
	var textBlockObject slack.TextBlockObject = slack.TextBlockObject{Type: "mrkdwn", Text: fmt.Sprintf("&gt;*Track Name*\n&gt;%s\n&gt;*Album Name*\n&gt;%s\n&gt;*Artist*\n&gt;%s\n", track.Title, track.Album, track.Artist)}
	var messageBlockType slack.MessageBlockType = "section"
	return slack.SectionBlock{Type: messageBlockType, Text: &textBlockObject}
}

func buildAction(track Spotify.Song) slack.ActionBlock {
	var messageBlockType slack.MessageBlockType = "actions"
	var TextBlockObject slack.TextBlockObject = slack.TextBlockObject{Type: "plain_text", Text: "Add to Playlist"}
	var blockElements slack.BlockElements = slack.BlockElements{ElementSet: []slack.BlockElement{slack.ButtonBlockElement{Type: "button", Text: &TextBlockObject, Value: track.Id}}}
	return slack.ActionBlock{Type: messageBlockType, Elements: &blockElements}
}
