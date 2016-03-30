package beater

import (
	"testing"

	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/stretchr/testify/assert"
	"github.com/elastic/beats/libbeat/common"
)

// Apache status example:
// Total Accesses: 16147
// Total kBytes: 12988
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

func TestStubCollector(t *testing.T) {
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Total Accesses: 16147")
		fmt.Fprintln(w, "Total kBytes: 12988")
		fmt.Fprintln(w, "CPULoad: .000408393")
		fmt.Fprintln(w, "Uptime: 3229728")
		fmt.Fprintln(w, "ReqPerSec: .00499949")
		fmt.Fprintln(w, "BytesPerSec: 4.1179")
		fmt.Fprintln(w, "BytesPerReq: 823.665")
		fmt.Fprintln(w, "BusyWorkers: 1")
		fmt.Fprintln(w, "IdleWorkers: 8")
		fmt.Fprintln(w, "ConnsTotal: 4940")
		fmt.Fprintln(w, "ConnsAsyncWriting: 527")
		fmt.Fprintln(w, "ConnsAsyncKeepAlive: 1321")
		fmt.Fprintln(w, "ConnsAsyncClosing: 2785")
		fmt.Fprintln(w, "ServerUptimeSeconds: 43")
		fmt.Fprintln(w, "Load1: 0.01")
		fmt.Fprintln(w, "Load5: 0.10")
		fmt.Fprintln(w, "Load15: 0.06")
		fmt.Fprintln(w, "CPUUser: 0")
		fmt.Fprintln(w, "CPUSystem: .01")
		fmt.Fprintln(w, "CPUChildrenUser: 0")
		fmt.Fprintln(w, "CPUChildrenSystem: 0")
		fmt.Fprintln(w, "Scoreboard: _W____........___...............................................................................................................................................................................................................................................")
	}))
	defer ts1.Close()

	ab1 := ApacheBeat{}
	u1, _ := url.Parse(ts1.URL)
	s1, _ := ab1.GetServerStatus(*u1)

	assert.Equal(t, s1["totalAccesses"].(int), 16147)
	assert.Equal(t, s1["totalKBytes"].(int), 12988)
	assert.Equal(t, s1["reqPerSec"].(float64), 0.00499949)
	assert.Equal(t, s1["bytesPerSec"].(float64), 4.1179)
	assert.Equal(t, s1["bytesPerReq"].(float64), 823.665)
	assert.Equal(t, s1["busyWorkers"].(int), 1)
	assert.Equal(t, s1["idleWorkers"].(int), 8)

	up1 := s1["uptime"].(common.MapStr)
	assert.Equal(t, up1["uptime"].(int), 3229728)
	assert.Equal(t, up1["serverUptimeSeconds"].(int), 43)

	c1 := s1["connections"].(common.MapStr)
	assert.Equal(t, c1["connsTotal"].(int), 4940)
	assert.Equal(t, c1["connsAsyncWriting"].(int), 527)
	assert.Equal(t, c1["connsAsyncKeepAlive"].(int), 1321)
	assert.Equal(t, c1["connsAsyncClosing"].(int), 2785)
	//assert.Equal(t, s1["host_url"], "localhost")

	l1 := s1["load"].(common.MapStr)
	assert.Equal(t, l1["load1"].(float64), 0.01)
	assert.Equal(t, l1["load5"].(float64), 0.10)
	assert.Equal(t, l1["load15"].(float64), 0.06)

	cpu1 := s1["cpu"].(common.MapStr)
	assert.Equal(t, cpu1["cpuLoad"].(float64), 0.000408393)
	assert.Equal(t, cpu1["cpuUser"].(float64), float64(0))
	assert.Equal(t, cpu1["cpuSystem"].(float64), 0.01)
	assert.Equal(t, cpu1["cpuChildrenUser"].(float64), float64(0))
	assert.Equal(t, cpu1["cpuChildrenSystem"].(float64), float64(0))

	scb1 := s1["scoreboard"].(common.MapStr)
	assert.Equal(t, scb1["startingUp"].(int), 0)
	assert.Equal(t, scb1["readingRequest"].(int), 0)
	assert.Equal(t, scb1["sendingReply"].(int), 1)
	assert.Equal(t, scb1["keepalive"].(int), 0)
	assert.Equal(t, scb1["dnsLookup"].(int), 0)
	assert.Equal(t, scb1["closingConnection"].(int), 0)
	assert.Equal(t, scb1["logging"].(int), 0)
	assert.Equal(t, scb1["gracefullyFinishing"].(int), 0)
	assert.Equal(t, scb1["idleCleanup"].(int), 0)
	assert.Equal(t, scb1["openSlot"].(int), 247)
	assert.Equal(t, scb1["waitingForConnection"].(int), 8)
	assert.Equal(t, scb1["total"].(int), 256)

	//  BusyWorkers: 1
	//  IdleWorkers: 4
	//  Scoreboard: W____...........................................................................................................................................................................................................................................................

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "BusyWorkers: 3")
		fmt.Fprintln(w, "IdleWorkers: 15")
		fmt.Fprintln(w, "Scoreboard: W____...........................................................................................................................................................................................................................................................")
	}))
	defer ts2.Close()

	ab2 := ApacheBeat{}
	u2, _ := url.Parse(ts2.URL)
	s2, _ := ab2.GetServerStatus(*u2)

	assert.Equal(t, s2["totalAccesses"].(int), 0)
	assert.Equal(t, s2["totalKBytes"].(int), 0)
	assert.Equal(t, s2["reqPerSec"].(float64), float64(0))
	assert.Equal(t, s2["bytesPerSec"].(float64), float64(0))
	assert.Equal(t, s2["bytesPerReq"].(float64), float64(0))
	assert.Equal(t, s2["busyWorkers"].(int), 3)
	assert.Equal(t, s2["idleWorkers"].(int), 15)

	up2 := s2["uptime"].(common.MapStr)
	assert.Equal(t, up2["uptime"].(int), 0)
	assert.Equal(t, up2["serverUptimeSeconds"].(int), 0)

	c2 := s2["connections"].(common.MapStr)
	assert.Equal(t, c2["connsTotal"].(int), 0)
	assert.Equal(t, c2["connsAsyncWriting"].(int), 0)
	assert.Equal(t, c2["connsAsyncKeepAlive"].(int), 0)
	assert.Equal(t, c2["connsAsyncClosing"].(int), 0)
	//assert.Equal(t, s2["host_url"], "localhost")

	l2 := s2["load"].(common.MapStr)
	assert.Equal(t, l2["load1"].(float64), float64(0))
	assert.Equal(t, l2["load5"].(float64), float64(0))
	assert.Equal(t, l2["load15"].(float64), float64(0))

	cpu2 := s2["cpu"].(common.MapStr)
	assert.Equal(t, cpu2["cpuLoad"].(float64), float64(0))
	assert.Equal(t, cpu2["cpuUser"].(float64), float64(0))
	assert.Equal(t, cpu2["cpuSystem"].(float64), float64(0))
	assert.Equal(t, cpu2["cpuChildrenUser"].(float64), float64(0))
	assert.Equal(t, cpu2["cpuChildrenSystem"].(float64), float64(0))

	scb2 := s2["scoreboard"].(common.MapStr)
	assert.Equal(t, scb2["startingUp"].(int), 0)
	assert.Equal(t, scb2["readingRequest"].(int), 0)
	assert.Equal(t, scb2["sendingReply"].(int), 1)
	assert.Equal(t, scb2["keepalive"].(int), 0)
	assert.Equal(t, scb2["dnsLookup"].(int), 0)
	assert.Equal(t, scb2["closingConnection"].(int), 0)
	assert.Equal(t, scb2["logging"].(int), 0)
	assert.Equal(t, scb2["gracefullyFinishing"].(int), 0)
	assert.Equal(t, scb2["idleCleanup"].(int), 0)
	assert.Equal(t, scb2["openSlot"].(int), 251)
	assert.Equal(t, scb2["waitingForConnection"].(int), 4)
	assert.Equal(t, scb2["total"].(int), 256)

	//  Total Accesses: 491803362
	//  Total kBytes: 21453176370
	//  CPULoad: .84122
	//  Uptime: 3761632
	//  ReqPerSec: 130.742
	//  BytesPerSec: 5840030
	//  BytesPerReq: 44668.4
	//  BusyWorkers: 1298
	//  IdleWorkers: 652
	//  ConnsTotal: 3264
	//  ConnsAsyncWriting: 82
	//  ConnsAsyncKeepAlive: 1574
	//  ConnsAsyncClosing: 309
	//  ServerUptimeSeconds: 5632
	//  Load1: 1.01
	//  Load5: 1.50
	//  Load15: 2.06
	//  CPUUser: 0.5
	//  CPUSystem: .01
	//  CPUChildrenUser: .1
	//  CPUChildrenSystem: 1.99
	//  Scoreboard: RRR_R__RR_RLRR_RRRRRRRRRR_RRR_RRRRR__R___RRR_R_RRWRR_R______RWRR___RR_RR_RRRR_WR_R__W_RR__RRRR_____RR___RRRWRR__RRR_R_RRR_RW_RR_RRR_WRRRR__R__RRRR_RRRW_RRRRR_RRRRRR_RRRR_RRRWRWRRRRRR_R_R____R___RRR_W__RRRRRR_RR_W__R_RRRRRRR_RR_RRR____W_RRRR_RRRW_RRRRRR_RR_RR___R_RRRRRRRR_RR__RRRR__RRRRR_RWR__RRWRRRR__RR_RRWWR_R_RWRR__R_RR_RRRRRWR___RW_RR__RR_R_RWRRR_RRRRRRW___R_RRRR_RRR__RRRRWR___R_R__RRR__RRRW_RRRR__RRWR_RRR__L_R_R__R_R_R_R__WRRRRRRRR__RWRRRW_RRWRRRWWL__RR_R__R_RRR____RRRRRR_RWRRRR_R_RRR_R__R_RR_R_RRRRR_RRRW_RR_RRRRRRRRR_RRRRR_RRRRR_RRRRRRRRWR_RRRRRR_R_RR_RR__R__R_RR_R_R_R_R__R_W__RRR_WRR_R_R__WRRR_RRRRR__RRRRRW_WRRR___RRRRRR_R_RRRRRRR_RR__R_RR_RR_RRRRR__W_R_WR__RRRRRR__RR_RR_R_R_RWRWRRRR_____WRRRRRRRRR_R_RRR_RRRRRRRRR_RR____RRRRWRRR_RR_R_WRRR___R_RRRR_RRRR_RR_W___RW_RRR_RRWRRRWWRRR_RRRRR_R_R_RR_RWR_R_R__RR_RR_RRWR_R_WRRWRWRRRR_RRWRR_RWRR_R_RRRR____R___RRR_W_RR__RR_R_RRRRRRRRRWRRRRRWR_RW__R__WR_RRW_R_RWRRWRR_R_R_RR__W_R___RR_RRR_RRWRRRRRR__RRRRRW_RRRR_WRRRRRRR__RRRWRWR_RR_RR___RRRW___RW__R__RRRRRRRWR_RR_RRRRR_RRRRRWRRRRR_RRR_RWR_RRRR__RR_RRR__RRRWWR_RRWRR__R_R_RWRRR___RRRRRRRR_RR_RR_WRW_RR_RRRRR_RRRRR_R_RRRWWRRRR_WRRRRRRRRRRR_RR_WRRRRW_RR_WRWRWRRRRW_WR__W_RR__RRR_W_RR_R___R_R____RR_WRRR_R_R_R_WRWRRRR_WR_RR_RRWR__RRRRRRRR_RRR__R_R___RRRRRRR__R_R_R__W__RRRRR_RRRRRRR_R_R_RW_W____RR__R_RRW_R_RRRRRWRRR_R__R_RRRRRRR__RRRRWRRRRR_RRRWRRRWW_RW_R_R_RR_____RRRW_RRRRRWRWRR_R__R_R_R__RRRRRRRRR_RRRRRRRRRR_RRRRRR_R_R__WRR_RR_RR__WRRR_RR_RRRR_R_RWWR___R_R_RRRR_R_RRRRWRRR_RRRRR_RRRRRRR_R___RR_____RRR__RRW__R__RR_____R_RR_____RR_W___R_R_____RRRR__R_RR_R__RR______R_RRR_RRRRWRW_R__RR__W__RRRR_WR_R__RRR__RR_RRRRR__R_W_________RR__WR_____R__W_RR_RR_RW_RRRRWRRRR___RWRR_WRRRRRR__R_RRR__R__RR__WWRR__RRRRR__RRR_RRR_RR_RR_RRRRR_RRRRRR_RRR_WRW_RRRW_RRRR__R_RWRRRRR_RR__RRRWRRR____RRWRR____RRRRRRW_......................................................................................................................................................R____R_RRRR_R__RR__R_RWR_RR__R__R_RRRR___RRRR_R_RRW__R__RRR_RRRR_R____RR_RRRRR_RR_RR___RRRR_RR_RR__R_RR___R_R_R_RR___RR_R__R____RR_RR__R__R__RRR__WRRW..................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................

	ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Total Accesses: 491803362")
		fmt.Fprintln(w, "Total kBytes: 21453176370")
		fmt.Fprintln(w, "CPULoad: .84122")
		fmt.Fprintln(w, "Uptime: 3761632")
		fmt.Fprintln(w, "ReqPerSec: 130.742")
		fmt.Fprintln(w, "BytesPerSec: 5840030")
		fmt.Fprintln(w, "BytesPerReq: 44668.4")
		fmt.Fprintln(w, "BusyWorkers: 1298")
		fmt.Fprintln(w, "IdleWorkers: 652")
		fmt.Fprintln(w, "ConnsTotal: 3264")
		fmt.Fprintln(w, "ConnsAsyncWriting: 82")
		fmt.Fprintln(w, "ConnsAsyncKeepAlive: 1574")
		fmt.Fprintln(w, "ConnsAsyncClosing: 309")
		fmt.Fprintln(w, "ServerUptimeSeconds: 5632")
		fmt.Fprintln(w, "Load1: 1.01")
		fmt.Fprintln(w, "Load5: 1.50")
		fmt.Fprintln(w, "Load15: 2.06")
		fmt.Fprintln(w, "CPUUser: 0.5")
		fmt.Fprintln(w, "CPUSystem: .01")
		fmt.Fprintln(w, "CPUChildrenUser: .1")
		fmt.Fprintln(w, "CPUChildrenSystem: 1.99")
		fmt.Fprintln(w, "Scoreboard: RRR_R__RR_RLRR_RRRRRRRRRR_RRR_RRRRR__R___RRR_R_RRWRR_R______RWRR___RR_RR_RRRR_WR_R__W_RR__RRRR_____RR___RRRWRR__RRR_R_RRR_RW_RR_RRR_WRRRR__R__RRRR_RRRW_RRRRR_RRRRRR_RRRR_RRRWRWRRRRRR_R_R____R___RRR_W__RRRRRR_RR_W__R_RRRRRRR_RR_RRR____W_RRRR_RRRW_RRRRRR_RR_RR___R_RRRRRRRR_RR__RRRR__RRRRR_RWR__RRWRRRR__RR_RRWWR_R_RWRR__R_RR_RRRRRWR___RW_RR__RR_R_RWRRR_RRRRRRW___R_RRRR_RRR__RRRRWR___R_R__RRR__RRRW_RRRR__RRWR_RRR__L_R_R__R_R_R_R__WRRRRRRRR__RWRRRW_RRWRRRWWL__RR_R__R_RRR____RRRRRR_RWRRRR_R_RRR_R__R_RR_R_RRRRR_RRRW_RR_RRRRRRRRR_RRRRR_RRRRR_RRRRRRRRWR_RRRRRR_R_RR_RR__R__R_RR_R_R_R_R__R_W__RRR_WRR_R_R__WRRR_RRRRR__RRRRRW_WRRR___RRRRRR_R_RRRRRRR_RR__R_RR_RR_RRRRR__W_R_WR__RRRRRR__RR_RR_R_R_RWRWRRRR_____WRRRRRRRRR_R_RRR_RRRRRRRRR_RR____RRRRWRRR_RR_R_WRRR___R_RRRR_RRRR_RR_W___RW_RRR_RRWRRRWWRRR_RRRRR_R_R_RR_RWR_R_R__RR_RR_RRWR_R_WRRWRWRRRR_RRWRR_RWRR_R_RRRR____R___RRR_W_RR__RR_R_RRRRRRRRRWRRRRRWR_RW__R__WR_RRW_R_RWRRWRR_R_R_RR__W_R___RR_RRR_RRWRRRRRR__RRRRRW_RRRR_WRRRRRRR__RRRWRWR_RR_RR___RRRW___RW__R__RRRRRRRWR_RR_RRRRR_RRRRRWRRRRR_RRR_RWR_RRRR__RR_RRR__RRRWWR_RRWRR__R_R_RWRRR___RRRRRRRR_RR_RR_WRW_RR_RRRRR_RRRRR_R_RRRWWRRRR_WRRRRRRRRRRR_RR_WRRRRW_RR_WRWRWRRRRW_WR__W_RR__RRR_W_RR_R___R_R____RR_WRRR_R_R_R_WRWRRRR_WR_RR_RRWR__RRRRRRRR_RRR__R_R___RRRRRRR__R_R_R__W__RRRRR_RRRRRRR_R_R_RW_W____RR__R_RRW_R_RRRRRWRRR_R__R_RRRRRRR__RRRRWRRRRR_RRRWRRRWW_RW_R_R_RR_____RRRW_RRRRRWRWRR_R__R_R_R__RRRRRRRRR_RRRRRRRRRR_RRRRRR_R_R__WRR_RR_RR__WRRR_RR_RRRR_R_RWWR___R_R_RRRR_R_RRRRWRRR_RRRRR_RRRRRRR_R___RR_____RRR__RRW__R__RR_____R_RR_____RR_W___R_R_____RRRR__R_RR_R__RR______R_RRR_RRRRWRW_R__RR__W__RRRR_WR_R__RRR__RR_RRRRR__R_W_________RR__WR_____R__W_RR_RR_RW_RRRRWRRRR___RWRR_WRRRRRR__R_RRR__R__RR__WWRR__RRRRR__RRR_RRR_RR_RR_RRRRR_RRRRRR_RRR_WRW_RRRW_RRRR__R_RWRRRRR_RR__RRRWRRR____RRWRR____RRRRRRW_......................................................................................................................................................R____R_RRRR_R__RR__R_RWR_RR__R__R_RRRR___RRRR_R_RRW__R__RRR_RRRR_R____RR_RRRRR_RR_RR___RRRR_RR_RR__R_RR___R_R_R_RR___RR_R__R____RR_RR__R__R__RRR__WRRW..................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................")
	}))
	defer ts3.Close()

	ab3 := ApacheBeat{}
	u3, _ := url.Parse(ts3.URL)
	s3, _ := ab3.GetServerStatus(*u3)

	assert.Equal(t, s3["totalAccesses"].(int), 491803362)
	assert.Equal(t, s3["totalKBytes"].(int), 21453176370)
	assert.Equal(t, s3["reqPerSec"].(float64), 130.742)
	assert.Equal(t, s3["bytesPerSec"].(float64), float64(5840030))
	assert.Equal(t, s3["bytesPerReq"].(float64), 44668.4)
	assert.Equal(t, s3["busyWorkers"].(int), 1298)
	assert.Equal(t, s3["idleWorkers"].(int), 652)

	up3 := s3["uptime"].(common.MapStr)
	assert.Equal(t, up3["uptime"].(int), 3761632)
	assert.Equal(t, up3["serverUptimeSeconds"].(int), 5632)

	c3 := s3["connections"].(common.MapStr)
	assert.Equal(t, c3["connsTotal"].(int), 3264)
	assert.Equal(t, c3["connsAsyncWriting"].(int), 82)
	assert.Equal(t, c3["connsAsyncKeepAlive"].(int), 1574)
	assert.Equal(t, c3["connsAsyncClosing"].(int), 309)
	//assert.Equal(t, s3["host_url"], "localhost")

	l3 := s3["load"].(common.MapStr)
	assert.Equal(t, l3["load1"].(float64), 1.01)
	assert.Equal(t, l3["load5"].(float64), 1.50)
	assert.Equal(t, l3["load15"].(float64), 2.06)

	cpu3 := s3["cpu"].(common.MapStr)
	assert.Equal(t, cpu3["cpuLoad"].(float64), 0.84122)
	assert.Equal(t, cpu3["cpuUser"].(float64), 0.5)
	assert.Equal(t, cpu3["cpuSystem"].(float64), 0.01)
	assert.Equal(t, cpu3["cpuChildrenUser"].(float64), 0.1)
	assert.Equal(t, cpu3["cpuChildrenSystem"].(float64), 1.99)

	scb3 := s3["scoreboard"].(common.MapStr)
	assert.Equal(t, scb3["startingUp"].(int), 0)
	assert.Equal(t, scb3["readingRequest"].(int), 1150)
	assert.Equal(t, scb3["sendingReply"].(int), 145)
	assert.Equal(t, scb3["keepalive"].(int), 0)
	assert.Equal(t, scb3["dnsLookup"].(int), 0)
	assert.Equal(t, scb3["closingConnection"].(int), 0)
	assert.Equal(t, scb3["logging"].(int), 3)
	assert.Equal(t, scb3["gracefullyFinishing"].(int), 0)
	assert.Equal(t, scb3["idleCleanup"].(int), 0)
	assert.Equal(t, scb3["openSlot"].(int), 1800)
	assert.Equal(t, scb3["waitingForConnection"].(int), 652)
	assert.Equal(t, scb3["total"].(int), 3750)
}
