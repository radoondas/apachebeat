package collector

import (
  "net/url"
)

// Collector collects status from Apache HTTPD server-status page.
type Collector interface {
  // Collect status from the given url.
  Collect(u url.URL) (map[string]interface{}, error)
}
