package musicwand

import (
	"fmt"
	"strings"
	"time"

	"github.com/shreve/musicwand/pkg/mpris"
)

func FormatStatus(template string, player *mpris.Player) string {
	findAndReplace(&template, "{artist}", func() string {
		return player.RawMetadata()["xesam:artist"].Value().([]string)[0]
	})

	findAndReplace(&template, "{album}", func() string {
		return player.RawMetadata()["xesam:album"].Value().([]string)[0]
	})

	findAndReplace(&template, "{track}", func() string {
		return player.RawMetadata()["xesam:title"].Value().(string)
	})

	findAndReplace(&template, "{length}", func() string {
		microseconds := player.RawMetadata()["mpris:length"].Value()
		switch microseconds.(type) {
		case int64:
			return formatTime(microseconds.(int64))
		case uint64:
			return formatTime(int64(microseconds.(uint64)))
		}
		return ""
	})

	findAndReplace(&template, "{position}", func() string {
		fmt.Println(player.Position())
		return formatTime(player.Position())
	})

	findAndReplace(&template, "{icon}", func() string {
		return string(Icon(player.Identity()))
	})

	return template
}

func formatTime(duration int64) string {
	length := time.Duration(duration) * time.Microsecond
	minutes := int(length.Minutes())
	length -= time.Duration(minutes) * time.Minute
	seconds := int(length.Seconds())
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func findAndReplace(source *string, query string, cb func() string) {
	if strings.Contains(*source, query) {
		*source = strings.ReplaceAll(*source, query, cb())
	}
}
