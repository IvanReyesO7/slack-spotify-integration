package slack

import (
	"fmt"
	"os"

	Spotify "slack-spotify-integration/spotify"

	"github.com/slack-go/slack"
	"github.com/tidwall/gjson"
)

func SendTracks(channel string, thread_ts string, tracks []Spotify.Song, action bool) {
	api := slack.New(os.Getenv("SLACK_TOKEN"), slack.OptionDebug(true))

	for _, track := range tracks {
		var attachments slack.Attachment
		if action == true {
			attachments = buildAttachment(track)
		} else {
			attachments = buildAttachmentNoAction(track)
		}

		_, _, err := api.PostMessage(channel, slack.MsgOptionTS(thread_ts), slack.MsgOptionAttachments(attachments))
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}
	return
}

func buildAttachment(track Spotify.Song) slack.Attachment {
	header := buildHeader(track)
	footer := buildFooter(track)
	action := slack.ActionBlock{Type: "actions", Elements: &slack.BlockElements{ElementSet: []slack.BlockElement{slack.ButtonBlockElement{Type: "button", Text: &slack.TextBlockObject{Type: "plain_text", Text: "Add To Playlist"}, Value: fmt.Sprintf(`{"id": "%s", "trackName": "%s", "trackArtist": "%s", "trasckAlbum": "%s", "imageUrl": "%s"}`, track.Id, track.Title, track.Artist, track.Album, track.UrlImage), Style: "primary"}}}}
	blocks := []slack.Block{header, footer, action}
	attachment := slack.Attachment{Color: "#1CDF63", Blocks: slack.Blocks{BlockSet: blocks}}
	return attachment
}

func buildAttachmentNoAction(track Spotify.Song) slack.Attachment {
	header := buildHeader(track)
	return slack.Attachment{Color: "#1CDF63", Blocks: slack.Blocks{BlockSet: []slack.Block{header}}}
}

func buildHeader(track Spotify.Song) slack.SectionBlock {
	text := fmt.Sprintf("*Track Name*\n%s\n*Artist*\n%s", track.Title, track.Artist)
	accessory := slack.Accessory{ImageElement: &slack.ImageBlockElement{Type: "image", ImageURL: track.UrlImage}}
	return slack.SectionBlock{Type: "section", Text: &slack.TextBlockObject{Type: "mrkdwn", Text: text}, Accessory: &accessory}
}

func buildFooter(track Spotify.Song) slack.SectionBlock {
	text := fmt.Sprintf("*Album Name*\n%s\n*ID*\n`%s`", track.Album, track.Id)
	return slack.SectionBlock{Type: "section", Text: &slack.TextBlockObject{Type: "mrkdwn", Text: text}}
}

func UpdateOriginalMessage(trackValue string, channelId string, messageTs string, responseUrl string) {
	api := slack.New(os.Getenv("SLACK_TOKEN"), slack.OptionDebug(true))

	Id, Title, Artist, Album, UrlImage := (gjson.Get(trackValue, "id")).String(), (gjson.Get(trackValue, "trackName")).String(), (gjson.Get(trackValue, "trackArtist")).String(), (gjson.Get(trackValue, "trackAlbum")).String(), (gjson.Get(trackValue, "imageUrl")).String()
	track := Spotify.Song{Id: Id, Title: Title, Artist: Artist, Album: Album, UrlImage: UrlImage}
	header := buildHeader(track)
	footer := buildFooter(track)

	image := slack.ImageBlockElement{
		Type:     "image",
		ImageURL: track.UrlImage,
	}

	text := slack.TextBlockObject{Type: "mrkdwn", Text: fmt.Sprintf("_*%s*_ by _*%s*_ added to the Playlist!", track.Title, track.Artist)}
	context := slack.NewContextBlock("", image, text)

	blocks := []slack.Block{header, footer, context}
	attachment := slack.Attachment{Color: "#1CDF63", Blocks: slack.Blocks{BlockSet: blocks}}

	_, _, _, err := api.UpdateMessage(channelId, messageTs, slack.MsgOptionReplaceOriginal(responseUrl), slack.MsgOptionAttachments(attachment))
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

func SendCommands(channel string, thread_ts string) error {
	api := slack.New(os.Getenv("SLACK_TOKEN"), slack.OptionDebug(true))

	attachments := slack.Attachment{Color: "#1CDF63", Text: "These are the current commands you can use to interact with Spotify:\n`Commands||Help` Lists all the available commands.\n`List` Lists all the current tracks in the playlist.\n`{Song name}` Searches for a song that matches the keyword typed."}
	_, _, err := api.PostMessage(channel, slack.MsgOptionTS(thread_ts), slack.MsgOptionAttachments(attachments))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
