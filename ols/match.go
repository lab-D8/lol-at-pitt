package ols

import (
	"github.com/TrevorSStone/goriot"
)

type Participant struct {
	Id            int64
	ParticipantId int
	ChampionId    int
}

type Match struct {
	BlueTeam     string
	RedTeam      string
	Participants []Participant
	Time         string
	Played       bool
	Winner       string
	Id           int64
	Week         int
	Scrim        bool
}

func (m *Match) BlueTeamWin() bool {
	return m.Winner == m.BlueTeam && m.Played
}

func (m *Match) DidRedTeamWin() bool {
	return !m.BlueTeamWin() && m.Played
}

type Matches []*Match

func (m Matches) Len() int {
	return len(m)
}

func (m Matches) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m Matches) Less(i, j int) bool {
	return m[i].Week < m[j].Week
}

func (m *Match) LeagueMatch() goriot.MatchDetail {
	return GetMatchesDAO().LoadLeagueGame(m.Id)
}
