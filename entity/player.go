package entity

type Player struct {
	id    int
	name  string
	kills int
}

func (p *Player) addKill() {
	p.kills++
}

func (p *Player) removeKill() {
	p.kills--
}

func (p *Player) rename(name string) {
	p.name = name
}
