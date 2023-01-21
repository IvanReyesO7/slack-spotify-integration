package spotify

import (
	"context"
	"fmt"
	"log"
	"os"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/zmb3/spotify/v2"
)

type Song struct {
	Title    string
	Album    string
	Artist   string
	Duration int
}

func GetSongs(hint string) ([]Song, error) {
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

	results, err := client.Search(ctx, hint, spotify.SearchTypeTrack, spotify.Limit(5))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	tracks := []Song{}
	if results.Tracks != nil {
		fmt.Println("Tracks:")
		for _, item := range results.Tracks.Tracks {
			tracks = append(tracks, Song{Title: item.Name, Album: item.Album.Name, Artist: item.Artists[0].Name, Duration: item.Duration})
		}
	}
	return tracks, nil
}
