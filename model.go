package main

type Game struct {
	Rooms   []*Room
	Players []*Player
}

type PlayerViewListener interface {
	channel() chan<- *PlayerView
	forPlayer() int
}

type Position struct {
	Room int
	X    int
	Y    int
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
	Commands *[]PlayerCommand
}

type RoomTransitionResponse struct {
	errors *[]*CommandFailure
}

type CommandFailure struct {
	Command PlayerCommand
	Cause   string
}

// TODO player memory
type PlayerMemory struct {
	Data string
}

type PlayerView struct {
	Rooms          []*Room
	CmdFails       *[]*CommandFailure
	ScriptResponse *ScriptResponse
}
