# Workshop: Monitor and debug a Go application with Prometheus and Grafana open source tooling

In this workshop we will extract metrics from our example Go application and use them for creating some graphs. We will query the logs and create some metrics and alerts from them. We will also collect some traces and profiles in order to detect the sources of failures and analyze its performance.

# Agenda

1. Some [Observability basics](docs/o11y.md)
2. Meet our [example web application](docs/wiki.md)
3. Integrate the wiki application into [Grafana MLTP demo](docs/mltp.md)
4. Add metrics into wiki application
5. Add logs
6. Add tracing
7. Continuous profiling
8. Create an alert

# Prerequisites
- [docker](https://docs.docker.com/engine/install/)
- [docker compose](https://docs.docker.com/compose/install/)

# Resources
- [Observability Whitepaper](https://github.com/cncf/tag-observability/blob/whitepaper-v1.0.0/whitepaper.md)
- [Introduction to Metrics, Logs, Traces and Profiles in Grafana](https://github.com/grafana/intro-to-mltp)
- [Writing Web Applications](https://go.dev/doc/articles/wiki/)
- [Instrumenting a Go application for Prometheus](https://prometheus.io/docs/guides/go-application/)