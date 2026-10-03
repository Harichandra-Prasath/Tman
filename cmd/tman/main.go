package main

import (
	"github.com/Harichandra-Prasath/Tman/pkg/tman"
	"github.com/Harichandra-Prasath/Tman/pkg/ui"
)

func main() {
	if r := tman.IsReachable(); !r {
		panic("tmux is not reachable. Is tmux installed??")
	}

	if err := ui.StartUI(); err != nil {
		panic(err)
	}
}
