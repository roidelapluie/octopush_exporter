# Octopush metrics for Prometheus

This is an exporter for getting your [Octopush](https://www.octopush-dm.com/)
balance as prometheus metrics.

## Configuration

```
---
- login: abcd@sub-accounts.com
  key: 123123123sadfasfdfdfds
  labels:
    route: 1
    route_env: prod
    type: direct-32
```

The configuration is a list, you can supply multiple accounts.

## Building

```
$ dep ensure
$ go build
```


## How to run

```
./octopush_exporter -config conf.yml
```
