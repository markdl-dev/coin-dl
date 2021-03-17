package main

import (
	"flag"
)

type RunConfig struct {
	showNotifications    bool
	playNotificationBeep bool
}

func main() {
	var runConfig RunConfig

	flag.BoolVar(&runConfig.showNotifications, "showNotifications", true, "show notifications")
	flag.BoolVar(&runConfig.playNotificationBeep, "playNotificationBeep", true, "plays a beep to notify you")
	flag.Parse()

}
