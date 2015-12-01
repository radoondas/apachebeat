package main

type ApacheConfig struct {
	// URLs to Apache status page.
	// Defaults to []string{"http://127.0.0.1/server-status?auto"}.
	URLs []string
	// Format of the status page, either "stub" or "plus".
	// Use "stub" for Nginx stub status page.
	// # Use "plus" for Nginx Plus status JSON document.
	// Defaults to "stub".
	//Format string

	// Period defines how often to read status in seconds.
	// Defaults to 1 second.
	Period *int64
}

type ConfigSettings struct {
	Input ApacheConfig
}
