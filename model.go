package main

import "fmt"

type Game struct {
	Rooms   []Room
	Players []Player
}

type Position struct {
	Room int
	X    int
	Y    int
}
type Room struct {
	RoomId    int
	Terrain   Terrain
	Units     []Unit
	Buildings []Building
	ExitTo    []int
}

type Terrain struct {
}

type Idable interface {
	GetId() int
}

type Locatable interface {
	locate(game *Game) Position
}

type RoomLocatable interface {
	locateRoom() int
}

type DirectLocatable interface {
	locate() Position
}

type HealthCheckable interface {
	GetHealth() int
}

type GameElementable interface {
	Idable
	Locatable
	HealthCheckable
}

type GameElement struct {
	Id       int
	RoomId   int
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

func (u Unit) GetId() int            { return u.Id }
func (u Unit) GetHealth() int        { return u.Health }
func (u Unit) GetPosition() Position { return u.Position }

type Blueprint struct {
}

type Player struct {
	Id     int
	Name   string
	Script string
}

type PlayerScriptTask struct {
	PlayerId int
	Script   string
	Game     *Game
}

type ScriptResponse struct {
	PlayerId int
	TimedOut bool
	Err      string
	Console  string
	Commands []PlayerCommand
}

type RoomTransitionTask struct {
	Room     *Room
	Commands *[]*RoomCommand
}

type RoomTransitionResponse struct {
	errors *[]*CommandFailure
}

type CommandFailure struct {
	command *RoomCommand
	cause   string
}

type RoomCommand struct {
	PlayerId int
	Command  PlayerCommand
}

type PlayerCommand interface {
	RoomLocatable
	describe() string
}

type MoveCommand struct {
	roomId int
	unit   int
	target Position
}

func (m MoveCommand) describe() string {
	return fmt.Sprintf("move %v to r%vx%vy%v", m.unit, m.roomId, m.target.X, m.target.Y)
}
func (move MoveCommand) locateRoom() int { return move.roomId }

type AttackCommand struct {
	roomId int
	unit   int
	target Position
}

func (m AttackCommand) describe() string {
	return fmt.Sprintf("attack %v to r%vx%vy%v", m.unit, m.roomId, m.target.X, m.target.Y)
}
func (move AttackCommand) locateRoom() int { return move.roomId }

type BuildCommand struct {
	roomId   int
	unit     int
	target   Position
	building string
}

func (m BuildCommand) describe() string {
	return fmt.Sprintf("build %v by %v on r%vx%vy%v", m.building, m.unit, m.roomId, m.target.X, m.target.Y)
}
func (move BuildCommand) locateRoom() int { return move.roomId }

type PlayerMemory struct {
	Data string
}
