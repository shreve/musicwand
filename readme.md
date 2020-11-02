MusicWand
=========

Control your local music players with the [MPRIS D-Bus
Interface](https://specifications.freedesktop.org/mpris-spec/latest/)

The current goal is to basically port
[playerctl](https://github.com/altdesktop/playerctl) to Go.

Once stable, some quality of life enhancements are planned on top.

## Usage

```
NAME:
   mw - magically control your local media players

USAGE:
   mw [global options] command [command options] [arguments...]

COMMANDS:
   play, y            Instruct the player to play
   pause, u           Instruct the player to pause
   play-pause, p      Instruct the player to play or pause based on current state
   next, n            Instruct the player to play the next media
   previous, prev, v  Instruct the player to play the previous media
   stop, s            Instruct the player to stop
   open, o            Instruct the player to open the provided URI
   metadata           Get all available metadata about the current media
   daemon             Run the musicwand control daemon
   help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --player value  
   --help, -h      show help (default: false)
```
