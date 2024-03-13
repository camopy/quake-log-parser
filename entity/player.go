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
	// TODO: should we decrement kills if it's already 0?
	p.kills--
}

func (p *Player) rename(name string) {
	p.name = name
}
