package collector

import (
	"net/url"
)

// Collector collects status from Apache HTTPD server-status module.
type Collector interface {
	// Collect status from the given url.
	Collect(u url.URL) (map[string]interface{}, error)
}
