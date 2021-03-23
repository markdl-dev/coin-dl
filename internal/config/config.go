package config

import (
	"flag"
	"time"

	"github.com/theckman/yacspin"
)

// Config is common for all the main command
type Config struct {
	ShowNotifications     bool
	PlayNotificationsBeep bool
	YakSpinConfig         yacspin.Config
}

func (c *Config) Setup(flags *flag.FlagSet) {
	flags.BoolVar(&c.ShowNotifications, "showNotifs", true, "show notification pop ups")
	flags.BoolVar(&c.PlayNotificationsBeep, "playNotifs", true, "plays a beep when a notification shows")

	c.YakSpinConfig = yacspin.Config{
		Frequency:         100 * time.Millisecond,
		CharSet:           yacspin.CharSets[10],
		Suffix:            " coindl",
		SuffixAutoColon:   true,
		Message:           "sit back and hodl.",
		StopCharacter:     "✓",
		StopColors:        []string{"fgGreen"},
		StopFailMessage:   "let's try again, the gecko might be out.",
		StopFailCharacter: "✗",
		StopFailColors:    []string{"fgRed"},
	}
}

func (c *Config) UpdateSpinner(suffix string, stopFailMessage string) {
	if len(suffix) > 0 {
		c.YakSpinConfig.Suffix = suffix
	}

	if len(stopFailMessage) > 0 {
		c.YakSpinConfig.StopFailMessage = stopFailMessage
	}
}
