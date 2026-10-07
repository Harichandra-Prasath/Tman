package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/Harichandra-Prasath/Tman/pkg/tman"
	"github.com/Harichandra-Prasath/Tman/pkg/ui"
)

func parseFlags(progName string, args []string) (*tman.TmanConfig, error) {
	fs := flag.NewFlagSet(progName, flag.ContinueOnError)
	cfg := &tman.TmanConfig{}

	fs.StringVar(&cfg.WorkDir, "work-dir", os.Getenv("HOME"), "Working Directory for creating Sessions")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if cfg.WorkDir == "" {
		return nil, fmt.Errorf("invalid workdir for creating sessions")
	}

	_, err := os.Stat(cfg.WorkDir)
	if err != nil {
		return nil, err
	}

	_, err = exec.LookPath("tmux")
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func main() {
	cfg, err := parseFlags(os.Args[0], os.Args[1:])
	if err != nil {
		flagReq := errors.Is(err, flag.ErrHelp)
		if !flagReq {
			fmt.Fprintln(os.Stderr, "tman:", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if err = ui.StartUI(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "tman:", err)
		os.Exit(1)

	}
}
