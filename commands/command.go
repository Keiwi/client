package commands

import "strings"

//region Output

// Output is the interface when a command outputs a message
type Output interface {
	Error() string
	Message() OutputMessage
}

// OutputMessage is a simple map[string]interface{} when outputting a message
type OutputMessage map[string]interface{}

//endregion

//region ICommand

// ICommand is the interface for all internal commands
type ICommand interface {
	Run(Command) Output
	Name() string
	Description() string
	Usage() string
}

/**
type EmptyCommand struct{}

func (e EmptyCommand) Run(cmd Command) Output {}
func (e EmptyCommand) Name() string {}
func (e EmptyCommand) Description() string {}
func (e EmptyCommand) Usage() string {}
*/

//endregion

//region Command

// Command is a struct for a command
type Command struct {
	Name       string
	Arguments  []*Argument
	RawMessage string
}

// GetArgument tries to find a specific argument in the argument list
func (c Command) GetArgument(name string) *Argument {
	name = strings.ToLower(name)
	for _, args := range c.Arguments {
		if strings.ToLower(args.Name) == name {
			return args
		}
	}
	return nil
}

// HasArgument first tries to find a specific argument and if found
// it will check if the value is a bool if so it will return the value
func (c Command) HasArgument(name string) bool {
	arg := c.GetArgument(name)
	if arg == nil {
		return false
	}
	if strings.ToLower(arg.Value) == "false" {
		return false
	}
	return true
}

// Argument is the struct for command arguments/flags
type Argument struct {
	Name  string
	Value string
}

// ParseCommand parses a raw command string and
// create a Command struct from the string
// Example:
//		Input: Check_cpu -cpu=1
// 		Output: {
//			Name: "Check_cpu"
//			Params: [
//				{ Name: "-cpu", Value: "1" }
//			]
//		}
func ParseCommand(command string) Command {
	command = strings.TrimRight(command, "\n")
	commandSplit := strings.Split(command, " ")
	cmd := Command{Name: commandSplit[0], Arguments: []*Argument{}, RawMessage: strings.Join(commandSplit[:1], " ")}

	// Check if the string contains flags
	if len(commandSplit) > 1 {
		for i := 1; i < len(commandSplit); i++ {
			argumentSplit := strings.Split(commandSplit[i], "=")
			val := ""
			if len(argumentSplit) > 1 {
				val = argumentSplit[1]
			}
			cmd.Arguments = append(cmd.Arguments, &Argument{Name: argumentSplit[0], Value: val})
		}
	}

	return cmd
}

//endregion

//region ErrorOutput

// Error Output is a simple struct for the Output interface when you need to return a simple error
type ErrorOutput struct{ err string }

func (e ErrorOutput) Error() string          { return e.err }
func (e ErrorOutput) Message() OutputMessage { return nil }

//endregion

//region CommandHandler

// NewCommandHandler creates a new instance of CommandHandler and adds all existing commands to the handler
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		commands: []ICommand{
			CPUCommand{},
			FileCommand{},
			InfoCommand{Version: "0.0.0"},
			MemoryCommand{},
			NetworkCommand{},
			PartitionCommand{},
			UptimeCommand{},
			PingCommand{},
		},
	}
}

// CommandHandler manages all of the commands
type CommandHandler struct{ commands []ICommand }

// GetHelp will return a string with a help page of all description, usage etc. of all the commands
func (c CommandHandler) GetHelp() string {
	return "Not implemented yet"
}

// RunCommand will try to find the correct command and run the command
func (c CommandHandler) RunCommand(command Command) Output {
	for _, cmd := range c.commands {
		if cmd.Name() == command.Name {
			return cmd.Run(command)
		}
	}
	return ErrorOutput{err: "invalid command"}
}

//endregion
