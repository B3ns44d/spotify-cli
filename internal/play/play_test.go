package play

import (
	"errors"
	"spotify/internal"
	"testing"

	"github.com/brianstrauch/spotify"
	"github.com/stretchr/testify/require"
)

func TestPlay(t *testing.T) {
	api := new(internal.MockAPI)

	playback1 := &spotify.Playback{
		IsPlaying:  false,
		ProgressMs: 0,
		Item: spotify.Item{
			ID:   "0",
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
	playback2.IsPlaying = true

	api.On("GetPlayback").Return(playback1, nil).Once()
	api.On("GetPlayback").Return(playback2, nil).Once()
	api.On("Play", []string(nil)).Return(nil)

	status, err := Play(api, "")
	require.NoError(t, err)
	require.Equal(t, "   Song\r🎵\n   Artist\r🎤\n   0:00 [                ] 0:01\r▶️\n", status)
}

func TestPlay_ErrAlreadyPlaying(t *testing.T) {
	api := new(internal.MockAPI)
	api.On("GetPlayback").Return(new(spotify.Playback), nil)
	api.On("Play", []string(nil)).Return(errors.New(internal.ErrRestrictionViolated))

	_, err := Play(api, "")
	require.Error(t, err)
	require.Equal(t, internal.ErrAlreadyPlaying, err.Error())
}

func TestPlay_ErrNoActiveDevice(t *testing.T) {
	api := new(internal.MockAPI)
	api.On("GetPlayback").Return(nil, nil)

	_, err := Play(api, "")
	require.Error(t, err)
	require.Equal(t, internal.ErrNoActiveDevice, err.Error())
}
