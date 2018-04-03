package commands

import (
	"github.com/shirou/gopsutil/mem"
)

//region MemoryOutput

// MemoryType is the type for which type of memory that should be checked
type MemoryType int

const (
	RAM  MemoryType = iota // 0
	SWAP                   // 1
)

// MemoryOutput contains information about RAM/Swap
type MemoryOutput struct {
	Size uint64
	Type MemoryType
}

func (MemoryOutput) Error() string { return "" }

func (i MemoryOutput) Message() OutputMessage {
	return map[string]interface{}{
		"size": i.Size,
		"type": i.Type,
	}
}

//endregion

//region MemoryCommand

// MemoryCommand checks the usage of a specific memory type (Swap or RAM)
type MemoryCommand struct {
	Version string
}

func (i MemoryCommand) Run(cmd Command) Output {
	resp := MemoryOutput{Type: RAM}

	if cmd.HasArgument("-swap") {
		resp.Type = SWAP
	}

	switch resp.Type {
	case RAM:
		v, err := mem.VirtualMemory()
		if err != nil {
			return ErrorOutput{err: err.Error()}
		}

		if cmd.HasArgument("-total") {
			resp.Size = v.Total
		} else {
			resp.Size = v.Used
		}
		break
	case SWAP:
		v, err := mem.SwapMemory()
		if err != nil {
			return ErrorOutput{err: err.Error()}
		}

		if cmd.HasArgument("-total") {
			resp.Size = v.Total
		} else {
			resp.Size = v.Used
		}
		break
	}

	return resp
}

func (MemoryCommand) Name() string { return "memory" }
func (MemoryCommand) Description() string {
	return "returns information about RAM/Swap memory"
}
func (MemoryCommand) Usage() string { return `[-total|-swap]` }

//endregion
