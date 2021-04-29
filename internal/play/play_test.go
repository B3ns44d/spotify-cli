package play

import (
	"errors"
	"spotify/internal"
	"spotify/pkg"
	"spotify/pkg/model"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPlayCommand(t *testing.T) {
	api := new(pkg.MockAPI)

	playback := &model.Playback{
		IsPlaying:  false,
		ProgressMs: 0,
		Item: model.Item{
			Name: "Song",
			Artists: []model.Artist{
				{Name: "Artist"},
			},
			DurationMs: 1000,
		},
	}

	api.On("Status").Run(func(_ mock.Arguments) {
		playback.IsPlaying = true
	}).Return(playback, nil)

	api.On("Play").Return(nil)

	status, err := Play(api)
	require.Equal(t, "🎵 Song\n🎤 Artist\n▶️  0:00 [                ] 0:01\n", status)
	require.NoError(t, err)
}

func TestAlreadyPlayingErr(t *testing.T) {
	api := new(pkg.MockAPI)
	api.On("Status").Return(new(model.Playback), nil)
	api.On("Play").Return(errors.New(internal.RestrictionViolatedSpotifyErr))

	_, err := Play(api)
	require.Equal(t, internal.AlreadyPlayingErr, err.Error())
}

func TestNoActiveDeviceErr(t *testing.T) {
	api := new(pkg.MockAPI)
	api.On("Status").Return(nil, nil)

	_, err := Play(api)
	require.Equal(t, internal.NoActiveDeviceErr, err.Error())
}
