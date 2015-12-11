package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/logp"
)

var Version = "1.0.0-beta2"
var Name = "apachebeat"

func main() {

	ab := &ApacheBeat{}

	beat := beat.NewBeat(Name, Version, ab)

	beat.CommandLineSetup()

	beat.LoadConfig()
	err := ab.Config(beat)
	if err != nil {
		logp.Critical("Config error: %v", err)
		os.Exit(1)
	}

	beat.Run()
}
