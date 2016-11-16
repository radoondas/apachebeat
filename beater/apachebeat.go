package beater

import (
	"fmt"
	"net/url"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	cfg "github.com/radoondas/apachebeat/config"
)

const selector = "apachebeat"
const AUTO_STRING = "?auto"

// ApacheBeat implements Beater interface and sends Apache HTTPD status using libbeat.
type ApacheBeat struct {
	config cfg.Config
	urls   []*url.URL
	auth   bool
	client publisher.Client
	done   chan struct{}
}

// Creates beater
func New(b *beat.Beat, rawCfg *common.Config) (beat.Beater, error) {
	config := cfg.DefaultConfig
	if err := rawCfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &ApacheBeat{
		done:   make(chan struct{}),
		config: config,
		auth:   true,
	}

	err := bt.init(b)
	if err != nil {
		return nil, err
	}

	return bt, nil
}

/// *** Beater interface methods ***///
func (ab *ApacheBeat) init(b *beat.Beat) error {

	ab.urls = make([]*url.URL, len(ab.config.URLs))
	for i := 0; i < len(ab.config.URLs); i++ {
		u, err := url.Parse(ab.config.URLs[i])
		if err != nil {
			logp.Err("Invalid Apache HTTPD server status url: %v", err)
			return err
		}
		ab.urls[i] = u
	}

	//Disable authentication when no username or password is set
	if ab.config.Authentication.Username == "" || ab.config.Authentication.Password == "" {
		logp.Info("One of username or password IS NOT set.")
		ab.auth = false
	}

	return nil
}

func (ab *ApacheBeat) Run(b *beat.Beat) error {
	logp.Debug(selector, "Run apachebeat")

	ab.client = b.Publisher.Connect()
	//for each url
	for _, u := range ab.urls {
		go func(u *url.URL) {
			ticker := time.NewTicker(ab.config.Period)
			defer ticker.Stop()

			for {
				select {
				case <-ab.done:
					goto GotoFinish
				case <-ticker.C:
				}

				timerStart := time.Now()

				logp.Debug(selector, "Cluster stats for url: %v", u)
				serverStatus, error := ab.GetServerStatus(*u)
				if error != nil {
					logp.Err("Error getting server-status for %s: %v", u.String(), error)
				} else {
					logp.Debug(selector, "Apache  detailfor %s : %+v", u.String(), serverStatus)

					event := common.MapStr{
						"@timestamp": common.Time(time.Now()),
						"type":       "apache_status", //TODO: NAMING??
						"url":        u.String(),
						"apache":     serverStatus, //TODO: NAMING??
					}

					logp.Debug(selector, "Server status event detail for %s: %v", u.String(), event)
					ab.client.PublishEvent(event)
				}

				timerEnd := time.Now()
				duration := timerEnd.Sub(timerStart)
				if duration.Nanoseconds() > ab.config.Period.Nanoseconds() {
					logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
				}
			}

		GotoFinish:
		}(u)
	}

	<-ab.done
	return nil
}

func (ab *ApacheBeat) Stop() {
	logp.Info("Stopping Apachebeat")
	if ab.done != nil {
		ab.client.Close()
		close(ab.done)
	}
}
