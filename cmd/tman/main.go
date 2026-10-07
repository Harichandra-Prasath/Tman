package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Harichandra-Prasath/Tman/pkg/tman"
	"github.com/Harichandra-Prasath/Tman/pkg/ui"
)

func parseFlags(progName string, args []string) error {
	fs := flag.NewFlagSet(progName, flag.ContinueOnError)
	cfg := &tman.TmanConfig{}

	fs.StringVar(&cfg.WorkDir, "work-dir", os.Getenv("HOME"), "Working Directory for creating Sessions")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if cfg.WorkDir == "" {
		return fmt.Errorf("invalid workdir for creating sessions")
	}

	_, err := os.Stat(cfg.WorkDir)
	if err != nil {
		return err
	}

	tman.InitialiseConfig(cfg)

	return nil
}

func main() {
	err := parseFlags(os.Args[0], os.Args[1:])
	if err != nil {
		panic(err)
	}

	if err = ui.StartUI(); err != nil {
		panic(err)
	}
}
