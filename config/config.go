package config

import "time"

type Config struct {
	// URLs to Apache status page.
	// Defaults to []string{"http://127.0.0.1/server-status?auto"}.
	URLs []string `config:"urls"`

	// Period defines how often to read status in seconds.
	// Defaults to 1 second.
	Period time.Duration `config:"period"`

	// Authentication for BasicAuth
	Authentication Authentication `config:"authentication"`
}

type Authentication struct {
	Username string
	Password string
}

var (
	DefaultConfig = Config{
		Period: 10 * time.Second,
		URLs:   []string{"http://127.0.0.1:8080"},
		Authentication: Authentication{
			Username: "",
			Password: "",
		},
	}
)
