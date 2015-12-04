# Short RUN guide

First, setup Go lang environment (https://golang.org/doc/install)
Add to your .bashrc important variables

```bash
export GOROOT="$HOME/opt/go"
export GOPATH="$HOME/workspace/go"
export PATH="$HOME/opt/go/bin:"$PATH
```

## Install apache beat

```bash
go get github.com/radoondas/apachebeat
```

## Elastic and Kibana
Meanwhile setup your ElasticSearch and Kibana (example dashbords coming soon)

## Build ApacheBeat 

```bash
cd ~/workspace/go/src/github.com/radoondas/apachebeat
go install
```

## Import template
```bash
cd ~/workspace/go/src/github.com/radoondas/apachebeat/etc
curl -XPUT 'http://localhost:9200/_template/apachebeat' -d@apachebeat.template.json
```

## Run ApacheBeat

Following command will execute ApacheBeat with debug option and will not index results in to ES. Instead, you will see output on the screen.
```bash
cd ~/workspace/go/bin
./apachebeat  -e -v -d apachebeat -c ~/workspace/go/src/github.com/radoondas/apachebeat/etc/apachebeat.yml
```

With no debug options - just do straight indexing to your ES installation

```bash
./apachebeat  -e -c ~/workspace/go/src/github.com/radoondas/apachebeat/etc/apachebeat.yml
```
