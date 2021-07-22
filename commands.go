package main

import "fmt"

type PlayerCommand interface {
	RoomLocatable
	describe() string
	execute(r *Room) error
	player() int
}

type MoveCommand struct {
	RoomId   int
	Unit     int
	Target   Position
	PlayerId int
}

func (m MoveCommand) describe() string {
	return fmt.Sprintf("move %v to r%vx%vy%v", m.Unit, m.RoomId, m.Target.X, m.Target.Y)
}
func (move MoveCommand) locateRoom() int { return move.RoomId }
func (move MoveCommand) player() int     { return move.PlayerId }

func (m MoveCommand) execute(r *Room) error {
	if m.Target.Room != r.RoomId {
		return fmt.Errorf("target must be in same room")
	}

	if !r.isFree(m.Target) {
		return fmt.Errorf("target position r%vx%vy%v is not free", r.RoomId, m.Target.X, m.Target.Y)
	}

	unit, ok := r.findUnit(m.Unit)
	if !ok {
		return fmt.Errorf("no unit with id %v in room %v", m.Unit, r.RoomId)
	}

	if unit.Player != m.player() {
		return fmt.Errorf("unit %v doesnt belong to you", m.Unit)
	}

	dist, _ := distance(m.Target, unit.Position)
	if dist > 1 {
		// TODO implement pathfinding instead ?
		return fmt.Errorf("teleporting not possible")
	}

	unit.Position = m.Target
	return nil
}

type AttackCommand struct {
	RoomId   int
	Unit     int
	Target   Position
	PlayerId int
}

func (a AttackCommand) describe() string {
	return fmt.Sprintf("attack %v to r%vx%vy%v", a.Unit, a.RoomId, a.Target.X, a.Target.Y)
}
func (a AttackCommand) locateRoom() int { return a.RoomId }
func (a AttackCommand) player() int     { return a.PlayerId }

func (a AttackCommand) execute(r *Room) error {
	if a.Target.Room != r.RoomId {
		return fmt.Errorf("target must be in same room")
	}

	unit, ok := r.findUnit(a.Unit)
	if !ok {
		return fmt.Errorf("no unit with id %v in room %v", a.Unit, r.RoomId)
	}

	if unit.Player != a.player() {
		return fmt.Errorf("unit %v doesnt belong to you", a.Unit)
	}

	dist, _ := distance(a.Target, unit.Position)
	if dist > unit.AttackRange {
		// TODO implement pathfinding instead ?
		return fmt.Errorf("target is out of range")
	}

	elem, found := r.getElementAt(a.Target)
	if found {
		elem.TakeDamage(unit.Attack)
	} else {
		return fmt.Errorf("there is nothing at r%vx%vy%v to attack", r.RoomId, a.Target.X, a.Target.Y)
	}

	return nil
}

type BuildCommand struct {
	RoomId   int
	Unit     int
	Target   Position
	Building string
	PlayerId int
}

func (b BuildCommand) describe() string {
	return fmt.Sprintf("build %v by %v on r%vx%vy%v", b.Building, b.Unit, b.RoomId, b.Target.X, b.Target.Y)
}
func (b BuildCommand) locateRoom() int { return b.RoomId }

func (b BuildCommand) player() int { return b.PlayerId }

func (b BuildCommand) execute(r *Room) error {
	if b.Target.Room != r.RoomId {
		return fmt.Errorf("target must be in same room")
	}

	if !r.isFree(b.Target) {
		return fmt.Errorf("build target r%vx%vy%v is already occupied", r.RoomId, b.Target.X, b.Target.Y)
	}

	unit, ok := r.findUnit(b.Unit)
	if !ok {
		return fmt.Errorf("no unit with id %v in room %v", b.Unit, r.RoomId)
	}

	if unit.Player != b.player() {
		return fmt.Errorf("unit %v doesnt belong to you", b.Unit)
	}

	dist, _ := distance(b.Target, unit.Position)
	if dist > 1 {
		// TODO implement pathfinding instead ?
		return fmt.Errorf("target is out of range")
	}

	// TODO non fixed player id command must now its players id
	r.spawnBuilding(1, b.Target, b.Building)

	return nil
}
