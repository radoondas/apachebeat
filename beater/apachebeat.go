package beater

import (
	"net/url"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
)

const selector = "apachebeat"
const AUTO_STRING = "?auto"

// ApacheBeat implements Beater interface and sends Apache HTTPD status using libbeat.
type ApacheBeat struct {
	// ApConfig holds configurations of Apachebeat parsed by libbeat.
	urls   []*url.URL
	period time.Duration

	AbConfig ConfigSettings
	events   publisher.Client
	auth     bool
	username string
	password string

	done chan struct{}
}

func New() *ApacheBeat {
	return &ApacheBeat{
		done: make(chan struct{}),
	}
}

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
		urlConfig = []string{"http://127.0.0.1/server-status"}
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

	if ab.AbConfig.Input.Authentication.Username == nil || ab.AbConfig.Input.Authentication.Password == nil {
		logp.Err("Username or password is not set.")
		ab.auth = false
	} else if *ab.AbConfig.Input.Authentication.Username == "" || *ab.AbConfig.Input.Authentication.Password == "" {
		logp.Err("Username or password is not set.")
		ab.auth = false
	} else {
		ab.username = *ab.AbConfig.Input.Authentication.Username
		ab.password = *ab.AbConfig.Input.Authentication.Password
		ab.auth = true
		logp.Debug(selector, "Username %v", ab.username)
		logp.Debug(selector, "Password %v", ab.password)
	}

	logp.Debug(selector, "Init apachebeat")
	logp.Debug(selector, "Watch %v", ab.urls)
	logp.Debug(selector, "Period %v", ab.period)

	return nil
}

func (ab *ApacheBeat) Setup(b *beat.Beat) error {
	ab.events = b.Publisher.Connect()
	ab.done = make(chan struct{})
	return nil
}

func (ab *ApacheBeat) Run(b *beat.Beat) error {
	logp.Debug(selector, "Run apachebeat")

	//for each url
	for _, u := range ab.urls {
		go func(u *url.URL) {
			ticker := time.NewTicker(ab.period)
			defer ticker.Stop()

			for {
				select {
				case <-ab.done:
					goto GotoFinish
				case <-ticker.C:
				}

				timerStart := time.Now()

				//if eb.clusterStats {
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
					ab.events.PublishEvent(event)
				}

				timerEnd := time.Now()
				duration := timerEnd.Sub(timerStart)
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

func (ab *ApacheBeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (ab *ApacheBeat) Stop() {
	logp.Debug(selector, "Stop Apachebeat")
	close(ab.done)
}
