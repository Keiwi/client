package commands

import (
	"os"
	"time"
)

//region PartitionOutput

// File contains information about a specific file
type File struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	IsDir   bool      `json:"is_dir"`
	Perms   string    `json:"perms"`
	LastMod time.Time `json:"last_mod"`
}

// TODO: Maybe add support for multiple files?

// FileOutput contains information about a file
type FileOutput struct {
	File File
}

func (FileOutput) Error() string { return "" }

func (f FileOutput) Message() OutputMessage {
	return map[string]interface{}{
		"file": f.File,
	}
}

//endregion

//region FileCommand

// FileCommand checks information about a specific file
type FileCommand struct {
}

func (c FileCommand) Run(cmd Command) Output {
	file := cmd.GetArgument("-file")
	if file == nil {
		return ErrorOutput{err: "Invalid value on the flag -file"}
	}

	f, err := os.Open(file.Value)
	if err != nil {
		return ErrorOutput{err: err.Error()}
	}

	fi, err := f.Stat()
	if err != nil {
		return ErrorOutput{err: err.Error()}
	}

	resp := FileOutput{
		File: File{
			Name:    fi.Name(),
			Size:    fi.Size(),
			IsDir:   fi.IsDir(),
			Perms:   fi.Mode().String(),
			LastMod: fi.ModTime(),
		},
	}
	return resp
}

func (c FileCommand) Name() string { return "file" }
func (c FileCommand) Description() string {
	return "returns information about a specific file"
}
func (c FileCommand) Usage() string { return `[-file="path"]` }

//endregion
