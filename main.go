package main

import (
	"github.com/elastic/beats/libbeat/beat"

	apachebeat "github.com/radoondas/apachebeat/beat"
)

var Version = "1.0.0-beta2"
var Name = "apachebeat"

func main() {

	beat.Run(Name, Version, apachebeat.New())
}
