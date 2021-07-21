package main

import "fmt"

type Room struct {
	RoomId   int
	Terrain  Terrain
	Elements []GameElement
	ExitTo   []int
}

func (r *Room) findUnit(unitId int) (u *Unit, ok bool) {
	for _, e := range r.Elements {
		unit, ok := e.(*Unit)
		if ok && unit.Id == unitId {
			return unit, true
		}
	}
	return nil, false
}

func (r *Room) getElementAt(pos Position) (u GameElement, ok bool) {
	for _, elem := range r.Elements {
		elemPos := elem.Locate()
		if elemPos.X == pos.X && elemPos.Y == pos.Y {
			return elem, true
		}
	}
	return nil, false
}

func (r *Room) isFree(pos Position) bool {
	_, found := r.getElementAt(pos)
	return !found
}

func (r *Room) spawnBuilding(player int, pos Position, buildingType string) (Building, error) {
	var newBuilding Building

	switch buildingType {
	case "spawner":
		newBuilding = &UnitSpawner{
			Id:       generateId(),
			Player:   player,
			Room:     r.RoomId,
			Position: pos,
			Health:   100,
			Energy:   0,
		}
	default:
		return nil, fmt.Errorf("unknown building type %v", buildingType)
	}

	r.Elements = append(r.Elements, newBuilding)
	return newBuilding, nil
}
