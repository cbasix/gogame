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
	defer jsVm.Dispose()

	for {
		for {
			select {
			case scriptTask := <-scriptTasks:
				commands, console, err := executeScript(jsVm, scriptTask)

				scriptResponses <- &ScriptResponse{
					PlayerId: scriptTask.PlayerId,
					Err:      errToString(err),
					Commands: commands,
					Console:  console,
				}

			case <-ctx.Done():
				return
			}
		}

	}
}

func executeScript(jsVm *v8go.Isolate, scriptTask *PlayerScriptTask) ([]PlayerCommand, string, *v8go.JSError) {
	jsCtx, _ := v8go.NewContext(jsVm) // creates a new V8 context with a new Isolate aka VM
	defer jsCtx.Close()

	vals := make(chan *v8go.Value, 1)
	errs := make(chan *v8go.JSError, 1)
	console := make(chan string, 50)
	commands := make(chan PlayerCommand, 50)

	// inject the game / log objects into the js global object
	injectGameObject(jsCtx, scriptTask.Game, commands)
	insertConsole(jsCtx, console)

	// execute script in its own goroutine to allow canceling it
	go func() {

		val, err := jsCtx.RunScript(scriptTask.Script, "playerscript.js") // execute the player given script
		if err != nil {
			errs <- err.(*v8go.JSError)
			return
		}
		vals <- val
	}()

	var consoleBuilder strings.Builder
	commandList := []PlayerCommand{}

	for {
		select {
		case <-vals:
			// successfull exit
			return commandList, consoleBuilder.String(), nil

		case err := <-errs:
			// js error exit
			return commandList, consoleBuilder.String(), err

		case <-time.After(500 * time.Millisecond):
			// timeout exit
			vm, _ := jsCtx.Isolate() // get the Isolate from the context
			vm.TerminateExecution()  // terminate the execution
			err := <-errs            // will get a termination error back from the running script
			return commandList, consoleBuilder.String(), err

		case text := <-console:
			// capture console log output
			consoleBuilder.WriteString(text)
			consoleBuilder.WriteString("\n")

		case command := <-commands:
			// capture command calls
			commandList = append(commandList, command)
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

func insertGameControls(jsCtx *v8go.Context, commandChannel chan<- PlayerCommand) {
	/*vm, err := jsCtx.Isolate()
	if err != nil {
		panic(err)
	}

	game, _ := v8go.NewObjectTemplate(vm)
	command, _ := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		commandType := args[0].String()
		switch commandType {
		case "move":
			commandChannel <- &MoveCommand{
				unit:   int(args[1].Int32()),
				target: Position{int(args[2].Int32()), int(args[3].Int32())},
			}
		case "build":
			commandChannel <- &BuildCommand{
				unit:   int(args[1].Int32()),
				target: Position{int(args[2].Int32()), int(args[3].Int32())},
				building: args[4].String(),
			}
		}

		return nil
	})
	game.Set("sendCommand", command)


	gameObj, _ := game.NewInstance(jsCtx)

	jsCtx.Global().Set("game", gameObj)*/
}
