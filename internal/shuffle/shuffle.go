package shuffle

import (
	"errors"
	"spotify/internal"

	"github.com/brianstrauch/spotify"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "shuffle",
		Short: "Turn shuffle on or off.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			api, err := internal.Authenticate()
			if err != nil {
				return err
			}

			var state bool

			switch args[0] {
			case "on":
				state = true
			case "off":
				state = false
			default:
				return errors.New(internal.ErrShuffleArg)
			}

			if err := Shuffle(api, state); err != nil {
				return err
			}

			if state {
				cmd.Println("🔀 Shuffle on")
			} else {
				cmd.Println("🔀 Shuffle off")
			}

			return nil
		},
	}
}

func Shuffle(api internal.APIInterface, state bool) error {
	if err := api.Shuffle(state); err != nil {
		return err
	}

	_, err := internal.WaitForUpdatedPlayback(api, func(playback *spotify.Playback) bool {
		return playback.ShuffleState == state
	})
	return err
}
