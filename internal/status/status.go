package status

import (
	"errors"
	"fmt"
	"github.com/brianstrauch/spotify"
	"github.com/spf13/cobra"
	"spotify/internal"
)

type EmojiLine struct {
	Emoji string
	Line  string
}

func (e EmojiLine) String() string {
	return fmt.Sprintf("   %s\r%s\n", e.Line, e.Emoji)
}

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "status",
		Aliases: []string{"s"},
		Short:   "show the current song or episode",
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

func status(api internal.APIInterface) (string, error) {
	playback, err := api.GetPlayback()
	if err != nil {
		return "", err
	}

	if playback == nil {
		return "", errors.New(internal.ErrNoActiveDevice)
	}

	return Show(playback), nil
}

func Show(playback *spotify.Playback) string {
	var artistLine string
	switch playback.Item.Type {
	case "track":
		artistLine = JoinArtists(playback.Item.Artists)
	case "episode":
		artistLine = playback.Item.Show.Name
	}

	var isPlayingEmoji string
	if playback.IsPlaying {
		isPlayingEmoji = "▶️"
	} else {
		isPlayingEmoji = "⏸"
	}

	progressBar := showProgressBar(playback.ProgressMs, playback.Item.Duration)

	output := EmojiLine{Emoji: "🎵", Line: playback.Item.Name}.String()
	output += EmojiLine{Emoji: "🎤", Line: artistLine}.String()
	output += EmojiLine{Emoji: isPlayingEmoji, Line: progressBar}.String()
	return output
}

func JoinArtists(artists []spotify.Artist) string {
	list := artists[0].Name
	for i := 1; i < len(artists); i++ {
		list += ", " + artists[i].Name
	}
	return list
}

func showProgressBar(progress int, duration *spotify.Duration) string {
	const length = 16
	bars := length * progress / int(duration.Milliseconds())

	status := fmt.Sprintf("%s [", formatTime(progress))
	for i := 0; i < bars; i++ {
		status += "="
	}
	for i := bars; i < length; i++ {
		status += " "
	}
	status += fmt.Sprintf("] %s", formatTime(int(duration.Milliseconds())))

	return status
}

func formatTime(ms int) string {
	s := ms / 1000
	m := s / 60
	h := m / 60

	if h == 0 {
		return fmt.Sprintf("%d:%02d", m, s%60)
	} else {
		return fmt.Sprintf("%d:%02d:%02d", h, m%60, s%60)
	}
}
