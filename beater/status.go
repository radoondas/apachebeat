package beater

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

var (
	scoreboardRegexp = regexp.MustCompile("(Scoreboard):\\s+((_|S|R|W|K|D|C|L|G|I|\\.)+)")

	// This should match: "CPUSystem: .01"
	matchNumber = regexp.MustCompile("(^[0-9a-zA-Z ]+):\\s+(\\d*\\.?\\d+)")
)

func (ab *ApacheBeat) GetServerStatus(u url.URL) (common.MapStr, error) {
	var (
		totalS          int
		totalR          int
		totalW          int
		totalK          int
		totalD          int
		totalC          int
		totalL          int
		totalG          int
		totalI          int
		totalDot        int
		totalUnderscore int
		totalAll        int
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String()+AUTO_STRING, nil)

	if ab.auth {
		req.SetBasicAuth(ab.username, ab.password)
	}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP%s", res.Status)
	}

	//set hostname from url
	hostname := u.Host
	logp.Debug(selector, "URL Hostname: %v", hostname)

	fullEvent := common.MapStr{}
	scanner := bufio.NewScanner(res.Body)

	// Iterate through all events to gather data
	for scanner.Scan() {
		if match := matchNumber.FindStringSubmatch(scanner.Text()); len(match) == 3 {
			// Total Accesses: 16147
			//Total kBytes: 12988
			// Uptime: 3229728
			// CPULoad: .000408393
			// CPUUser: 0
			// CPUSystem: .01
			// CPUChildrenUser: 0
			// CPUChildrenSystem: 0
			// ReqPerSec: .00499949
			// BytesPerSec: 4.1179
			// BytesPerReq: 823.665
			// BusyWorkers: 1
			// IdleWorkers: 8
			// ConnsTotal: 4940
			// ConnsAsyncWriting: 527
			// ConnsAsyncKeepAlive: 1321
			// ConnsAsyncClosing: 2785
			// ServerUptimeSeconds: 43
			//Load1: 0.01
			//Load5: 0.10
			//Load15: 0.06
			fullEvent[match[1]] = match[2]

		} else if match := scoreboardRegexp.FindStringSubmatch(scanner.Text()); len(match) == 4 {
			// Scoreboard Key:
			// "_" Waiting for Connection, "S" Starting up, "R" Reading Request,
			// "W" Sending Reply, "K" Keepalive (read), "D" DNS Lookup,
			// "C" Closing connection, "L" Logging, "G" Gracefully finishing,
			// "I" Idle cleanup of worker, "." Open slot with no current process
			// Scoreboard: _W____........___...............................................................................................................................................................................................................................................

			totalUnderscore = strings.Count(match[2], "_")
			totalS = strings.Count(match[2], "S")
			totalR = strings.Count(match[2], "R")
			totalW = strings.Count(match[2], "W")
			totalK = strings.Count(match[2], "K")
			totalD = strings.Count(match[2], "D")
			totalC = strings.Count(match[2], "C")
			totalL = strings.Count(match[2], "L")
			totalG = strings.Count(match[2], "G")
			totalI = strings.Count(match[2], "I")
			totalDot = strings.Count(match[2], ".")
			totalAll = totalUnderscore + totalS + totalR + totalW + totalK + totalD + totalC + totalL + totalG + totalI + totalDot

		} else {

			logp.Debug("Unexpected line in apache server-status output: %s", scanner.Text())
		}
	}

	event := common.MapStr{
		"hostname":      hostname,
		"totalAccesses": toInt(fullEvent["Total Accesses"]),
		"totalKBytes":   toInt(fullEvent["Total kBytes"]),
		"reqPerSec":     parseMatchFloat(fullEvent["ReqPerSec"], hostname, "ReqPerSec"),
		"bytesPerSec":   parseMatchFloat(fullEvent["BytesPerSec"], hostname, "BytesPerSec"),
		"bytesPerReq":   parseMatchFloat(fullEvent["BytesPerReq"], hostname, "BytesPerReq"),
		"busyWorkers":   toInt(fullEvent["BusyWorkers"]),
		"idleWorkers":   toInt(fullEvent["IdleWorkers"]),
		"uptime": common.MapStr{
			"serverUptimeSeconds": toInt(fullEvent["ServerUptimeSeconds"]),
			"uptime":              toInt(fullEvent["Uptime"]),
		},
		"cpu": common.MapStr{
			"cpuLoad":           parseMatchFloat(fullEvent["CPULoad"], hostname, "CPULoad"),
			"cpuUser":           parseMatchFloat(fullEvent["CPUUser"], hostname, "CPUUser"),
			"cpuSystem":         parseMatchFloat(fullEvent["CPUSystem"], hostname, "CPUSystem"),
			"cpuChildrenUser":   parseMatchFloat(fullEvent["CPUChildrenUser"], hostname, "CPUChildrenUser"),
			"cpuChildrenSystem": parseMatchFloat(fullEvent["CPUChildrenSystem"], hostname, "CPUChildrenSystem"),
		},
		"connections": common.MapStr{
			"connsTotal":          toInt(fullEvent["ConnsTotal"]),
			"connsAsyncWriting":   toInt(fullEvent["ConnsAsyncWriting"]),
			"connsAsyncKeepAlive": toInt(fullEvent["ConnsAsyncKeepAlive"]),
			"connsAsyncClosing":   toInt(fullEvent["ConnsAsyncClosing"]),
		},
		"load": common.MapStr{
			"load1":  parseMatchFloat(fullEvent["Load1"], hostname, "Load1"),
			"load5":  parseMatchFloat(fullEvent["Load5"], hostname, "Load5"),
			"load15": parseMatchFloat(fullEvent["Load15"], hostname, "Load15"),
		},
		"scoreboard": common.MapStr{
			"startingUp":           totalS,
			"readingRequest":       totalR,
			"sendingReply":         totalW,
			"keepalive":            totalK,
			"dnsLookup":            totalD,
			"closingConnection":    totalC,
			"logging":              totalL,
			"gracefullyFinishing":  totalG,
			"idleCleanup":          totalI,
			"openSlot":             totalDot,
			"waitingForConnection": totalUnderscore,
			"total":                totalAll,
		},
	}

	return event, nil
}

func parseMatchFloat(input interface{}, hostname, fieldName string) float64 {
	var parseString string

	if input != nil {
		if strings.HasPrefix(input.(string), ".") {
			parseString = strings.Replace(input.(string), ".", "0.", 1)
		} else {
			parseString = input.(string)
		}
		outputFloat, er := strconv.ParseFloat(parseString, 64)

		/* Do we need to log failure? */
		if er != nil {
			logp.Debug("Host: %s - cannot parse string %s: %s to float.", hostname, fieldName, input)
			return 0.0
		}
		return outputFloat
	} else {
		return 0.0
	}
}

// toInt converts value to int. In case of error, returns 0
func toInt(param interface{}) int {
	if param == nil {
		return 0
	}

	value, err := strconv.Atoi(param.(string))

	if err != nil {
		logp.Err("Error converting param to int: %s", param)
		value = 0
	}

	return value
}
