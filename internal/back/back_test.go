package back

import (
	"errors"
	"spotify/internal"
	"spotify/pkg"
	"spotify/pkg/model"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBackCommand(t *testing.T) {
	api := new(pkg.MockAPI)

	playback1 := &model.Playback{
		IsPlaying:  true,
		ProgressMs: 0,
		Item: model.Item{
			ID:   "1",
			Type: "track",
			Name: "Song",
			Artists: []model.Artist{
				{Name: "Artist"},
			},
			DurationMs: 1000,
		},
	}

	playback2 := new(model.Playback)
	*playback2 = *playback1
	playback2.Item.ID = "0"

	api.On("Status").Return(playback1, nil)
	api.On("WaitForUpdatedPlayback", mock.AnythingOfType("func(*model.Playback) bool")).Return(playback2, nil)
	api.On("Back").Return(nil)

	status, err := back(api)
	require.Equal(t, "🎵 Song\n🎤 Artist\n▶️  0:00 [                ] 0:01\n", status)
	require.NoError(t, err)
}

func TestNoPreviousErr(t *testing.T) {
	api := new(pkg.MockAPI)
	api.On("Status").Return(new(model.Playback), nil)
	api.On("Back").Return(errors.New(internal.RestrictionViolatedSpotifyErr))

	_, err := back(api)
	require.Equal(t, internal.NoPreviousErr, err.Error())
}

func TestNoActiveDeviceErr(t *testing.T) {
	api := new(pkg.MockAPI)
	api.On("Status").Return(nil, nil)

	_, err := back(api)
	require.Equal(t, internal.NoActiveDeviceErr, err.Error())
}
