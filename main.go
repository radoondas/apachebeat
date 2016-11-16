package main

import (
	"github.com/elastic/beats/libbeat/beat"

	"github.com/radoondas/apachebeat/beater"
	"os"
)

var Name = "apachebeat"

func main() {
	if err := beat.Run(Name, "", beater.New); err != nil {
		os.Exit(1)
	}
}
