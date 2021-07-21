package main

import (
	"context"
	"testing"

	"rogchap.com/v8go"
)

func TestConsoleLog(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scriptTasks := make(chan *PlayerScriptTask)
	scriptResponses := make(chan *ScriptResponse)

	go PlayerScriptExecutor(ctx, scriptTasks, scriptResponses)

	scriptTasks <- &PlayerScriptTask{
		PlayerId: 1,
		Script:   "console.log('test')",
		Game:     &Game{},
	}
	response := <-scriptResponses

	if response.Console != "test\n" {
		t.Errorf("Did not get correct console, got %v", response.Console)
	}
	if response.Err != "" {
		t.Errorf("Did get unexpected jsError %v", response.Err)
	}
}

func TestJsError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scriptTasks := make(chan *PlayerScriptTask)
	scriptResponses := make(chan *ScriptResponse)

	go PlayerScriptExecutor(ctx, scriptTasks, scriptResponses)

	scriptTasks <- &PlayerScriptTask{
		PlayerId: 1,
		Script:   "console.log(')",
		Game:     &Game{},
	}
	response := <-scriptResponses

	if response.Err == "" {
		t.Error("Js should show an error but is empty")
	} else {
		t.Logf("Got correct JsError\n%v", response.Err)
	}

}

func TestMoveCommand(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scriptTasks := make(chan *PlayerScriptTask)
	scriptResponses := make(chan *ScriptResponse)

	go PlayerScriptExecutor(ctx, scriptTasks, scriptResponses)

	scriptTasks <- &PlayerScriptTask{
		PlayerId: 1,
		Script:   "cmd.move(2, 1, 3, 5); cmd.attack(2, 1, 3, 5); cmd.build(2, 1, 3, 5, \"tower\")",
		Game:     &Game{},
	}
	response := <-scriptResponses

	if response.Err != "" {
		t.Errorf("Js error: %v", response.Err)
	}

	cmd := response.Commands[0]
	switch v := cmd.(type) {
	case *MoveCommand:
		if response.PlayerId != 1 {
			t.Errorf("playerId should be 1 but is %v", response.PlayerId)
		}
		if v.unit != 2 {
			t.Errorf("unit should be 2 but is %v", v.unit)
		}
		if v.target.X != 3 {
			t.Errorf("X should be 3 but is %v", v.target.X)
		}
		if v.target.Y != 5 {
			t.Errorf("Y should be 5 but is %v", v.target.Y)
		}
	case *BuildCommand:
		if response.PlayerId != 1 {
			t.Errorf("playerId should be 1 but is %v", response.PlayerId)
		}
		if v.unit != 2 {
			t.Errorf("unit should be 2 but is %v", v.unit)
		}
		if v.target.X != 3 {
			t.Errorf("X should be 3 but is %v", v.target.X)
		}
		if v.target.Y != 5 {
			t.Errorf("Y should be 5 but is %v", v.target.Y)
		}
		if v.building != "tower" {
			t.Errorf("bulding should be tower but is %v", v.building)
		}
	case *AttackCommand:
		if response.PlayerId != 1 {
			t.Errorf("playerId should be 1 but is %v", response.PlayerId)
		}
		if v.unit != 2 {
			t.Errorf("unit should be 2 but is %v", v.unit)
		}
		if v.target.X != 3 {
			t.Errorf("X should be 3 but is %v", v.target.X)
		}
		if v.target.Y != 5 {
			t.Errorf("Y should be 5 but is %v", v.target.Y)
		}
	default:
		t.Error("Recived command is not a move command pointer")
	}

}

func TestJsGameInput(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scriptTasks := make(chan *PlayerScriptTask)
	scriptResponses := make(chan *ScriptResponse)

	go PlayerScriptExecutor(ctx, scriptTasks, scriptResponses)

	unit := Unit{
		Id:            554,
		Position:      Position{X: 1, Y: 6},
		Health:        100,
		Attack:        10,
		AttackRange:   1,
		Heal:          0,
		Speed:         1,
		Build:         10,
		Carry:         50,
		Manufactoring: 10,
	}

	room := &Room{
		RoomId:   15,
		Elements: []GameElement{unit},
	}

	scriptTasks <- &PlayerScriptTask{
		PlayerId: 1,
		Script:   "console.log(game.Rooms[0].Elements[0].Id)",
		Game:     &Game{Rooms: []*Room{room}},
	}
	response := <-scriptResponses

	if response.Err != "" {
		t.Errorf("Js error: %v", response.Err)
	}

	if response.Console != "554\n" {
		t.Errorf("Log output differs, got : %v", response.Console)
	}
}

func TestSetArray(t *testing.T) {
	ctx, _ := v8go.NewContext()
	ctx.RunScript("myarray = []", "init.js")
	//vm, _ := ctx.Isolate()
	myarrayV, _ := ctx.Global().Get("myarray")
	myarray, _ := myarrayV.AsObject()

	//array, _ := v8go.NewObjectTemplate(vm)
	//myarray, _ := array.NewInstance(ctx)
	//myarray.Set("length", int32(2))
	myarray.SetIdx(0, int32(42))
	myarray.SetIdx(1, int32(21))

	//global.Set("testarray", myarray)
	val, _ := ctx.RunScript("myarray", "test.js")
	json, _ := v8go.JSONStringify(ctx, val)

	if json != "[42,21]" {
		t.Errorf("Array != [42,21] is %v", json)
	}
}
