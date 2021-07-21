package main

import (
	"strings"
	"testing"
)

func TestAttackExecute(t *testing.T) {
	room := &Room{
		RoomId: 1,
	}
	a := AttackCommand{
		roomId: 1,
	}

	err := a.execute(room)

	if err == nil {
		t.Error("Attack with non existing unit should fail")
	}

	if strings.Contains(err.Error(), "nothing") {
		t.Errorf("Wrong error message: %v", err.Error())
	}
}

func TestMoveExecute(t *testing.T) {
	unit := &Unit{
		Id:       8,
		Position: Position{Room: 0, X: 6, Y: 7},
	}

	room := &Room{
		RoomId:   0,
		Elements: []GameElement{unit},
	}

	move := MoveCommand{
		roomId: 0,
		unit:   8,
		target: Position{Room: 0, X: 6, Y: 8},
	}

	err := move.execute(room)
	if err != nil {
		t.Fatal(err)
	}

	if unit.Position.Y != 8 {
		t.Errorf("unit was not moved to y8 but has y%v", unit.Position.Y)
	}
}

func TestMoveExecuteFailOnDist(t *testing.T) {
	unit := &Unit{
		Id:       8,
		Position: Position{Room: 0, X: 6, Y: 7},
	}

	room := &Room{
		Elements: []GameElement{unit},
	}

	// invalid command dist to current unit position must be 1
	move := MoveCommand{
		roomId: 0,
		unit:   8,
		target: Position{Room: 0, X: 6, Y: 9},
	}

	err := move.execute(room)
	if err == nil || !strings.Contains(err.Error(), "Teleport") {
		t.Error("did not get the teleport exception")
	}
}
