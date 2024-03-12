package entity

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGame(t *testing.T) {
	g := NewGame()
	assert.NotNil(t, g)
}

func TestGame_AddPlayer(t *testing.T) {
	g := NewGame()
	g.AddPlayer(1)
	assert.Equal(t, 1, len(g.playersById))

	g.AddPlayer(2)
	assert.Equal(t, 2, len(g.playersById))
}

func TestGame_RenamePlayer(t *testing.T) {
	g := NewGame()
	g.AddPlayer(1)
	err := g.RenamePlayer(1, "Player 1")
	assert.NoError(t, err)
	p, err := g.getPlayerById(1)
	assert.Equal(t, "Player 1", p.name)

	err = g.RenamePlayer(1, "Player 1 Renamed")
	assert.NoError(t, err)
	p, err = g.getPlayerById(1)
	assert.Equal(t, "Player 1 Renamed", p.name)
	p, err = g.getPlayerByName("Player 1 Renamed")
	assert.Equal(t, 1, p.id)
}

func setupGame() (*Game, error) {
	g := NewGame()
	g.AddPlayer(1)
	g.AddPlayer(2)
	g.AddPlayer(3)

	err := g.RenamePlayer(1, "Player 1")
	if err != nil {
		return nil, err
	}
	err = g.RenamePlayer(2, "Player 2")
	if err != nil {
		return nil, err
	}
	err = g.RenamePlayer(3, "Player 3")
	if err != nil {
		return nil, err
	}

	return g, nil
}

func TestGame_AddKill(t *testing.T) {
	g, err := setupGame()
	require.NoError(t, err)

	tests := []struct {
		killer               string
		killed               string
		means                string
		expectedTotalKills   int
		expectedKillsByMeans map[string]int
		expectedKillerKills  int
		expectedKilledKills  int
	}{
		{
			"Player 1",
			"Player 2",
			"MOD_RAILGUN",
			1,
			map[string]int{"MOD_RAILGUN": 1},
			1,
			0,
		},
		{
			"Player 2",
			"Player 3",
			"MOD_TARGET_LASER",
			2,
			map[string]int{
				"MOD_RAILGUN":      1,
				"MOD_TARGET_LASER": 1,
			},
			1,
			0,
		},
		{
			"Player 1",
			"Player 3",
			"MOD_RAILGUN",
			3,
			map[string]int{
				"MOD_RAILGUN":      2,
				"MOD_TARGET_LASER": 1,
			},
			2,
			0,
		},
		{
			worldKill,
			"Player 1",
			"MOD_TRIGGER_HURT",
			4,
			map[string]int{
				"MOD_RAILGUN":      2,
				"MOD_TARGET_LASER": 1,
				"MOD_TRIGGER_HURT": 1,
			},
			0,
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.killer+" kills "+tt.killed, func(t *testing.T) {
			err := g.AddKill(tt.killer, tt.killed, tt.means)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedTotalKills, g.totalKills)
			assert.Equal(t, tt.expectedKillsByMeans, g.killsByMeans)
			if tt.killer != worldKill {
				p, err := g.getPlayerByName(tt.killer)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedKillerKills, p.kills)
			}
			p, err := g.getPlayerByName(tt.killed)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedKilledKills, p.kills)
		})
	}
}

func TestGame_GenerateReport(t *testing.T) {
	g, err := setupGame()
	require.NoError(t, err)

	err = g.AddKill("Player 1", "Player 2", "MOD_RAILGUN")
	require.NoError(t, err)
	err = g.AddKill("Player 2", "Player 3", "MOD_TARGET_LASER")
	require.NoError(t, err)
	err = g.AddKill("Player 1", "Player 3", "MOD_RAILGUN")
	require.NoError(t, err)
	err = g.AddKill(worldKill, "Player 1", "MOD_TRIGGER_HURT")
	require.NoError(t, err)

	report := g.GenerateReport()
	assert.Equal(t, 4, report.TotalKills)
	assert.ElementsMatch(t, []string{"Player 1", "Player 2", "Player 3"}, report.Players)
	assert.True(t, reflect.DeepEqual(map[string]int{
		"Player 1": 1,
		"Player 2": 1,
		"Player 3": 0,
	}, report.Kills))
	assert.True(t, reflect.DeepEqual(map[string]int{
		"MOD_RAILGUN":      2,
		"MOD_TARGET_LASER": 1,
		"MOD_TRIGGER_HURT": 1,
	}, report.KillsByMeans))
}
