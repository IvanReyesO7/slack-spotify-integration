package spotify

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/tidwall/gjson"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/zmb3/spotify/v2"
)

type Song struct {
	Id       string
	Title    string
	Album    string
	Artist   string
	Duration int
}

func GetSongs(keyword string) ([]Song, error) {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
		return nil, err
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	results, err := client.Search(ctx, keyword, spotify.SearchTypeTrack, spotify.Limit(5))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	tracks := []Song{}
	if results.Tracks != nil {
		fmt.Println("Tracks:")
		for _, item := range results.Tracks.Tracks {
			tracks = append(tracks, Song{Title: item.Name, Album: item.Album.Name, Artist: item.Artists[0].Name, Duration: item.Duration, Id: string(item.ID)})
		}
	}
	return tracks, nil
}

func AddTrackToPlaylist(track_id string) (*string, error) {
	ctx := context.Background()
	token := oauth2.Token{AccessToken: RefreshSpotifyAccessToken()}

	httpClient := spotifyauth.New().Client(ctx, &token)
	client := spotify.New(httpClient)
	playlist_id := spotify.ID(os.Getenv("SPOTIFY_PLAYLIST_ID"))
	track := spotify.ID(track_id)

	snapshot, err := client.AddTracksToPlaylist(ctx, playlist_id, track)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &snapshot, nil
}

func RefreshSpotifyAccessToken() string {
	requestUrl := "https://accounts.spotify.com/api/token"
	buffer := fmt.Sprintf("Basic %s", os.Getenv("SPOTIFY_BUFFER"))

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", os.Getenv("SPOTIFY_REFRESH_TOKEN"))

	req, err := http.NewRequest(http.MethodPost, requestUrl, strings.NewReader(data.Encode()))
	req.Header.Set("Authorization", buffer)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	return gjson.Get(bodyString, "access_token").String()
}
