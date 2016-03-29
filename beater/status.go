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

func (ab *ApacheBeat) GetServerStatus(u url.URL) (common.MapStr, error) {
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
	// ServerUptimeSeconds: 43
	// Load1: 0.01
	// Load5: 0.10
	// Load15: 0.06
	// CPUUser: 0
	// CPUSystem: .01
	// CPUChildrenUser: 0
	// CPUChildrenSystem: 0

	var re *regexp.Regexp
	scanner := bufio.NewScanner(res.Body)

	// apachebeat - while we have something to read
	// default values for type int is 0
	var (
		hostname            string
		totalAccesses       int
		totalKBytes         int
		uptime              int
		cpuLoad             float64
		cpuUser             float64
		cpuSystem           float64
		cpuChildrenUser     float64
		cpuChildrenSystem   float64
		reqPerSec           float64
		bytesPerSec         float64
		bytesPerReq         float64
		busyWorkers         int
		idleWorkers         int
		connsTotal          int
		connsAsyncWriting   int
		connsAsyncKeepAlive int
		connsAsyncClosing   int
		serverUptimeSeconds int
		load1               float64
		load5               float64
		load15              float64
		totalS              int
		totalR              int
		totalW              int
		totalK              int
		totalD              int
		totalC              int
		totalL              int
		totalG              int
		totalI              int
		totalDot            int
		totalUnderscore     int
		totalAll            int
	)

	//set hostname from url
	hostname = u.Host
	logp.Debug(selector, "URL Hostname: %v", hostname)

	//for scanner.Scan() {
	//	// fmt.Println("read: ", scanner.Text())
	//	logp.Debug(selector, "%v: Reading from body: %v", hostname, scanner.Text())
	//
	//	// Total Accesses: 16147
	//	re = regexp.MustCompile("Total Accesses: (\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		total_access, _ = strconv.Atoi(matches[1])
	//		logp.Debug(selector, "%v: Total Accesses: %v", hostname, total_access)
	//	}
	//
	//	//Total kBytes: 12988
	//	re = regexp.MustCompile("Total kBytes: (\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		total_kbytes, _ = strconv.Atoi(matches[1])
	//		logp.Debug(selector, "%v: Total kBytes: %v", hostname, total_kbytes)
	//	}
	//
	//	// CPULoad: .000408393
	//	re = regexp.MustCompile("CPULoad: (\\d*.*\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		if strings.HasPrefix(matches[1], ".") {
	//			cpu_load = strings.Replace(matches[1], ".", "0.", 1)
	//		} else {
	//			cpu_load = matches[1]
	//		}
	//		logp.Debug(selector, "%v: CPULoad: %v", hostname, cpu_load)
	//	}
	//
	//	// CPUUser: 0
	//	re = regexp.MustCompile("CPUUser: (\\d*.*\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		if strings.HasPrefix(matches[1], ".") {
	//			cpu_user = strings.Replace(matches[1], ".", "0.", 1)
	//		} else {
	//			cpu_user = matches[1]
	//		}
	//		logp.Debug(selector, "%v: CPUUser: %v", hostname, cpu_user)
	//	}
	//
	//	// CPUSystem: .01
	//	re = regexp.MustCompile("CPUSystem: (\\d*.*\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		if strings.HasPrefix(matches[1], ".") {
	//			cpu_system = strings.Replace(matches[1], ".", "0.", 1)
	//		} else {
	//			cpu_system = matches[1]
	//		}
	//		logp.Debug(selector, "%v: CPUSystem: %v", hostname, cpu_system)
	//	}
	//
	//	// CPUChildrenUser: 0
	//	re = regexp.MustCompile("CPUChildrenUser: (\\d*.*\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		if strings.HasPrefix(matches[1], ".") {
	//			cpu_children_user = strings.Replace(matches[1], ".", "0.", 1)
	//		} else {
	//			cpu_children_user = matches[1]
	//		}
	//		logp.Debug(selector, "%v: CPUChildrenUser: %v", hostname, cpu_children_user)
	//	}
	//
	//	// CPUChildrenSystem: 0
	//	re = regexp.MustCompile("CPUChildrenSystem: (\\d*.*\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		if strings.HasPrefix(matches[1], ".") {
	//			cpu_children_system = strings.Replace(matches[1], ".", "0.", 1)
	//		} else {
	//			cpu_children_system = matches[1]
	//		}
	//		logp.Debug(selector, "%v: CPUChildrenSystem: %v", hostname, cpu_children_system)
	//	}
	//
	//	// Uptime: 3229728
	//	re = regexp.MustCompile("Uptime: (\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		uptime, _ = strconv.Atoi(matches[1])
	//		logp.Debug(selector, "%v: Uptime: %v", hostname, uptime)
	//	}
	//
	//	// ReqPerSec: .00499949
	//	re = regexp.MustCompile("ReqPerSec: (\\d*.*\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		if strings.HasPrefix(matches[1], ".") {
	//			req_per_sec = strings.Replace(matches[1], ".", "0.", 1)
	//		} else {
	//			req_per_sec = matches[1]
	//		}
	//		logp.Debug(selector, "%v: ReqPerSec: %v", hostname, req_per_sec)
	//	}
	//
	//	// BytesPerSec: 4.1179
	//	re = regexp.MustCompile("BytesPerSec: (\\d*.*\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		if strings.HasPrefix(matches[1], ".") {
	//			bytes_per_sec = strings.Replace(matches[1], ".", "0.", 1)
	//		} else {
	//			bytes_per_sec = matches[1]
	//		}
	//		logp.Debug(selector, "%v: BytesPerSec: %v", hostname, bytes_per_sec)
	//	}
	//
	//	// BytesPerReq: 823.665
	//	re = regexp.MustCompile("BytesPerReq: (\\d*.*\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		if strings.HasPrefix(matches[1], ".") {
	//			bytes_per_req = strings.Replace(matches[1], ".", "0.", 1)
	//		} else {
	//			bytes_per_req = matches[1]
	//		}
	//		logp.Debug(selector, "%v: BytesPerReq: %v", hostname, bytes_per_req)
	//	}
	//
	//	// BusyWorkers: 1
	//	re = regexp.MustCompile("BusyWorkers: (\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		busy_workers, _ = strconv.Atoi(matches[1])
	//		logp.Debug(selector, "%v: BusyWorkers: %v", hostname, busy_workers)
	//	}
	//
	//	// IdleWorkers: 8
	//	re = regexp.MustCompile("IdleWorkers: (\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		idle_workers, _ = strconv.Atoi(matches[1])
	//		logp.Debug(selector, "%v: IdleWorkers: %v", hostname, idle_workers)
	//	}
	//
	//	// ConnsTotal: 4940
	//	re = regexp.MustCompile("ConnsTotal: (\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		conns_total, _ = strconv.Atoi(matches[1])
	//		logp.Debug(selector, "%v: ConnsTotal: %v", hostname, conns_total)
	//	}
	//
	//	// ConnsAsyncWriting: 527
	//	re = regexp.MustCompile("ConnsAsyncWriting: (\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		conns_async_writing, _ = strconv.Atoi(matches[1])
	//		logp.Debug(selector, "%v: ConnsAsyncWriting: %v", hostname, conns_async_writing)
	//	}
	//
	//	// ConnsAsyncKeepAlive: 1321
	//	re = regexp.MustCompile("ConnsAsyncKeepAlive: (\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		conns_async_keep_alive, _ = strconv.Atoi(matches[1])
	//		logp.Debug(selector, "%v: ConnsAsyncKeepAlive: %v", hostname, conns_async_keep_alive)
	//	}
	//
	//	// ConnsAsyncClosing: 2785
	//	re = regexp.MustCompile("ConnsAsyncClosing: (\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		conns_async_closing, _ = strconv.Atoi(matches[1])
	//		logp.Debug(selector, "%v: ConnsAsyncClosing: %v", hostname, conns_async_closing)
	//	}
	//
	//	// ServerUptimeSeconds: 43
	//	re = regexp.MustCompile("ServerUptimeSeconds: (\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		server_uptime_seconds, _ = strconv.Atoi(matches[1])
	//		logp.Debug(selector, "%v: ServerUptimeSeconds: %v", hostname, server_uptime_seconds)
	//	}
	//
	//	//Load1: 0.01
	//	re = regexp.MustCompile("Load1: (\\d*.*\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		if strings.HasPrefix(matches[1], ".") {
	//			load1 = strings.Replace(matches[1], ".", "0.", 1)
	//		} else {
	//			load1 = matches[1]
	//		}
	//		logp.Debug(selector, "%v: Load1: %v", hostname, load1)
	//	}
	//
	//	//Load5: 0.10
	//	re = regexp.MustCompile("Load5: (\\d*.*\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		if strings.HasPrefix(matches[1], ".") {
	//			load5 = strings.Replace(matches[1], ".", "0.", 1)
	//		} else {
	//			load5 = matches[1]
	//		}
	//		logp.Debug(selector, "%v: Load5: %v", hostname, load5)
	//	}
	//
	//	//Load15: 0.06
	//	re = regexp.MustCompile("Load15: (\\d*.*\\d+)")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		if strings.HasPrefix(matches[1], ".") {
	//			load15 = strings.Replace(matches[1], ".", "0.", 1)
	//		} else {
	//			load15 = matches[1]
	//		}
	//		logp.Debug(selector, "%v: Load15: %v", hostname, load15)
	//	}
	//
	//	// Scoreboard Key:
	//	// "_" Waiting for Connection, "S" Starting up, "R" Reading Request,
	//	// "W" Sending Reply, "K" Keepalive (read), "D" DNS Lookup,
	//	// "C" Closing connection, "L" Logging, "G" Gracefully finishing,
	//	// "I" Idle cleanup of worker, "." Open slot with no current process
	//	// Scoreboard: _W____........___...............................................................................................................................................................................................................................................
	//	re = regexp.MustCompile("Scoreboard: (_|S|R|W|K|D|C|L|G|I|\\.)+")
	//	if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
	//		scr := strings.Split(scanner.Text(), " ")
	//
	//		tot_underscore = strings.Count(scr[1], "_")
	//		tot_s = strings.Count(scr[1], "S")
	//		tot_r = strings.Count(scr[1], "R")
	//		tot_w = strings.Count(scr[1], "W")
	//		tot_k = strings.Count(scr[1], "K")
	//		tot_d = strings.Count(scr[1], "D")
	//		tot_c = strings.Count(scr[1], "C")
	//		tot_l = strings.Count(scr[1], "L")
	//		tot_g = strings.Count(scr[1], "G")
	//		tot_i = strings.Count(scr[1], "I")
	//		tot_dot = strings.Count(scr[1], ".")
	//		tot_total = tot_underscore + tot_s + tot_r + tot_w + tot_k + tot_d + tot_c + tot_l + tot_g + tot_i + tot_dot
	//
	//		logp.Debug(selector, "%v: Waiting for Connection (_): %v", hostname, tot_underscore)
	//		logp.Debug(selector, "%v: Starting up (S): %v", hostname, tot_s)
	//		logp.Debug(selector, "%v: Reading Request (R): %v", hostname, tot_r)
	//		logp.Debug(selector, "%v: Sending Reply (W): %v", hostname, tot_w)
	//		logp.Debug(selector, "%v: Keepalive (read) (K): %v", hostname, tot_k)
	//		logp.Debug(selector, "%v: DNS Lookup (D): %v", hostname, tot_d)
	//		logp.Debug(selector, "%v: Closing connection (C): %v", hostname, tot_c)
	//		logp.Debug(selector, "%v: Logging (L): %v", hostname, tot_l)
	//		logp.Debug(selector, "%v: Gracefully finishing (G): %v", hostname, tot_g)
	//		logp.Debug(selector, "%v: Idle cleanup of worker (I): %v", hostname, tot_i)
	//		logp.Debug(selector, "%v: Open slot with no current process (.): %v", hostname, tot_dot)
	//	}
	//}

	for scanner.Scan() {
		logp.Debug(selector, "%v: Reading from body: %v", hostname, scanner.Text())

		// Total Accesses: 16147
		re = regexp.MustCompile("Total Accesses: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			totalAccesses, _ = strconv.Atoi(matches[1])
		}

		//Total kBytes: 12988
		re = regexp.MustCompile("Total kBytes: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			totalKBytes, _ = strconv.Atoi(matches[1])
		}

		// Uptime: 3229728
		re = regexp.MustCompile("Uptime: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			uptime, _ = strconv.Atoi(matches[1])
		}

		// CPULoad: .000408393
		re = regexp.MustCompile("CPULoad: (\\d*.*\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			cpuLoad = ParseMatchFloat(matches[1], hostname, "cpuLoad")
		}

		// CPUUser: 0
		re = regexp.MustCompile("CPUUser: (\\d*.*\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			cpuUser = ParseMatchFloat(matches[1], hostname, "cpuUser")
		}

		// CPUSystem: .01
		re = regexp.MustCompile("CPUSystem: (\\d*.*\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			cpuSystem = ParseMatchFloat(matches[1], hostname, "cpuSystem")
		}

		// CPUChildrenUser: 0
		re = regexp.MustCompile("CPUChildrenUser: (\\d*.*\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			cpuChildrenUser = ParseMatchFloat(matches[1], hostname, "cpuChildrenUser")
		}

		// CPUChildrenSystem: 0
		re = regexp.MustCompile("CPUChildrenSystem: (\\d*.*\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			cpuChildrenSystem = ParseMatchFloat(matches[1], hostname, "cpuChildrenSystem")
		}

		// ReqPerSec: .00499949
		re = regexp.MustCompile("ReqPerSec: (\\d*.*\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			reqPerSec = ParseMatchFloat(matches[1], hostname, "reqPerSec")
		}

		// BytesPerSec: 4.1179
		re = regexp.MustCompile("BytesPerSec: (\\d*.*\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			bytesPerSec = ParseMatchFloat(matches[1], hostname, "bytesPerSec")
		}

		// BytesPerReq: 823.665
		re = regexp.MustCompile("BytesPerReq: (\\d*.*\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			bytesPerReq = ParseMatchFloat(matches[1], hostname, "bytesPerReq")
		}

		// BusyWorkers: 1
		re = regexp.MustCompile("BusyWorkers: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			busyWorkers, _ = strconv.Atoi(matches[1])
		}

		// IdleWorkers: 8
		re = regexp.MustCompile("IdleWorkers: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			idleWorkers, _ = strconv.Atoi(matches[1])
		}

		// ConnsTotal: 4940
		re = regexp.MustCompile("ConnsTotal: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			connsTotal, _ = strconv.Atoi(matches[1])
		}

		// ConnsAsyncWriting: 527
		re = regexp.MustCompile("ConnsAsyncWriting: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			connsAsyncWriting, _ = strconv.Atoi(matches[1])
		}

		// ConnsAsyncKeepAlive: 1321
		re = regexp.MustCompile("ConnsAsyncKeepAlive: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			connsAsyncKeepAlive, _ = strconv.Atoi(matches[1])
		}

		// ConnsAsyncClosing: 2785
		re = regexp.MustCompile("ConnsAsyncClosing: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			connsAsyncClosing, _ = strconv.Atoi(matches[1])
		}

		// ServerUptimeSeconds: 43
		re = regexp.MustCompile("ServerUptimeSeconds: (\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			serverUptimeSeconds, _ = strconv.Atoi(matches[1])
		}

		//Load1: 0.01
		re = regexp.MustCompile("Load1: (\\d*.*\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			load1 = ParseMatchFloat(matches[1], hostname, "load1")
		}

		//Load5: 0.10
		re = regexp.MustCompile("Load5: (\\d*.*\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			load5 = ParseMatchFloat(matches[1], hostname, "load5")
		}

		//Load15: 0.06
		re = regexp.MustCompile("Load15: (\\d*.*\\d+)")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			load15 = ParseMatchFloat(matches[1], hostname, "load15")
		}

		// Scoreboard Key:
		// "_" Waiting for Connection, "S" Starting up, "R" Reading Request,
		// "W" Sending Reply, "K" Keepalive (read), "D" DNS Lookup,
		// "C" Closing connection, "L" Logging, "G" Gracefully finishing,
		// "I" Idle cleanup of worker, "." Open slot with no current process
		// Scoreboard: _W____........___...............................................................................................................................................................................................................................................
		re = regexp.MustCompile("Scoreboard: (_|S|R|W|K|D|C|L|G|I|\\.)+")
		if matches := re.FindStringSubmatch(scanner.Text()); matches != nil {
			scr := strings.Split(scanner.Text(), " ")

			totalUnderscore = strings.Count(scr[1], "_")
			totalS = strings.Count(scr[1], "S")
			totalR = strings.Count(scr[1], "R")
			totalW = strings.Count(scr[1], "W")
			totalK = strings.Count(scr[1], "K")
			totalD = strings.Count(scr[1], "D")
			totalC = strings.Count(scr[1], "C")
			totalL = strings.Count(scr[1], "L")
			totalG = strings.Count(scr[1], "G")
			totalI = strings.Count(scr[1], "I")
			totalDot = strings.Count(scr[1], ".")
			totalAll = totalUnderscore + totalS + totalR + totalW + totalK + totalD + totalC + totalL + totalG + totalI + totalDot
		}
	}

	event := common.MapStr{
		"hostname":      hostname,
		"totalAccesses": totalAccesses,
		"totalKBytes":   totalKBytes,
		"reqPerSec":     reqPerSec,
		"bytesPerSec":   bytesPerSec,
		"bytesPerReq":   bytesPerReq,
		"busyWorkers":   busyWorkers,
		"idleWorkers":   idleWorkers,
		"uptime": common.MapStr{
			"serverUptimeSeconds": serverUptimeSeconds,
			"uptime":              uptime,
		},
		"cpu": common.MapStr{
			"cpuLoad":           cpuLoad,
			"cpuUser":           cpuUser,
			"cpuSystem":         cpuSystem,
			"cpuChildrenUser":   cpuChildrenUser,
			"cpuChildrenSystem": cpuChildrenSystem,
		},
		"connections": common.MapStr{
			"connsTotal":          connsTotal,
			"connsAsyncWriting":   connsAsyncWriting,
			"connsAsyncKeepAlive": connsAsyncKeepAlive,
			"connsAsyncClosing":   connsAsyncClosing,
		},
		"load": common.MapStr{
			"load1":  load1,
			"load5":  load5,
			"load15": load15,
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

func ParseMatchFloat(inputString, hostname, fieldName string) float64 {
	var parseString string
	if strings.HasPrefix(inputString, ".") {
		parseString = strings.Replace(inputString, ".", "0.", 1)
	} else {
		parseString = inputString
	}
	outputFloat, er := strconv.ParseFloat(parseString, 64)

	/* Do we need to log failure? */
	if er != nil {
		logp.Warn("Host: %s - cannot parse string %s: %s to float.", hostname, fieldName, inputString)
		return 0.0
	}
	return outputFloat
}
