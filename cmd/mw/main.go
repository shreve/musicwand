package main

import (
	"fmt"
	"os"
	"sort"
	"time"

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
	var player *mpris.Player

	cliApp := cli.App{
		Name:  "mw",
		Usage: "magically control your local media players",
		Before: func(c *cli.Context) (err error) {
			client, err = mpris.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
				os.Exit(1)
			}

			player = client.FindPlayer("musicwand")
			if player == nil {
				// Autostarting daemon
				StartDaemon()

				attempts := 3
				for attempts > 0 {
					player = client.FindPlayer("musicwand")
					if player != nil {
						break
					}
					time.Sleep(1 * time.Second)
					attempts -= 1
				}
			}

			if player == nil {
				fmt.Fprintf(os.Stderr, "Couldn't start the musicwand daemon.\n")
				os.Exit(1)
			}

			return nil
		},
		After: func(c *cli.Context) error {
			client.Close()
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
					fmt.Println("player.Play()")
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
					list := make([]string, 0)
					for key, _ := range meta {
						list = append(list, key)
					}
					sort.Strings(list)
					for _, key := range list {
						fmt.Println(key, "\t", meta[key])
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
						player := client.PlayerWithOwner(event.Sender)
						fmt.Println("CHANGE", player.Name, event.Body[1])
					}
				},
			},
			{
				Name:  "status",
				Usage: "Get a pretty formatted status of current music player",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "watch",
						Aliases: []string{"w"},
						Usage:   "Maintain a process that updates the status when it changes",
					},
					&cli.StringFlag{
						Name:    "format",
						Aliases: []string{"f"},
						Usage:   "Provide a format string to be filled in with data",
						Value:   "{icon} {artist} :: {track}",
					},
				},
				Action: func(c *cli.Context) error {
					fmt.Println(musicwand.FormatStatus(c.String("format"), player))
					if c.Bool("watch") {
						events, _ := client.OnAnyPlayerChange()
						for {
							event := <-events
							player := client.PlayerWithOwner(event.Sender)
							fmt.Println(musicwand.FormatStatus(c.String("format"), player))
						}
					}
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
