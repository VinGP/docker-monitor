package main

import (
	"pinger/internal/app"
	"pinger/internal/config"
)

func main() {
	app.Run(config.Get())
}
