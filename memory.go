package client

import (
	"strings"

	"github.com/shirou/gopsutil/mem"
)

// MemoryType is the type for which type of memory that should be checked
type MemoryType int

const (
	RAM  MemoryType = iota // 0
	Swap                   // 1
)

// MemoryResponse contains information about RAM/Swap
type MemoryResponse struct {
	Error string     `json:"error"`
	Size  uint64     `json:"size"`
	Type  MemoryType `json:"type"`
}

// MemoryCheck checks the usage of a specific memory type (Swap or RAM)
func MemoryCheck(cmd Command) MemoryResponse {
	t := RAM
	total := false
	resp := MemoryResponse{Type: t}

	for _, args := range cmd.Params {
		if strings.ToLower(args.Name) == "-swap" {
			t = Swap
		} else if strings.ToLower(args.Name) == "-total" {
			total = true
		}
	}

	switch t {
	case RAM:
		v, err := mem.VirtualMemory()
		if err != nil {
			resp.Error = err.Error()
			return resp
		}

		if total {
			resp.Size = v.Total
		} else {
			resp.Size = v.Used
		}
		break
	case Swap:
		v, err := mem.SwapMemory()
		if err != nil {
			resp.Error = err.Error()
			return resp
		}

		if total {
			resp.Size = v.Total
		} else {
			resp.Size = v.Used
		}
		break
	}

	return resp
}
