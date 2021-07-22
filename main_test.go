package main

import (
	"strconv"
	"strings"
	"testing"
)

func TestRoomsForPlayer(t *testing.T) {
	player := &Player{
		Id: 6,
	}
	rooms := &[]*Room{
		{
			RoomId: 1,
			Elements: []GameElement{
				&Unit{Player: 6},
			},
		},
		{
			RoomId: 2,
			Elements: []GameElement{
				&Unit{Player: 4},
				&Unit{Player: 6},
			},
		},
		{
			RoomId:   3,
			Elements: []GameElement{},
		},
		{
			RoomId:   4,
			Elements: []GameElement{&Unit{Player: 4}},
		},
	}

	res := roomsForPlayer(player, rooms)

	var sb strings.Builder
	for _, room := range res {
		sb.WriteString(strconv.Itoa(room.RoomId))
		sb.WriteString(", ")
	}

	if sb.String() != "1, 2, " {
		t.Errorf("expected rooms 1, 2, but got %v", sb.String())
	}
}

func TestScriptsForPlayer(t *testing.T) {
	player := &Player{Id: 8}
	scriptResponses := &[]*ScriptResponse{
		{PlayerId: 7, Console: "1"},
		{PlayerId: 8, Console: "2"},
		{PlayerId: 6, Console: "3"},
		{PlayerId: 8, Console: "4"},
	}

	res := scriptResponseForPlayer(player, scriptResponses)

	if res.Console != "2" {
		t.Errorf("expected response with console 2 but got %v", res.Console)
	}
}

func TestCmdFailsForPlayer(t *testing.T) {
	player := &Player{Id: 8}
	cmdFails := &[]*CommandFailure{
		{Command: AttackCommand{PlayerId: 8, Unit: 1}},
		{Command: AttackCommand{PlayerId: 7, Unit: 2}},
		{Command: BuildCommand{PlayerId: 8, Unit: 3}},
		{Command: AttackCommand{PlayerId: -7, Unit: 4}},
	}

	res := cmdFailsForPlayer(player, cmdFails)

	cmd0Unit := (*res)[0].Command.(AttackCommand).Unit
	if (*res)[0].Command.(AttackCommand).Unit != 1 {
		t.Errorf("expected cmd fail with unit 1 but got %v", cmd0Unit)
	}

	cmd1Unit := (*res)[1].Command.(BuildCommand).Unit
	if cmd1Unit != 3 {
		t.Errorf("expected cmd fail with unit 3 but got %v", cmd1Unit)
	}
}

func TestBuildViewForPlayer(t *testing.T) {
	views := *buildPlayerViews(&Game{Players: []*Player{{Id: 1}}}, &[]*CommandFailure{}, &[]*ScriptResponse{})

	view := views[1]

	if len(*view.CmdFails) != 0 {
		t.Error("expected empty cmd fails")
	}

	if len(view.Rooms) != 0 {
		t.Error("expected empty rooms list")
	}

	if view.ScriptResponse != nil {
		t.Error("expected script response nil")
	}
}
