package main

func transitionToNextTick(room *Room, commands *[]PlayerCommand) *[]*CommandFailure {

	cmdFails := []*CommandFailure{}

	for _, command := range *commands {
		err := command.execute(room)
		cmdFails = append(cmdFails, &CommandFailure{command: command, cause: err.Error()})
	}

	// remove dead bodies / ruins
	/*for _, unit := range room.Units {
		if unit.Health <= 0 {
			room.Units.
		}
	}*/
	return &cmdFails
}
