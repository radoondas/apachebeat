[![Travis Build Status](https://travis-ci.org/elastic/libbeat.svg?branch=master)](https://travis-ci.org/radoondas/apachebeat)

# ApacheBeat
Current status: **beta release**.

## Description
This is beat for Apache HTTPD [server-status](https://httpd.apache.org/docs/2.4/mod/mod_status.html) page. ApacheBeat polls Apache HTTPD server-status page every 'defined' period. You can poll more URL's at once and save your results to ElasticSearch. Each document in ElasticSearch consists of metrics gathered from server-status page and add url.host to see which URL answered. Each document in ElasticSearch is flat document with no nested objects.

Document example:
```json
{
  "_index": "apachebeat-2015.12.05",
  "_type": "apache_status",
  "_id": "AVFvpKJ21NqxaroAvAlC",
  "_score": null,
  "_source": {
    "@timestamp": "2015-12-05T00:57:18.887Z",
    "apache": {
      "busy_workers": 184,
      "bytes_per_req": "42878.3",
      "bytes_per_sec": "5678720",
      "conns_async_closing": 153,
      "conns_async_keep_alive": 486,
      "conns_async_writing": 18,
      "conns_total": 841,
      "cpu_load": "0.817271",
      "host_url": "www.apache.org",
      "idle_workers": 416,
      "req_per_sec": "132.438",
      "scb_closing_connection": 0,
      "scb_dns_lookup": 0,
      "scb_gracefully_finishing": 0,
      "scb_idle_cleanup": 0,
      "scb_keepalive": 0,
      "scb_logging": 1,
      "scb_open_slot": 3150,
      "scb_reading_request": 62,
      "scb_sending_reply": 121,
      "scb_starting_up": 0,
      "scb_waiting_for_connection": 416,
      "total_access": 368347986,
      "total_kbytes": 15423958662,
      "uptime": 2781282
    },
    "beat": {
      "hostname": "hostname",
      "name": "hostname"
    },
    "count": 1,
    "source": "http://www.apache.org/server-status?auto",
    "type": "apache_status"
  }
```

More about beats platform: https://www.elastic.co/products/beats

## To apply ApacheBeat template:

```bash
curl -XPUT 'http://localhost:9200/_template/apachebeat' -d@apachebeat.template.json
```

## Thanks note
Beat is highly motivated by [nginxbeat](https://github.com/mrkschan/nginxbeat). In fact nginxbeat served as a template. Thanks!!
