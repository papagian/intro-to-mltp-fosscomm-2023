# Add logs into wiki web application

1. Import `log/slog`

2. Modify [makeHandler()](../../gowiki/wiki.go#L86) by wrapping `fn()` execution like:

```
start := time.Now()
fn(w, r, m[2])
logger.Info("request completed", "handler", m[1], "title", m[2], "duration", time.Now().Sub(start))
```

3. Modify [main()](../../gowiki/wiki.go#L105) like:

```
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("failed to start server", "err", err)
	}
```

At the end your wiki.go should look like [this](../../gowiki/wiki_logs.go)

2. [Install the Docker driver client](https://grafana.com/docs/loki/latest/send-data/docker-driver/#install-the-docker-driver-client):

```
docker plugin install grafana/loki-docker-driver:2.9.1 --alias loki --grant-all-permissions
```

3. Modify gowiki service in [](../../docker-compose-otel.yml) to use the loki logging driver:

```
  gowiki:
    build:
      context: ./gowiki
      dockerfile: Dockerfile
    logging:
      driver: loki
      options:
        loki-url: http://loki:3100/loki/api/v1/push
        loki-pipeline-stages: |
          - regex:
              expression: '(level|lvl|severity)=(?P<level>\w+)'
          - labels:
              level:
    ports:
      - 8080:8080
```

6. Rebuild and restart the services. Switch to the [root](../..) and run:

```bash
docker-compose -f docker-compose-otel.yml up --force-recreate --build -d
```

7. Access [the web application](http://localhost:8080/view/fosscomm2023)

8. Check gowiki logs:

```bash
docker-compose -f docker-compose-otel.yml logs gowiki
```

9. Open Grafana [Explore](http://localhost:3000/explore).

10. Change the time range to `Last 5 minutes`.

11. Select `Loki` data source.

12. Switch to `Code` view and paste the following query:

```
{compose_service="gowiki"} |= `` | json | duration > 100000
```

13. Toggle `Explain` and try to understand what the above query does.

14. Switch to `Builder` view and experiment with other filtering options.

# Resources
- [Grafana Loki](https://grafana.com/oss/loki/)
- [Send log data to Loki](https://grafana.com/docs/loki/latest/send-data/)
- [Docker driver client](https://grafana.com/docs/loki/latest/send-data/docker-driver/)
- [LogQL](https://grafana.com/docs/loki/latest/query/)