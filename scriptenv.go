package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"rogchap.com/v8go"
)

func PlayerScriptExecutor(
	ctx context.Context,
	scriptTasks <-chan *PlayerScriptTask,
	scriptResponses chan<- *ScriptResponse) {

	jsVm, _ := v8go.NewIsolate() // creates a new JavaScript VM

	for {
		scriptTask := <-scriptTasks
		commands, console, err := executeScript(jsVm, scriptTask)

		scriptResponses <- &ScriptResponse{
			PlayerId: scriptTask.PlayerId,
			Err:      errToString(err),
			Commands: commands,
			Console:  console,
		}
	}
}

func executeScript(jsVm *v8go.Isolate, scriptTask *PlayerScriptTask) ([]*PlayerCommand, string, *v8go.JSError) {
	jsCtx, _ := v8go.NewContext(jsVm) // creates a new V8 context with a new Isolate aka VM

	vals := make(chan *v8go.Value, 1)
	errs := make(chan *v8go.JSError, 1)
	console := make(chan string, 50)

	// inject the game object into the js global object
	jsCtx.Global().Set("game", "v1.0.0")

	go func() {
		insertConsole(jsCtx, console)
		val, err := jsCtx.RunScript(scriptTask.Script, "playerscript.js") // execute the player given script
		if err != nil {
			errs <- err.(*v8go.JSError)
			return
		}
		vals <- val
	}()

	var consoleBuilder strings.Builder

	for {
		select {
		case _ = <-vals:
			return []*PlayerCommand{}, consoleBuilder.String(), nil
		case err := <-errs:
			return nil, consoleBuilder.String(), err
		case <-time.After(500 * time.Millisecond):
			vm, _ := jsCtx.Isolate() // get the Isolate from the context
			vm.TerminateExecution()  // terminate the execution
			err := <-errs            // will get a termination error back from the running script
			return nil, consoleBuilder.String(), err
		case text := <-console:
			consoleBuilder.WriteString(text)
			consoleBuilder.WriteString("\n")
		}
	}
}

func errToString(err *v8go.JSError) string {
	if err != nil {
		return fmt.Sprintf("%v %v \n%v ", err.Message, err.Location, err.StackTrace)
	} else {
		return ""
	}
}

func insertConsole(jsCtx *v8go.Context, consoleChannel chan<- string) {
	vm, err := jsCtx.Isolate()
	if err != nil {
		panic(err)
	}

	console, _ := v8go.NewObjectTemplate(vm)
	logfn, _ := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		consoleChannel <- fmt.Sprintf("%v", info.Args()[0])
		return nil
	})
	console.Set("log", logfn)
	consoleObj, _ := console.NewInstance(jsCtx)

	jsCtx.Global().Set("console", consoleObj)
}
