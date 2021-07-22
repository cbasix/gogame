package main

import (
	"context"
	"testing"
)

func TestConvertToRoomCommands(t *testing.T) {
	scriptResponses := []*ScriptResponse{
		{
			PlayerId: 1,
			Commands: []PlayerCommand{
				MoveCommand{
					RoomId: 0,
					Unit:   33,
					Target: Position{0, 4, 5},
				},
			},
		},
		{
			PlayerId: 1,
			Commands: []PlayerCommand{
				MoveCommand{
					RoomId: 0,
					Unit:   34,
					Target: Position{0, 4, 7},
				},
			},
		},
		{
			PlayerId: 2,
			Commands: []PlayerCommand{
				MoveCommand{
					RoomId: 1,
					Unit:   35,
					Target: Position{1, 2, 3},
				},
			},
		},
	}

	roomCommands := *groupCommandsByRoom(scriptResponses)
	room1Cmds, ok := roomCommands[0]
	if !ok {
		t.Fatalf("commands for room 0 is nil")
	}
	room2Cmds, ok := roomCommands[1]
	if !ok {
		t.Fatalf("commands for room 1 is nil")
	}

	if (*room1Cmds)[0].describe() != "move 33 to r0x4y5" {
		t.Errorf("Room 0 Cmd 0 should be 'move 33 to r0x4y5' but is %v", (*room1Cmds)[0].describe())
	}

	if (*room1Cmds)[1].describe() != "move 34 to r0x4y7" {
		t.Errorf("Room 0 Cmd 1 should be 'move 34 to r0x4y7' but is %v", (*room1Cmds)[1].describe())
	}

	if (*room2Cmds)[0].describe() != "move 35 to r1x2y3" {
		t.Errorf("Room 1 Cmd 0 should be 'move 35 to r1x2y3' but is %v", (*room2Cmds)[0].describe())
	}
}

func TestRoomTransitionExecutor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tasks := make(chan *RoomTransitionTask, 1)
	responses := make(chan *RoomTransitionResponse, 1)

	command := &MoveCommand{}

	go RoomTransitionExecutor(ctx, tasks, responses, func(r *Room, rc *[]PlayerCommand) *[]*CommandFailure {
		return &[]*CommandFailure{{Command: (*rc)[0]}}
	})

	tasks <- &RoomTransitionTask{&Room{}, &[]PlayerCommand{command}}

	result := <-responses
	if (*result.errors)[0].Command != command {
		t.Error("Result should match passed command")
	}
}

func TestIntegrationTick(t *testing.T) {
	unit := &Unit{
		Id:       2,
		Position: Position{Room: 1, X: 3, Y: 4},
	}
	game := &Game{
		Rooms: []*Room{{
			RoomId:   1,
			Elements: []GameElement{unit},
		}},
		Players: []*Player{{Id: 2, Name: "Me", Script: "cmd.move(2, 1, 3, 5);"}},
	}

	tick(game)

	if unit.Position.Y != 5 {
		t.Error("unit is not at expected position")
	}
}
