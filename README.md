# Prometheus Demo Exporter

This is a quick container that can be used for demoing or debugging a
prometheus config. It was thrown together while trying to debug some issues
with running prometheus in Kubernetes and we wanted a minimal application that
responds to /metrics.

```
curl localhost:1845/metrics
```

## Docker Image

```
docker pull gaffneyc/prom-demo-exporter
```
