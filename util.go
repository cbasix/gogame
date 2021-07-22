package main

import (
	"fmt"
	"sync"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func distance(a Position, b Position) (int, error) {
	if a.Room != b.Room {
		return -1, fmt.Errorf("Position are not in same room no distance measuring possible")
	}
	// Units can only move up down left right so distance must respect that
	return Abs(a.X-b.X) + Abs(a.Y-b.Y), nil

}

var lastId = 1
var lastIdMutex sync.Mutex

func generateId() int {
	lastIdMutex.Lock()
	defer lastIdMutex.Unlock()

	currentId := lastId
	lastId++

	return currentId
}

/*func removeUnordered(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}*/

type GameListenerSet map[PlayerViewListener]struct{}
