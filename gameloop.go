package main

import (
	"context"
)

func convertToRoomCommands(scriptResponses []*ScriptResponse) *map[int]*[]*RoomCommand {
	roomCommands := make(map[int]*[]*RoomCommand)

	for _, scriptResponse := range scriptResponses {
		for _, command := range scriptResponse.Commands {
			roomId := command.locateRoom()
			cmdList, exists := roomCommands[roomId]
			if !exists {
				cmdList = &[]*RoomCommand{}
				roomCommands[roomId] = cmdList
			}
			*cmdList = append(*cmdList, &RoomCommand{
				PlayerId: scriptResponse.PlayerId,
				Command:  command,
			})
		}
	}

	return &roomCommands
}

const SCRIPT_ROUTINES = 5
const TRANSITION_ROUTINES = 5

func tick(game *Game) *[]*CommandFailure {
	// setup background execution env for player scripts
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scriptTasks := make(chan *PlayerScriptTask, 50)
	scriptResponses := make(chan *ScriptResponse, 50)

	for i := 0; i < SCRIPT_ROUTINES; i++ {
		go PlayerScriptExecutor(ctx, scriptTasks, scriptResponses)
	}

	// execute player scripts
	for _, p := range game.Players {
		scriptTasks <- &PlayerScriptTask{
			PlayerId: p.Id,
			Game:     game,
			Script:   p.Script,
		}
	}

	// receive script answers
	responses := []*ScriptResponse{}
	for range game.Players {
		responses = append(responses, <-scriptResponses)
	}

	// setup room transition backround workers

	roomTransitionTasks := make(chan *RoomTransitionTask, 50)
	transitionResponses := make(chan *RoomTransitionResponse, 50)

	for i := 0; i < TRANSITION_ROUTINES; i++ {
		go RoomTransitionExecutor(ctx, roomTransitionTasks, transitionResponses, transitionToNextTick)
	}

	// divide commands to their rooms and appy transition to next tick
	roomCommandMap := convertToRoomCommands(responses)

	for rId, room := range game.Rooms {
		cmd, hasCmd := (*roomCommandMap)[rId]
		if !hasCmd {
			cmd = &[]*RoomCommand{}
		}
		roomTransitionTasks <- &RoomTransitionTask{&room, cmd}

	}

	// wait for all rooms to be transitioned
	commandFailures := []*CommandFailure{}
	for range game.Rooms {
		resp := <-transitionResponses
		commandFailures = append(commandFailures, *resp.errors...)
	}

	return &commandFailures

}

type transition func(*Room, *[]*RoomCommand) *[]*CommandFailure

func RoomTransitionExecutor(
	ctx context.Context,
	transitionTasks <-chan *RoomTransitionTask,
	responses chan<- *RoomTransitionResponse,
	transition transition) {

	for {
		select {
		case transitionTask := <-transitionTasks:
			errs := transition(transitionTask.Room, transitionTask.Commands)
			responses <- &RoomTransitionResponse{errors: errs}

		case <-ctx.Done():
			return
		}
	}
}
