package main

type GameElement interface {
	GetId() int
	Locate() Position
	GetHealth() int
	GetPlayer() int
	TakeDamage(int)
}

type Building interface {
	GameElement
}

type Storage struct {
	Items map[Item]int
}

type Item interface {
	Idable
}

type Energy struct {
	Idable
}

type Food struct {
	Item
}

type Material struct {
	Item
}

type Unit struct {
	Id            int
	Player        int
	Room          int
	Position      Position
	Health        int
	Attack        int
	AttackRange   int
	Heal          int
	Build         int
	Speed         int
	Carry         int
	Manufactoring int
}

func (u Unit) GetId() int       { return u.Id }
func (u Unit) GetHealth() int   { return u.Health }
func (u Unit) GetPlayer() int   { return u.Player }
func (u Unit) Locate() Position { return u.Position }
func (u Unit) TakeDamage(damage int) {
	u.Health -= damage
	if u.Health < 0 {
		u.Health = 0
	}
}

type Blueprint struct {
}

type UnitSpawner struct {
	Id       int
	Player   int
	Room     int
	Position Position
	Health   int
	Energy   int
	Storage  Storage
}

func (u UnitSpawner) GetId() int       { return u.Id }
func (u UnitSpawner) GetHealth() int   { return u.Health }
func (u UnitSpawner) GetPlayer() int   { return u.Player }
func (u UnitSpawner) Locate() Position { return u.Position }
func (u UnitSpawner) TakeDamage(damage int) {
	u.Health -= damage
	if u.Health < 0 {
		u.Health = 0
	}
}
