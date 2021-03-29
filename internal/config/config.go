package config

import (
	"flag"
	"time"

	"github.com/theckman/yacspin"
)

const (
	showNotifsFlagValue           = true
	playNotifsFlagValue           = true
	showNotifsFlagDescription     = "show notification pop ups"
	showPlayNotifsFlagDescription = "plays a beep when a notification shows"
)

// Config is common for all the main command
type Config struct {
	ShowNotifications     bool
	PlayNotificationsBeep bool
	YakSpinConfig         yacspin.Config
}

func (c *Config) Setup(flagSet *flag.FlagSet) {
	flagSet.BoolVar(&c.ShowNotifications, "showNotifs", showNotifsFlagValue, showNotifsFlagDescription)
	flagSet.BoolVar(&c.PlayNotificationsBeep, "playNotifs", playNotifsFlagValue, showPlayNotifsFlagDescription)
	flagSet.BoolVar(&c.ShowNotifications, "sn", showNotifsFlagValue, showNotifsFlagDescription+" (shorthand)")
	flagSet.BoolVar(&c.PlayNotificationsBeep, "pn", playNotifsFlagValue, showPlayNotifsFlagDescription+" (shorthand)")

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
