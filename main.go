package main

import (
	"github.com/elastic/beats/libbeat/beat"

	apachebeat "github.com/radoondas/apachebeat/beat"
	"os"
)

var Version = "1.0.0-beta2"
var Name = "apachebeat"

func main() {
	err := beat.Run(Name, Version, apachebeat.New())
	if err != nil {
		os.Exit(1)
	}
}
