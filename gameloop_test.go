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
					roomId: 0,
					unit:   33,
					target: Position{0, 4, 5},
				},
			},
		},
		{
			PlayerId: 1,
			Commands: []PlayerCommand{
				MoveCommand{
					roomId: 0,
					unit:   34,
					target: Position{0, 4, 7},
				},
			},
		},
		{
			PlayerId: 2,
			Commands: []PlayerCommand{
				MoveCommand{
					roomId: 1,
					unit:   35,
					target: Position{1, 2, 3},
				},
			},
		},
	}

	roomCommands := *convertToRoomCommands(scriptResponses)
	room1Cmds, ok := roomCommands[0]
	if !ok {
		t.Fatalf("commands for room 0 is nil")
	}
	room2Cmds, ok := roomCommands[1]
	if !ok {
		t.Fatalf("commands for room 1 is nil")
	}

	if (*room1Cmds)[0].Command.describe() != "move 33 to r0x4y5" {
		t.Errorf("Room 0 Cmd 0 should be 'move 33 to r0x4y5' but is %v", (*room1Cmds)[0].Command.describe())
	}

	if (*room1Cmds)[1].Command.describe() != "move 34 to r0x4y7" {
		t.Errorf("Room 0 Cmd 1 should be 'move 34 to r0x4y7' but is %v", (*room1Cmds)[1].Command.describe())
	}

	if (*room2Cmds)[0].Command.describe() != "move 35 to r1x2y3" {
		t.Errorf("Room 1 Cmd 0 should be 'move 35 to r1x2y3' but is %v", (*room2Cmds)[0].Command.describe())
	}
}

func TestRoomTransitionExecutor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tasks := make(chan *RoomTransitionTask, 1)
	responses := make(chan *RoomTransitionResponse, 1)

	command := &RoomCommand{}

	go RoomTransitionExecutor(ctx, tasks, responses, func(r *Room, rc *[]*RoomCommand) *[]*CommandFailure {
		return &[]*CommandFailure{&CommandFailure{command: (*rc)[0]}}
	})

	tasks <- &RoomTransitionTask{&Room{}, &[]*RoomCommand{command}}

	result := <-responses
	if (*result.errors)[0].command != command {
		t.Error("Result should match passed command")
	}

}
