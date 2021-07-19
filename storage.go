package main

type PermanantStorage interface {
	storePlayerMemory(playerId int, memory PlayerMemory) error
	loadPlayerMemory(playerId int) (PlayerMemory, error)
	storeGameState(gameName string, game Game) error
	loadGameState(gameName string) (Game, error)
}
