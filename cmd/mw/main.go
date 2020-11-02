package main

import (
	"fmt"
	"os"

	"github.com/shreve/musicwand/internal/pkg/musicwand"
	"github.com/shreve/musicwand/pkg/mpris"
	"github.com/urfave/cli/v2"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "daemon" {
		RunDaemon()
		os.Exit(0)
	}

	var client *mpris.Client
	var app *mpris.App
	var player mpris.Player

	cliApp := cli.App{
		Name:  "mw",
		Usage: "magically control your local media players",
		Before: func(c *cli.Context) error {
			client, _ = mpris.NewClient()
			// if err != nil {
			// 	fmt.Fprintf(os.Stderr, err.Error())
			// 	os.Exit(1)
			// }

			app = client.FindApp("musicwand")
			if app == nil {
				fmt.Fprintf(os.Stderr, "Couldn't find the musicwand daemon\n")
				os.Exit(1)
			}
			player = app.Player()
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "player"},
		},
		Commands: []*cli.Command{
			{
				Name:    "play",
				Aliases: []string{"y"},
				Usage:   "Instruct the player to play",
				Action: func(c *cli.Context) error {
					player.Play()
					return nil
				},
			},
			{
				Name:    "pause",
				Aliases: []string{"u"},
				Usage:   "Instruct the player to pause",
				Action: func(c *cli.Context) error {
					player.Pause()
					return nil
				},
			},
			{
				Name:    "play-pause",
				Aliases: []string{"p"},
				Usage:   "Instruct the player to play or pause based on current state",
				Action: func(c *cli.Context) error {
					player.PlayPause()
					return nil
				},
			},
			{
				Name:    "next",
				Aliases: []string{"n"},
				Usage:   "Instruct the player to play the next media",
				Action: func(c *cli.Context) error {
					player.Next()
					return nil
				},
			},
			{
				Name:    "previous",
				Aliases: []string{"prev", "v"},
				Usage:   "Instruct the player to play the previous media",
				Action: func(c *cli.Context) error {
					player.Previous()
					return nil
				},
			},
			{
				Name:    "stop",
				Aliases: []string{"s"},
				Usage:   "Instruct the player to stop",
				Action: func(c *cli.Context) error {
					player.Stop()
					return nil
				},
			},
			{
				Name:      "open",
				Aliases:   []string{"o"},
				Usage:     "Instruct the player to open the provided URI",
				ArgsUsage: "[URI]",
				Action: func(c *cli.Context) error {
					player.OpenUri(c.Args().Get(0))
					return nil
				},
			},
			{
				Name:  "metadata",
				Usage: "Get all available metadata about the current media",
				Action: func(c *cli.Context) error {
					meta := player.RawMetadata()
					for key, value := range meta {
						fmt.Println(key, "\t", value)
					}
					return nil
				},
			},
			{
				Name:  "daemon",
				Usage: "Run the musicwand control daemon",
				Action: func(c *cli.Context) error {
					RunDaemon()
					return nil
				},
			},
			{
				Name:  "watch",
				Usage: "Tail a log of events monitored by this library",
				Action: func(c *cli.Context) error {
					events, _ := client.OnAnyPlayerChange()
					for {
						event := <-events
						app := client.AppWithOwner(event.Sender)
						fmt.Println("CHANGE", app.Name, event.Body[1])
					}
				},
			},
			{
				Name:  "status",
				Usage: "Get a pretty formatted status of current music player",
				Action: func(c *cli.Context) error {
					fmt.Println(musicwand.FormatStatus("{icon} {artist} :: {track} ({position}/{length})", app))
					return nil
				},
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
