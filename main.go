package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"runtime"
	"time"
)

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
	log.Print("Starting up")
	start := time.Now()

	coordinator := newCoordinator()
	listeners := GameListenerSet{}
	timer := time.Tick(500 * time.Millisecond)
	statsTimer := time.Tick(30 * time.Second)

	game := &Game{
		Players: []*Player{{Id: 0, Script: `
		Array.prototype.random = function () {
			return this[Math.floor((Math.random()*this.length))];
		}
		

		console.log('HeyHo!'); 
		let unit = game.Rooms[0].Elements[1];
		cmd.move(unit.Id, 0, unit.Position.X + [1, 0, -1].random(), unit.Position.Y + [1, 0, -1].random());`}},
		Rooms: []*Room{
			{
				RoomId: 0,
				Elements: []GameElement{
					&Unit{
						Id:       0,
						Player:   0,
						Position: Position{0, 4, 3},
					},
					&Unit{
						Id:       1,
						Player:   0,
						Position: Position{0, 50, 50},
					},
				},
			},
		},
	}

	go startWebserver(coordinator)

	log.Print("Entering gameloop")
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
		case <-statsTimer:
			PrintMemUsage(start)
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
		if cmdFail.Command.player() == player.Id {
			playerCmdFails = append(playerCmdFails, cmdFail)
		}
	}
	return &playerCmdFails
}

func PrintMemUsage(start time.Time) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	elapsed := time.Since(start)

	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tAlloc/h (since start) = %v MiB", bToMbfloat(float32(m.TotalAlloc)/float32(elapsed.Hours())))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
func bToMbfloat(b float32) float32 {
	return b / 1024 / 1024
}
