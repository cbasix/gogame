package main

func transitionToNextTick(room *Room, commands *[]PlayerCommand) *[]*CommandFailure {

	cmdFails := []*CommandFailure{}

	for _, command := range *commands {
		err := command.execute(room)
		if err != nil {
			fail := &CommandFailure{
				Command: command,
				Cause:   err.Error(),
			}
			cmdFails = append(cmdFails, fail)
		}
	}

	// remove dead bodies / ruins
	/*for _, unit := range room.Units {
		if unit.Health <= 0 {
			room.Units.
		}
	}*/
	return &cmdFails
}
