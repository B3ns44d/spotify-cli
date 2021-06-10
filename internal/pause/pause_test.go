package pause

import (
	"errors"
	"spotify/internal"
	"testing"

	"github.com/brianstrauch/spotify"
	"github.com/stretchr/testify/require"
)

func TestPauseCommand(t *testing.T) {
	api := new(spotify.MockAPI)

	playback1 := &spotify.Playback{
		IsPlaying:  true,
		ProgressMs: 0,
		Item: spotify.Item{
			Type: "track",
			Name: "Song",
			Artists: []spotify.Artist{
				{Name: "Artist"},
			},
			DurationMs: 1000,
		},
	}

	playback2 := new(spotify.Playback)
	*playback2 = *playback1
	playback2.IsPlaying = false

	api.On("GetPlayback").Return(playback1, nil).Once()
	api.On("GetPlayback").Return(playback2, nil).Once()
	api.On("Pause").Return(nil)

	status, err := Pause(api)
	require.NoError(t, err)
	require.Equal(t, "   Song\r🎵\n   Artist\r🎤\n   0:00 [                ] 0:01\r⏸\n", status)
}

func TestAlreadyPaused(t *testing.T) {
	api := new(spotify.MockAPI)
	api.On("GetPlayback").Return(new(spotify.Playback), nil)
	api.On("Pause").Return(errors.New(internal.RestrictionViolatedSpotifyErr))

	_, err := Pause(api)
	require.Equal(t, internal.AlreadyPausedErr, err.Error())
}

func TestNoActiveDeviceErr(t *testing.T) {
	api := new(spotify.MockAPI)
	api.On("GetPlayback").Return(nil, nil)

	_, err := Pause(api)
	require.Equal(t, internal.NoActiveDeviceErr, err.Error())
}
