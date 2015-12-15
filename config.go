package main

type ApacheConfig struct {
  // URLs to Apache status page.
  // Defaults to []string{"http://127.0.0.1/server-status?auto"}.
  URLs   []string

  // Period defines how often to read status in seconds.
  // Defaults to 1 second.
  Period *int64
}

type ConfigSettings struct {
  Input ApacheConfig
}
