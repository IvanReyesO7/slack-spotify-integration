# slack-spotify-integration
A working Slack bot made with Gin-Gonic framework that integrates a Slack DM channel with a Spotify Account.
Developed to play some music at the office and democratise what everyone is listening to.
The bot is hosted in an instance of render.com, and is always listening to message events supported by the Slack API.

The bot currently supports a few, but important commands.

## `help` or `Commands`

List the available commands

<img width="564" alt="image" src="https://user-images.githubusercontent.com/74359151/220633262-5598ca74-86a4-4bbe-b8d6-a6ceef86d95f.png">


## `{Song name}`
Searches for a song that matches the keyword typed and returns the 5 most relevant results.

<img width="564" alt="image" src="https://user-images.githubusercontent.com/74359151/220633617-bb6ed452-f075-4ba8-9b5f-c28840bb4cbf.png">

If you click on `Add to Playlist`... well, the song gets added to the playlist, Du'h.

## `List`
Lists all the current tracks in the playlist.

<img width="564" alt="image" src="https://user-images.githubusercontent.com/74359151/220634232-08e9bd0c-69fa-433e-92e1-d6aabd613d04.png">

If the `Remove` button is clicked, the track is excluded from the playlist.
