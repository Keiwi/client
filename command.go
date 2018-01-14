package client

import (
	"strings"
)

// Command is a struct for a command
type Command struct {
	Name   string
	Params []Argument
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
	commandSplit := strings.Split(command, " ")
	cmd := Command{Name: commandSplit[0], Params: []Argument{}}

	// Check if the string contains flags
	if len(commandSplit) > 1 {
		for i := 1; i < len(commandSplit); i++ {
			argumentSplit := strings.Split(commandSplit[i], "=")
			val := ""
			if len(argumentSplit) > 1 {
				val = argumentSplit[1]
			}
			cmd.Params = append(cmd.Params, Argument{Name: argumentSplit[0], Value: val})
		}
	}

	return cmd
}
