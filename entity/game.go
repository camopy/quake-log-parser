package entity

import (
	"errors"
	"fmt"
)

var (
	ErrPlayerNotFound = errors.New("Player not found")
	worldKill         = "<world>"
)

type Games []*Game

func (g Games) GenerateReport() GamesReport {
	reports := make(GamesReport)
	for i, game := range g {
		reports[fmt.Sprintf("game-%d", i+1)] = game.GenerateReport()
	}
	return reports
}

type Game struct {
	totalKills    int
	playersById   map[int]*Player
	playersByName map[string]*Player
	killsByMeans  map[string]int
}

func NewGame() *Game {
	return &Game{
		playersById:   make(map[int]*Player),
		playersByName: make(map[string]*Player),
		killsByMeans:  make(map[string]int),
	}
}

func (g *Game) PlayerCount() int {
	return len(g.playersById)
}

func (g *Game) AddPlayer(id int) {
	g.playersById[id] = &Player{
		id: id,
	}
}

func (g *Game) RenamePlayer(id int, name string) error {
	p, err := g.getPlayerById(id)
	if err != nil {
		return err
	}
	if p.name == "" {
		g.playersByName[name] = p
	} else {
		delete(g.playersByName, p.name)
		g.playersByName[name] = p
	}
	p.rename(name)
	return nil
}

func (g *Game) getPlayerById(id int) (*Player, error) {
	if p, ok := g.playersById[id]; ok {
		return p, nil
	}
	return nil, ErrPlayerNotFound
}

func (g *Game) AddKill(killer string, killed string, means string) error {
	g.incrementTotalKills()
	g.incrementKillsByMeans(means)
	// TODO: should we ignore kills where killer == killed?
	//if killer == killed {
	//	return nil
	//}

	if killer == worldKill {
		p, err := g.getPlayerByName(killed)
		if err != nil {
			return err
		}
		p.removeKill()
	} else {
		p, err := g.getPlayerByName(killer)
		if err != nil {
			return err
		}
		p.addKill()
	}

	return nil
}

func (g *Game) getPlayerByName(name string) (*Player, error) {
	if p, ok := g.playersByName[name]; ok {
		return p, nil
	}
	return nil, ErrPlayerNotFound
}

func (g *Game) incrementTotalKills() {
	g.totalKills++
}

func (g *Game) incrementKillsByMeans(means string) {
	g.killsByMeans[means]++
}

func (g *Game) GenerateReport() GameReport {
	return GameReport{
		TotalKills:   g.totalKills,
		Players:      g.playerNames(),
		Kills:        g.killsByPlayer(),
		KillsByMeans: g.killsByMeans,
	}
}

func (g *Game) playerNames() []string {
	var names []string
	for _, p := range g.playersByName {
		names = append(names, p.name)
	}
	return names
}

func (g *Game) killsByPlayer() map[string]int {
	kills := make(map[string]int)
	for _, p := range g.playersByName {
		kills[p.name] = p.kills
	}
	return kills
}
