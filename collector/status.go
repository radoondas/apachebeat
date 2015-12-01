package collector

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/elastic/libbeat/logp"
)

const selector = "apachebeat"

// StubCollector is a Collector that collects Apache HTTPD server-status page.
type StubCollector struct {
	requests int
}

// NewStubCollector constructs a new StubCollector.
func NewStubCollector() Collector {
	return &StubCollector{requests: 0}
}

// Collect Apache HTTPD server-status from given url.
func (c *StubCollector) Collect(u url.URL) (map[string]interface{}, error) {
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP%s", res.Status)
	}

	// Apache status example:
	// Total Accesses: 16147
	// Total kBytes: 12988i
	// CPULoad: .000408393
	// Uptime: 3229728
	// ReqPerSec: .00499949
	// BytesPerSec: 4.1179
	// BytesPerReq: 823.665
	// BusyWorkers: 1
	// IdleWorkers: 8
	// ConnsTotal: 4940
	// ConnsAsyncWriting: 527
	// ConnsAsyncKeepAlive: 1321
	// ConnsAsyncClosing: 2785

	var re *regexp.Regexp
	scanner := bufio.NewScanner(res.Body)

	// apachebeat - while we have something to read
	// default values for type int is 0
	var (
		total_access           int
		total_kbytes           int
		cpu_load               float64
		uptime                 int
		req_per_sec            int
		bytes_per_sec          int
		bytes_per_req          int
		busy_workers           int
		idle_workers           int
		conns_total            int
		conns_async_writing    int
		conns_async_keep_alive int
		conns_async_closing    int
	)

	//cpuload - string to float!!

	var (
		hostname string
	)

	//set hostname from url
	hostname = u.Host
	logp.Debug(selector, "URL Hostname: %v", hostname)

	var (
		tot_s          int
		tot_r          int
		tot_w          int
		tot_k          int
		tot_d          int
		tot_c          int
		tot_l          int
		tot_g          int
		tot_i          int
		tot_dot        int
		tot_underscore int
	)

	for scanner.Scan() {
		// fmt.Println("read: ", scanner.Text())
		logp.Debug(selector, "Reading from body: %v", scanner.Text())

		// Total Accesses: 16147
		re = regexp.MustCompile("Total Accesses: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
			//
		} else {
			total_access, _ = strconv.Atoi(matches[1])
			logp.Debug(selector, "Total Accesses: %v", total_access)
		}

		//Total kBytes: 12988
		re = regexp.MustCompile("Total kBytes: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
			//
		} else {
			total_kbytes, _ = strconv.Atoi(matches[1])
			logp.Debug(selector, "Total kBytes: %v", total_kbytes)
		}

		// CPULoad: .000408393
		re = regexp.MustCompile("CPULoad: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
			//
		} else {
			//TODO
			// cpu_load, _ = strconv.Atoi(matches[1])
			cpu_load, _ = strconv.ParseFloat(matches[1], 64)
			logp.Debug(selector, "CPULoad: %v", cpu_load)
		}

		// Uptime: 3229728
		re = regexp.MustCompile("Uptime: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
			//
		} else {
			uptime, _ = strconv.Atoi(matches[1])
			logp.Debug(selector, "Uptime: %v", uptime)
		}

		// ReqPerSec: .00499949
		re = regexp.MustCompile("ReqPerSec: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
			//
		} else {
			req_per_sec, _ = strconv.Atoi(matches[1])
			logp.Debug(selector, "ReqPerSec: %v", req_per_sec)
		}

		// BytesPerSec: 4.1179
		re = regexp.MustCompile("BytesPerSec: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
			//
		} else {
			bytes_per_sec, _ = strconv.Atoi(matches[1])
			logp.Debug(selector, "BytesPerSec: %v", bytes_per_sec)
		}

		// BytesPerReq: 823.665
		re = regexp.MustCompile("BytesPerReq: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
			//
		} else {
			bytes_per_req, _ = strconv.Atoi(matches[1])
			logp.Debug(selector, "BytesPerReq: %v", bytes_per_req)
		}

		// BusyWorkers: 1
		re = regexp.MustCompile("BusyWorkers: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
			//
		} else {
			busy_workers, _ = strconv.Atoi(matches[1])
			logp.Debug(selector, "BusyWorkers: %v", busy_workers)
		}

		// IdleWorkers: 8
		re = regexp.MustCompile("IdleWorkers: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
			//
		} else {
			idle_workers, _ = strconv.Atoi(matches[1])
			logp.Debug(selector, "IdleWorkers: %v", idle_workers)
		}

		// ConnsTotal: 4940
		re = regexp.MustCompile("ConnsTotal: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
			//
		} else {
			conns_total, _ = strconv.Atoi(matches[1])
			logp.Debug(selector, "ConnsTotal: %v", conns_total)
		}

		// ConnsAsyncWriting: 527
		re = regexp.MustCompile("ConnsAsyncWriting: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
			//
		} else {
			conns_async_writing, _ = strconv.Atoi(matches[1])
			logp.Debug(selector, "ConnsAsyncWriting: %v", conns_async_writing)
		}

		// ConnsAsyncKeepAlive: 1321
		re = regexp.MustCompile("ConnsAsyncKeepAlive: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
			//
		} else {
			conns_async_keep_alive, _ = strconv.Atoi(matches[1])
			logp.Debug(selector, "ConnsAsyncKeepAlive: %v", conns_async_keep_alive)
		}

		// ConnsAsyncClosing: 2785
		re = regexp.MustCompile("ConnsAsyncClosing: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
			//
		} else {
			conns_async_closing, _ = strconv.Atoi(matches[1])
			logp.Debug(selector, "ConnsAsyncClosing: %v", conns_async_closing)
		}

		// Scoreboard Key:
		// "_" Waiting for Connection, "S" Starting up, "R" Reading Request,
		// "W" Sending Reply, "K" Keepalive (read), "D" DNS Lookup,
		// "C" Closing connection, "L" Logging, "G" Gracefully finishing,
		// "I" Idle cleanup of worker, "." Open slot with no current process
		// Scoreboard: _W____........___...............................................................................................................................................................................................................................................
		re = regexp.MustCompile("Scoreboard: (\\w+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches == nil {
			//
		} else {
			scr := strings.Split(scanner.Text(), " ")

			tot_underscore = strings.Count(scr[1], "_")
			tot_s = strings.Count(scr[1], "S")
			tot_r = strings.Count(scr[1], "R")
			tot_w = strings.Count(scr[1], "W")
			tot_k = strings.Count(scr[1], "K")
			tot_d = strings.Count(scr[1], "D")
			tot_c = strings.Count(scr[1], "C")
			tot_l = strings.Count(scr[1], "L")
			tot_g = strings.Count(scr[1], "G")
			tot_i = strings.Count(scr[1], "I")
			tot_dot = strings.Count(scr[1], ".")

			logp.Debug(selector, "Waiting for Connection (_): %v", tot_underscore)
			logp.Debug(selector, "Starting up (S): %v", tot_s)
			logp.Debug(selector, "Reading Request (R): %v", tot_r)
			logp.Debug(selector, "Sending Reply (W): %v", tot_w)
			logp.Debug(selector, "Keepalive (read) (K): %v", tot_k)
			logp.Debug(selector, "DNS Lookup (D): %v", tot_d)
			logp.Debug(selector, "Closing connection (C): %v", tot_c)
			logp.Debug(selector, "Logging (L): %v", tot_l)
			logp.Debug(selector, "Gracefully finishing (G): %v", tot_g)
			logp.Debug(selector, "Idle cleanup of worker (I): %v", tot_i)
			logp.Debug(selector, "Open slot with no current process (.): %v", tot_dot)
		}
	}

	return map[string]interface{}{
		"total_access":           total_access,
		"total_kbytes":           total_kbytes,
		"cpu_load":               cpu_load,
		"uptime":                 uptime,
		"req_per_sec":            req_per_sec,
		"bytes_per_sec":          bytes_per_sec,
		"bytes_per_req":          bytes_per_req,
		"busy_workers":           busy_workers,
		"idle_workers":           idle_workers,
		"conns_total":            conns_total,
		"conns_async_writing":    conns_async_writing,
		"conns_async_keep_alive": conns_async_keep_alive,
		"conns_async_closing":    conns_async_closing,
		"host_url":               hostname,
		"tot_s":                  tot_s,
		"tot_r":                  tot_r,
		"tot_w":                  tot_w,
		"tot_k":                  tot_k,
		"tot_d":                  tot_d,
		"tot_c":                  tot_c,
		"tot_l":                  tot_l,
		"tot_g":                  tot_g,
		"tot_i":                  tot_i,
		"tot_dot":                tot_dot,
		"tot_underscore":         tot_underscore,
	}, nil
}
