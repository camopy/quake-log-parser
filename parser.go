package quake_log_parser

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strconv"

	"github.com/camopy/quake-log-parser/entity"
)

var (
	ErrNoGameFound    = errors.New("no game found")
	ErrNoPlayerFound  = errors.New("no player found on game")
	ErrParsingLogLine = errors.New("error when parsing log line")

	commandInitGame              = "InitGame"
	commandKill                  = "Kill"
	commandClientConnected       = "ClientConnect"
	commandClientUserinfoChanged = "ClientUserinfoChanged"
)

type Parser struct {
	logPath string
	games   []*entity.Game
}

func NewParser(logPath string) *Parser {
	return &Parser{
		logPath: logPath,
		games:   make([]*entity.Game, 0, 10),
	}
}

func (p *Parser) Parse() ([]*entity.Game, error) {
	f, err := os.Open(p.logPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := regexp.Compile(`\d+:\d+ (\w+): (.*)`)
	if err != nil {
		return nil, ErrParsingLogLine
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		match := r.MatchString(line)
		if !match {
			continue
		}
		if err := p.parseLine(line); err != nil {
			return nil, err
		}
	}

	return p.games, nil
}

func (p *Parser) parseLine(line string) error {
	if len(line) == 0 {
		return nil
	}

	r, err := regexp.Compile(`\d+:\d+ (\w+): (.*)`)
	if err != nil {
		return ErrParsingLogLine
	}
	matches := r.FindStringSubmatch(line)
	if len(matches) < 3 {
		return ErrParsingLogLine
	}
	switch matches[1] {
	case commandInitGame:
		p.parseInitGame()
	case commandClientConnected:
		return p.parseClientConnected(matches[2])
	case commandClientUserinfoChanged:
		return p.parseClientUserInfoChanged(matches[2])
	case commandKill:
		return p.parseKill(matches[2])
	}
	return nil
}

func (p *Parser) parseInitGame() {
	p.games = append(p.games, entity.NewGame())
}

func (p *Parser) parseClientConnected(clientId string) error {
	if len(p.games) == 0 {
		return ErrNoGameFound
	}
	game := p.games[len(p.games)-1]

	id, err := strconv.Atoi(clientId)
	if err != nil {
		return err
	}

	game.AddPlayer(id)
	return nil
}

func (p *Parser) parseClientUserInfoChanged(clientInfo string) error {
	if len(p.games) == 0 {
		return ErrNoGameFound
	}
	game := p.games[len(p.games)-1]
	if game.PlayerCount() == 0 {
		return ErrNoPlayerFound
	}

	r, err := regexp.Compile(`(\d+)\sn\\(.*?)\\t\\`)
	if err != nil {
		return err
	}
	matches := r.FindStringSubmatch(clientInfo)
	if len(matches) <= 0 {
		return ErrParsingLogLine
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return err
	}
	return game.RenamePlayer(id, matches[2])
}

func (p *Parser) parseKill(killInfo string) error {
	if len(p.games) == 0 {
		return ErrNoGameFound
	}
	game := p.games[len(p.games)-1]

	r, err := regexp.Compile(`(\d+) (\d+) (\d+): (.*) killed (.*) by (.*)`)
	if err != nil {
		return ErrParsingLogLine
	}
	matches := r.FindStringSubmatch(killInfo)
	if len(matches) < 7 {
		return ErrParsingLogLine
	}
	return game.AddKill(matches[4], matches[5], matches[6])
}
