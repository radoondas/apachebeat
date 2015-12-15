package main

import (
  "net/url"
  "time"

  "github.com/elastic/beats/libbeat/beat"
  "github.com/elastic/beats/libbeat/cfgfile"
  "github.com/elastic/beats/libbeat/logp"

  "github.com/radoondas/apachebeat/collector"
  "github.com/radoondas/apachebeat/publisher"
)

const selector = "apachebeat"

// ApacheBeat implements Beater interface and sends Apache HTTPD status using libbeat.
type ApacheBeat struct {
  // ApConfig holds configurations of Apachebeat parsed by libbeat.
  AbConfig ConfigSettings

  done     chan uint

  urls     []*url.URL

  period   time.Duration
}

// func New() *Apachebeat {
// 	return &Apachebeat{}
// }

// Config ApacheBeat according to apachebeat.yml.
func (ab *ApacheBeat) Config(b *beat.Beat) error {
  //read config file
  err := cfgfile.Read(&ab.AbConfig, "")
  if err != nil {
    logp.Err("Error reading configuration file: %v", err)
    return err
  }

  //define default URL if none provided
  var urlConfig []string
  if ab.AbConfig.Input.URLs != nil {
    urlConfig = ab.AbConfig.Input.URLs
  } else {
    urlConfig = []string{"http://127.0.0.1/server-status?auto"}
  }

  ab.urls = make([]*url.URL, len(urlConfig))
  for i := 0; i < len(urlConfig); i++ {
    u, err := url.Parse(urlConfig[i])
    if err != nil {
      logp.Err("Invalid Apache HTTPD server status page: %v", err)
      return err
    }
    ab.urls[i] = u
  }

  if ab.AbConfig.Input.Period != nil {
    ab.period = time.Duration(*ab.AbConfig.Input.Period) * time.Second
  } else {
    ab.period = 10 * time.Second
  }

  logp.Debug(selector, "Init apachebeat")
  logp.Debug(selector, "Watch %v", ab.urls)
  logp.Debug(selector, "Period %v", ab.period)

  return nil
}

// Setup ApacheBeat.
func (ab *ApacheBeat) Setup(b *beat.Beat) error {
  ab.done = make(chan uint)

  return nil
}

// Run Apachebeat.
func (ab *ApacheBeat) Run(b *beat.Beat) error {
  logp.Debug(selector, "Run apachebeat")

  for _, u := range ab.urls {
    go func(u *url.URL) {
      var c collector.Collector
      var p publisher.Publisher

      c = collector.NewStubCollector()
      p = publisher.NewStubPublisher(b.Events)

      ticker := time.NewTicker(ab.period)
      defer ticker.Stop()

      for {
        select {
        case <-ab.done:
          goto GotoFinish
        case <-ticker.C:
        }

        start := time.Now()

        s, err := c.Collect(*u)

        if err != nil {
          logp.Err("Fail to read Apache HTTPD status: %v", err)
          goto GotoNext
        }
        p.Publish(s, u.String())

        GotoNext:
        end := time.Now()
        duration := end.Sub(start)
        if duration.Nanoseconds() > ab.period.Nanoseconds() {
          logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
        }
      }

      GotoFinish:
    }(u)
  }

  <-ab.done
  return nil
}

// Cleanup Apachebeat.
func (ab *ApacheBeat) Cleanup(b *beat.Beat) error {
  return nil
}

// Stop Apachebeat.
func (ab *ApacheBeat) Stop() {
  logp.Debug(selector, "Stop Apachebeat")
  close(ab.done)
}
