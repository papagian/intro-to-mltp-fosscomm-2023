# Workshop: Monitor and debug a Go application with Prometheus and Grafana open source tooling

In this workshop we will extract metrics from our example Go application and use them for creating some graphs. We will query the logs and create some metrics and alerts from them. We will also collect some traces and profiles in order to detect the sources of failures and analyze its performance.

# Agenda

1. Some [Observability basics](docs/o11y.md)
2. Meet our [example web application](docs/wiki.md)
3. Integrate the wiki application into [Grafana MLTP demo](docs/mltp.md)
4. [Add metrics](./docs/add_metrics.md) into wiki application and create an alert
5. [Add logs](./docs/add_logs.md) into wiki application
6. Add [continuous profiling](./docs/add_profiles.md)
7. [Add tracing](./docs/add_traces.md)

# Prerequisites
- [docker](https://docs.docker.com/engine/install/)
- [docker compose](https://docs.docker.com/compose/install/)

# Resources
- [Grafana Loki](https://grafana.com/oss/loki/)
- [Send log data to Loki](https://grafana.com/docs/loki/latest/send-data/)
- [Docker driver client](https://grafana.com/docs/loki/latest/send-data/docker-driver/)