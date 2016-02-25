# Short RUN guide

First, setup [Go lang environment](https://golang.org/doc/install) (if you don't have it already)
Add to your .bashrc important variables

```bash
export GOROOT="$HOME/opt/go"
export GOPATH="$HOME/workspace/go"
export PATH="$GOROOT/bin:$PATH"
# to enable vendor experiment in GO 1.5
export GO15VENDOREXPERIMENT=1
```

## Build ApacheBeat

```bash
cd $GOPATH
mkdir -p src/github.com/radoondas
cd src/github.com/radoondas
git clone https://github.com/radoondas/apachebeat.git
```

## Elastic and Kibana
Meanwhile setup your ElasticSearch and Kibana (example [dashbords](https://github.com/radoondas/apachebeat/tree/master/kibana))

## Build ApacheBeat

```bash
cd $GOPATH/src/github.com/radoondas/apachebeat
make
```

## Delete template (Optional)
If you need for any reason to delete old template, use following method.

```bash
curl -XDELETE 'http://localhost:9200/_template/apachebeat'
```

## Import template
```bash
cd ~/workspace/go/src/github.com/radoondas/apachebeat/etc
curl -XPUT 'http://localhost:9200/_template/apachebeat' -d@apachebeat.template.json
```

## Run ApacheBeat

Following command will execute ApacheBeat with debug option and will not index results in to ES. Instead, you will see output on the screen.
```bash
cd $GOPATH/src/github.com/radoondas/apachebeat
./apachebeat  -e -v -d apachebeat -c apachebeat.yml
```

With no debug options - just do straight indexing to your ES installation

```bash
./apachebeat -c apachebeat.yml
```
