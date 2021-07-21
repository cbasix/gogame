package main

import "fmt"

type PlayerCommand interface {
	RoomLocatable
	describe() string
	execute(r *Room) error
	player() int
}

type MoveCommand struct {
	roomId   int
	unit     int
	target   Position
	playerId int
}

func (m MoveCommand) describe() string {
	return fmt.Sprintf("move %v to r%vx%vy%v", m.unit, m.roomId, m.target.X, m.target.Y)
}
func (move MoveCommand) locateRoom() int { return move.roomId }
func (move MoveCommand) player() int     { return move.playerId }

func (m MoveCommand) execute(r *Room) error {
	if m.target.Room != r.RoomId {
		return fmt.Errorf("Target must be in same room")
	}

	if !r.isFree(m.target) {
		return fmt.Errorf("Target position r%vx%vy%v is not free", r.RoomId, m.target.X, m.target.Y)
	}

	unit, ok := r.findUnit(m.unit)
	if !ok {
		return fmt.Errorf("no unit with id %v in room %v", m.unit, r.RoomId)
	}

	if unit.Player != m.player() {
		return fmt.Errorf("unit %v doesnt belong to you", m.unit)
	}

	dist, _ := distance(m.target, unit.Position)
	if dist > 1 {
		// TODO implement pathfinding instead ?
		return fmt.Errorf("Teleporting not possible")
	}

	unit.Position = m.target
	return nil
}

type AttackCommand struct {
	roomId   int
	unit     int
	target   Position
	playerId int
}

func (a AttackCommand) describe() string {
	return fmt.Sprintf("attack %v to r%vx%vy%v", a.unit, a.roomId, a.target.X, a.target.Y)
}
func (a AttackCommand) locateRoom() int { return a.roomId }
func (a AttackCommand) player() int     { return a.playerId }

func (a AttackCommand) execute(r *Room) error {
	if a.target.Room != r.RoomId {
		return fmt.Errorf("Target must be in same room")
	}

	unit, ok := r.findUnit(a.unit)
	if !ok {
		return fmt.Errorf("No unit with id %v in room %v", a.unit, r.RoomId)
	}

	if unit.Player != a.player() {
		return fmt.Errorf("unit %v doesnt belong to you", a.unit)
	}

	dist, _ := distance(a.target, unit.Position)
	if dist > unit.AttackRange {
		// TODO implement pathfinding instead ?
		return fmt.Errorf("Target is out of range")
	}

	elem, found := r.getElementAt(a.target)
	if found {
		elem.TakeDamage(unit.Attack)
	} else {
		return fmt.Errorf("There is nothing at r%vx%vy%v to attack", r.RoomId, a.target.X, a.target.Y)
	}

	return nil
}

type BuildCommand struct {
	roomId   int
	unit     int
	target   Position
	building string
	playerId int
}

func (b BuildCommand) describe() string {
	return fmt.Sprintf("build %v by %v on r%vx%vy%v", b.building, b.unit, b.roomId, b.target.X, b.target.Y)
}
func (b BuildCommand) locateRoom() int { return b.roomId }

func (b BuildCommand) player() int { return b.playerId }

func (b BuildCommand) execute(r *Room) error {
	if b.target.Room != r.RoomId {
		return fmt.Errorf("Target must be in same room")
	}

	if !r.isFree(b.target) {
		return fmt.Errorf("Build target r%vx%vy%v is already occupied", r.RoomId, b.target.X, b.target.Y)
	}

	unit, ok := r.findUnit(b.unit)
	if !ok {
		return fmt.Errorf("No unit with id %v in room %v", b.unit, r.RoomId)
	}

	if unit.Player != b.player() {
		return fmt.Errorf("unit %v doesnt belong to you", b.unit)
	}

	dist, _ := distance(b.target, unit.Position)
	if dist > 1 {
		// TODO implement pathfinding instead ?
		return fmt.Errorf("Target is out of range")
	}

	// TODO non fixed player id command must now its players id
	r.spawnBuilding(1, b.target, b.building)

	return nil
}
