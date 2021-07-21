package main

import "time"

type ListenerCoordinator struct {
	registerListener   chan PlayerViewListener
	unregisterListener chan PlayerViewListener
}

func (lc ListenerCoordinator) register() chan<- PlayerViewListener   { return lc.registerListener }
func (lc ListenerCoordinator) unregister() chan<- PlayerViewListener { return lc.unregisterListener }

func newCoordinator() *ListenerCoordinator {
	return &ListenerCoordinator{
		registerListener:   make(chan PlayerViewListener, 25),
		unregisterListener: make(chan PlayerViewListener, 25),
	}
}

func main() {
	coordinator := newCoordinator()
	listeners := GameListenerSet{}
	timer := time.Tick(250 * time.Millisecond)

	game := &Game{Players: []*Player{{Id: 0, Script: "console.log('HeyHo!')"}}}

	go startWebserver(coordinator)

	for {
		select {
		case listener := <-coordinator.registerListener:
			//value part of map is unused, the map is just used as a set
			listeners[listener] = struct{}{}
		case listener := <-coordinator.unregisterListener:
			delete(listeners, listener)
		case <-timer:
			cmdFails, scriptResps := tick(game)

			// view per player
			playerViews := buildPlayerViews(game, cmdFails, scriptResps)
			for listener := range listeners {
				listener.channel() <- (*playerViews)[listener.forPlayer()]
			}
		}
	}
}

func buildPlayerViews(game *Game, cmdFails *[]*CommandFailure, scriptResps *[]*ScriptResponse) *map[int]*PlayerView {
	playerViews := make(map[int]*PlayerView)
	for _, player := range game.Players {

		playerViews[player.Id] = &PlayerView{
			Rooms:          roomsForPlayer(player, &game.Rooms),
			ScriptResponse: scriptResponseForPlayer(player, scriptResps),
			CmdFails:       cmdFailsForPlayer(player, cmdFails),
		}
	}
	return &playerViews
}

func roomsForPlayer(player *Player, rooms *[]*Room) []*Room {
	playerRooms := []*Room{}
	for _, room := range *rooms {
		for _, elem := range room.Elements {
			if elem.GetPlayer() == player.Id {
				playerRooms = append(playerRooms, room)
				break
			}
		}
	}

	return playerRooms
}

func scriptResponseForPlayer(player *Player, scriptResponses *[]*ScriptResponse) *ScriptResponse {
	for _, scriptResp := range *scriptResponses {
		if scriptResp.PlayerId == player.Id {
			return scriptResp
		}
	}
	return nil
}

func cmdFailsForPlayer(player *Player, cmdFails *[]*CommandFailure) *[]*CommandFailure {
	playerCmdFails := []*CommandFailure{}
	for _, cmdFail := range *cmdFails {
		if cmdFail.command.player() == player.Id {
			playerCmdFails = append(playerCmdFails, cmdFail)
		}
	}
	return &playerCmdFails
}
