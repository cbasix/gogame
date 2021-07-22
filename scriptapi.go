package main

import (
	"encoding/json"
	"fmt"

	"rogchap.com/v8go"
)

func injectGameObject(jsCtx *v8go.Context, game *Game, commands chan<- PlayerCommand) {
	vm, _ := jsCtx.Isolate()

	gameJs, err := json.Marshal(game)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("game = %v; game", string(gameJs))
	_, jsErr := jsCtx.RunScript(fmt.Sprintf("game = %v; game", string(gameJs)), "init.js")
	if jsErr != nil {
		panic(jsErr)
	}
	//gameObj, _ := gameObjVal.AsObject()

	unitMove, _ := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) == 4 {
			// TODO Error handling here
			commands <- &MoveCommand{
				Unit:   int(args[0].Int32()),
				Target: Position{int(args[1].Int32()), int(args[2].Int32()), int(args[3].Int32())},
			}

		}
		return nil

	})

	unitAttack, _ := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) == 4 {
			commands <- &AttackCommand{
				Unit:   int(args[0].Int32()),
				Target: Position{int(args[1].Int32()), int(args[2].Int32()), int(args[3].Int32())},
			}
		}
		return nil
	})

	unitBuild, _ := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) == 5 {
			commands <- &BuildCommand{
				Unit:     int(args[0].Int32()),
				Target:   Position{int(args[1].Int32()), int(args[2].Int32()), int(args[3].Int32())},
				Building: args[1].String(),
			}
		}
		return nil
	})

	cmdTmp, _ := v8go.NewObjectTemplate(vm)
	cmdTmp.Set("move", unitMove)
	cmdTmp.Set("attack", unitAttack)
	cmdTmp.Set("build", unitBuild)

	cmd, _ := cmdTmp.NewInstance(jsCtx)
	jsCtx.Global().Set("cmd", cmd)
}
