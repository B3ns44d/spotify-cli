package status

import (
	"errors"
	"fmt"
	"spotify/internal"

	"github.com/brianstrauch/spotify"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Aliases: []string{"s"},
		Short:   "Show the current song or episode.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			api, err := internal.Authenticate()
			if err != nil {
				return err
			}

			status, err := status(api)
			if err != nil {
				return err
			}

			cmd.Print(status)
			return nil
		},
	}
}

func status(api spotify.APIInterface) (string, error) {
	playback, err := api.GetPlayback()
	if err != nil {
		return "", err
	}

	if playback == nil {
		return "", errors.New(internal.NoActiveDeviceErr)
	}

	return Show(playback), nil
}

func Show(playback *spotify.Playback) string {
	var artistLine string
	switch playback.Item.Type {
	case "track":
		artistLine = joinArtists(playback.Item.Artists)
	case "episode":
		artistLine = playback.Item.Show.Name
	}

	var isPlayingEmoji string
	if playback.IsPlaying {
		isPlayingEmoji = "▶️"
	} else {
		isPlayingEmoji = "⏸"
	}

	progressBar := showProgressBar(playback.ProgressMs, playback.Item.DurationMs)

	status := prefixLineWithEmoji("🎵", playback.Item.Name)
	status += prefixLineWithEmoji("🎤", artistLine)
	status += prefixLineWithEmoji(isPlayingEmoji, progressBar)

	return status
}

func joinArtists(artists []spotify.Artist) string {
	list := artists[0].Name
	for i := 1; i < len(artists); i++ {
		list += ", " + artists[i].Name
	}
	return list
}

func showProgressBar(progress, duration int) string {
	const length = 16
	bars := length * progress / duration

	status := fmt.Sprintf("%s [", formatTime(progress))
	for i := 0; i < bars; i++ {
		status += "="
	}
	for i := bars; i < length; i++ {
		status += " "
	}
	status += fmt.Sprintf("] %s", formatTime(duration))

	return status
}

func formatTime(ms int) string {
	s := ms / 1000
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

func prefixLineWithEmoji(emoji, line string) string {
	// Carriage return jumps to start of line because emojis can have variable widths
	return fmt.Sprintf("   %s\r%s\n", line, emoji)
}