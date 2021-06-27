package back

import (
	"errors"
	"spotify/internal"
	"testing"

	"github.com/brianstrauch/spotify"
	"github.com/stretchr/testify/require"
)

func TestBack(t *testing.T) {
	api := new(internal.MockAPI)

	playback1 := &spotify.Playback{
		IsPlaying:  true,
		ProgressMs: 0,
		Item: spotify.Item{
			ID:   "1",
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
	playback2.Item.ID = "0"

	api.On("GetPlayback").Return(playback1, nil).Once()
	api.On("GetPlayback").Return(playback2, nil).Once()
	api.On("SkipToPreviousTrack").Return(nil)

	status, err := back(api)
	require.NoError(t, err)
	require.Equal(t, "   Song\r🎵\n   Artist\r🎤\n   0:00 [                ] 0:01\r▶️\n", status)
}

func TestBack_ErrNoPrevious(t *testing.T) {
	api := new(internal.MockAPI)
	api.On("GetPlayback").Return(new(spotify.Playback), nil)
	api.On("SkipToPreviousTrack").Return(errors.New(internal.ErrRestrictionViolated))

	_, err := back(api)
	require.Error(t, err)
	require.Equal(t, internal.ErrNoPrevious, err.Error())
}

func TestBack_ErrNoActiveDevice(t *testing.T) {
	api := new(internal.MockAPI)
	api.On("GetPlayback").Return(nil, nil)

	_, err := back(api)
	require.Error(t, err)
	require.Equal(t, internal.ErrNoActiveDevice, err.Error())
}
