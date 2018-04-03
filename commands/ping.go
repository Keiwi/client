package commands

//region PingOutput

// PingOutput returns "pong" if the client is responsive
type PingOutput struct {
	message string
}

func (PingOutput) Error() string { return "" }

func (i PingOutput) Message() OutputMessage {
	return map[string]interface{}{
		"response": i.message,
	}
}

//endregion

//region PingCommand

// PingCommand returns "pong" if the client is responsive
type PingCommand struct {
}

func (i PingCommand) Run(cmd Command) Output {
	return PingOutput{"pong"}
}

func (PingCommand) Name() string { return "ping" }
func (PingCommand) Description() string {
	return `returns "pong" if the client is responsive`
}
func (PingCommand) Usage() string { return `` }

//endregion
