package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/hay-kot/gojsonsumgen/internal/commands"
)

var (
	// Build information. Populated at build-time via -ldflags flag.
	version = "dev"
	commit  = "HEAD"
	date    = "now"
)

func build() string {
	short := commit
	if len(commit) > 7 {
		short = commit[:7]
	}

	return fmt.Sprintf("%s (%s) %s", version, short, date)
}

func main() {
	ctrl := &commands.Controller{
		Flags: &commands.Flags{},
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	app := &cli.App{
		Name:    "Go JSON Sum Type Generator",
		Usage:   `This is an experimental code generator for generating sum types in go that are JSON marshable`,
		Version: build(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "log level (debug, info, warn, error, fatal, panic)",
				Value:       "debug",
				Destination: &ctrl.Flags.LogLevel,
			},
		},
		Before: func(ctx *cli.Context) error {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			level, err := zerolog.ParseLevel(ctx.String("log-level"))
			if err != nil {
				return fmt.Errorf("failed to parse log level: %w", err)
			}

			log.Logger = log.Level(level)

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "gen",
				Usage: "",
				Action: func(ctx *cli.Context) error {
					return ctrl.Gen(ctx.Context, ctx.Args().Slice())
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run Go JSON Sum Type Generator")
	}
}
