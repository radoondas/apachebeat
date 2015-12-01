# ApacheBeat
Current status: **DEVELOPMENT release**.
Do not use in Production. Beat will be highly developed in the following days.

## Description
This is beat for Apache HTTPD server-status page. 
More about beats platform: https://www.elastic.co/products/beats

ApacheBeat polls Apache httpd server-status page every 'defined' period. You can poll more URL's at once and save your results to ElasticSearch. Each document in Elasticsearch consists of metrics gathered from server-status page and add url.host to see which URL answered. Each document in ElasticSearch is flat document with no nested objects.

## Thanks note
Beat is highly motivated by nginxbeat (https://github.com/mrkschan/nginxbeat). In fact nginxbeat served as a template. Thanks!!