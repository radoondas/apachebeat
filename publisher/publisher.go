package publisher

// Publisher publishes Apache HTTPD server-status via libbeat.
type Publisher interface {
	Publish(s map[string]interface{}, source string)
}
