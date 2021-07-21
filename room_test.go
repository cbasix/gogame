package main

import "testing"

func TestSpawnBuilding(t *testing.T) {
	room := &Room{}
	pos := Position{Room: 1, X: 3, Y: 4}
	building, err := room.spawnBuilding(1, pos, "spawner")

	elementInRoom, found := room.getElementAt(pos)
	buildingInRoom, isBuilding := elementInRoom.(Building)

	if err != nil || !found || !isBuilding || building != buildingInRoom {
		t.Errorf("Building is not what it was expected to be")
	}
}

func TestSpawnNonExistendBuilding(t *testing.T) {
	room := &Room{}
	_, err := room.spawnBuilding(1, Position{}, "nonExistentType")

	if err == nil {
		t.Errorf("there should be an error here saying that this type doesnt exist")
	}
}
