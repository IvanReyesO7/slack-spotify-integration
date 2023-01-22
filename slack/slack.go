package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

// func SendTracks(channel string, tracks []Spotify.Song) {

// }

func Ping(channel string, thread_ts string) {
	// api := slack.New("xoxb-2882488176691-4678049065891-d9eT4MiNpUPpoUCLo41YNNs7")
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	api := slack.New("xoxb-2882488176691-4678049065891-d9eT4MiNpUPpoUCLo41YNNs7", slack.OptionDebug(true))
	_, _, err := api.PostMessage(channel, slack.MsgOptionTS(thread_ts), slack.MsgOptionText("Got it!", true))

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	return
}
