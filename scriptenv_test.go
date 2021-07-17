package main

import (
	"context"
	"testing"
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
		GameJson: "",
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
		GameJson: "",
	}
	response := <-scriptResponses

	if response.Err == "" {
		t.Error("Js should show an error but is empty")
	} else {
		t.Logf("Got correct JsError\n%v", response.Err)
	}

}
