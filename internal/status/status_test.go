package status

import (
	"spotify/internal"
	"testing"

	"github.com/brianstrauch/spotify"
	"github.com/brianstrauch/spotify/model"
	"github.com/stretchr/testify/require"
)

func TestStatusCommand(t *testing.T) {
	api := new(spotify.MockAPI)

	playback := &model.Playback{
		IsPlaying:  true,
		ProgressMs: 0,
		Item: model.Item{
			Type: "track",
			Name: "Song",
			Artists: []model.Artist{
				{Name: "Artist"},
			},
			DurationMs: 1000,
		},
	}

	api.On("Status").Return(playback, nil)

	status, err := status(api)
	require.Equal(t, "   Song\r🎵\n   Artist\r🎤\n   0:00 [                ] 0:01\r▶️\n", status)
	require.NoError(t, err)
}

func TestMultipleArtists(t *testing.T) {
	api := new(spotify.MockAPI)

	playback := &model.Playback{
		IsPlaying:  true,
		ProgressMs: 0,
		Item: model.Item{
			Type: "track",
			Name: "Song",
			Artists: []model.Artist{
				{Name: "Artist 1"},
				{Name: "Artist 2"},
			},
			DurationMs: 1000,
		},
	}

	api.On("Status").Return(playback, nil)

	status, err := status(api)
	require.Equal(t, "   Song\r🎵\n   Artist 1, Artist 2\r🎤\n   0:00 [                ] 0:01\r▶️\n", status)
	require.NoError(t, err)
}

func TestPodcast(t *testing.T) {
	api := new(spotify.MockAPI)

	playback := &model.Playback{
		IsPlaying:  true,
		ProgressMs: 0,
		Item: model.Item{
			Type: "episode",
			Name: "Episode",
			Show: model.Show{
				Name: "Podcast",
			},
			DurationMs: 1000,
		},
	}

	api.On("Status").Return(playback, nil)

	status, err := status(api)
	require.Equal(t, "   Episode\r🎵\n   Podcast\r🎤\n   0:00 [                ] 0:01\r▶️\n", status)
	require.NoError(t, err)
}

func TestNoActiveDeviceErr(t *testing.T) {
	api := new(spotify.MockAPI)
	api.On("Status").Return(nil, nil)

	_, err := status(api)
	require.Equal(t, internal.NoActiveDeviceErr, err.Error())
}
