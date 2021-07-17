package main

type Game struct {
	Rooms [][]Room
}

type Position struct {
	X int
	Y int
}
type Room struct {
	Terrain  Terrain
	Elements []GameElement
}

type Terrain struct {
}

type GameElement struct {
	Position Position
	Health   int
}

type Building struct {
	GameElement
	Energy  int
	Storage Storage
}

type Storage struct {
	Items map[Item]int
}

type Item struct {
	Id int
}

type Energy struct {
	Item
}

type Food struct {
	Item
}

type Material struct {
	Item
}

type Unit struct {
	GameElement
	Attack        int
	AttackRange   int
	Heal          int
	Build         int
	Speed         int
	Carry         int
	Manufactoring int
}

type Blueprint struct {
}

type Player struct {
	Id   int
	Name string
}

type PlayerScriptTask struct {
	PlayerId int
	Script   string
	GameJson string
}

type ScriptResponse struct {
	PlayerId int
	TimedOut bool
	Err      string
	Console  string
	Commands []*PlayerCommand
}

type PlayerCommand struct{}

type MoveCommand struct {
	unit   int
	target Position
}

type BuildCommand struct {
	unit     int
	target   Position
	building string
}
