package publisher

import (
	"time"

	"github.com/elastic/libbeat/common"
	"github.com/elastic/libbeat/publisher"
)

// StubPublisher is a Publisher that publishes Apache server-status.
type StubPublisher struct {
	client publisher.Client
}

// NewStubPublisher constructs a new StubPublisher.
func NewStubPublisher(c publisher.Client) *StubPublisher {
	return &StubPublisher{client: c}
}

// Publish Apache server-status.
func (p *StubPublisher) Publish(s map[string]interface{}, source string) {

	p.client.PublishEvent(common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       "apache_status",
		"source":     source,
		"apache":     s,
	})
}
