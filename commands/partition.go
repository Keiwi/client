package commands

import (
	"github.com/shirou/gopsutil/disk"
)

//region PartitionOutput

// Partition contains information about a specific partition
type Partition struct {
	Name string `json:"name"`
	Size uint64 `json:"size"`
	Type string `json:"type"`
}

// PartitionOutput contains information about partitions
type PartitionOutput struct {
	Partitions []Partition
}

func (PartitionOutput) Error() string { return "" }

func (p PartitionOutput) Message() OutputMessage {
	return map[string]interface{}{
		"partitions": p.Partitions,
	}
}

//endregion

//region PartitionCommand

// PartitionCommand gathers information about partitions on the current system
type PartitionCommand struct {
}

func (c PartitionCommand) Run(cmd Command) Output {
	part, err := disk.Partitions(false)
	if err != nil {
		return ErrorOutput{err: err.Error()}
	}

	resp := PartitionOutput{}
	total := cmd.HasArgument("-total")

	for _, p := range part {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			return ErrorOutput{err: err.Error()}
		}

		partition := Partition{
			Name: p.Mountpoint,
			Type: p.Fstype,
		}

		if total {
			partition.Size = usage.Total
		} else {
			partition.Size = usage.Used
		}

		resp.Partitions = append(resp.Partitions, partition)
	}
	return resp
}

func (c PartitionCommand) Name() string { return "partition" }
func (c PartitionCommand) Description() string {
	return "returns information about partitions on the current system"
}
func (c PartitionCommand) Usage() string { return "[-total]" }

//endregion
