package main

import (
	"github.com/Harichandra-Prasath/Tman/pkg/ui"
)

func main() {
	if err := ui.StartUI(); err != nil {
		panic(err)
	}
}
