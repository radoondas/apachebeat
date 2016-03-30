[![Travis Build Status](https://travis-ci.org/radoondas/apachebeat.svg?branch=master)](https://travis-ci.org/radoondas/apachebeat)

# ApacheBeat
Current status: **beta release**.

## Description
This is beat for Apache HTTPD [server-status](https://httpd.apache.org/docs/2.4/mod/mod_status.html) page. ApacheBeat polls Apache HTTPD server-status page every 'defined' period. You can poll more URL's at once and save your results to ElasticSearch. Each document in ElasticSearch consists of metrics gathered from server-status page and add url.host to see which URL answered. Each document in ElasticSearch is flat document with no nested objects.

Document example:
```json
{
    "_index": "apachebeat-2016.03.29",
    "_type": "apache_status",
    "_id": "AVPBdyRHD3Lkxnx3Btq7",
    "_version": 1,
    "_score": 1,
    "_source": {
        "@timestamp": "2016-03-29T08:22:04.102Z",
        "apache": {
            "busyWorkers": 263,
            "bytesPerReq": 44679.9,
            "bytesPerSec": 4895320,
            "connections": {
                "connsAsyncClosing": 612,
                "connsAsyncKeepAlive": 1483,
                "connsAsyncWriting": 475,
                "connsTotal": 2839
            },
            "cpu": {
                "cpuChildrenSystem": 0,
                "cpuChildrenUser": 0,
                "cpuLoad": 2.70362,
                "cpuSystem": 0,
                "cpuUser": 0
            },
            "hostname": "www.apache.org",
            "idleWorkers": 637,
            "load": {
                "load1": 0,
                "load15": 0,
                "load5": 0
            },
            "reqPerSec": 109.564,
            "scoreboard": {
                "closingConnection": 0,
                "dnsLookup": 0,
                "gracefullyFinishing": 113,
                "idleCleanup": 0,
                "keepalive": 0,
                "logging": 337,
                "openSlot": 2400,
                "readingRequest": 232,
                "sendingReply": 31,
                "startingUp": 0,
                "total": 3750,
                "waitingForConnection": 637
            },
            "totalAccesses": 141301656,
            "totalKBytes": 6165381544,
            "uptime": {
                "serverUptimeSeconds": 0,
                "uptime": 1289672
            }
        },
        "beat": {
            "hostname": "hostname",
            "name": "name"
        },
        "type": "apache_status",
        "url": "http://www.apache.org/server-status?auto"
    }
}
```

More about beats platform: https://www.elastic.co/products/beats

## To apply ApacheBeat template:

```bash
curl -XPUT 'http://localhost:9200/_template/apachebeat' -d@apachebeat.template.json
```

## Example Kibana dashboard
![Apache HTTPD server-status](/docs/images/apache-server-status.png)

## Links
* [Simple Run guide](/RUN.md)
* [Kibana guide](/KIBANA.md)
* [Kibana samples](/kibana/dashboards)

## Thanks note
Beat is highly motivated by [nginxbeat](https://github.com/mrkschan/nginxbeat). In fact nginxbeat served as a template. Thanks!!
