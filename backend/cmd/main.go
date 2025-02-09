package main

import (
	"backend/internal/app"
	"backend/internal/config"
)

func main() {
	app.Run(config.Get())
}
