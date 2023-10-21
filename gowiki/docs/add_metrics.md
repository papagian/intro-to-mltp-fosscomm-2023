# Add metrics to the wiki application

1. Modify [wiki.go](../source/wiki.go) and import prometheus Go libraries:

```
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
```

2. Register our metric for counting the number of pages accessed. Our new metric is called: `wiki_pages_total` and we define also one label: `handler`.

```
var (
	totalPages = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "wiki_pages_total",
		Help: "The total number of processed events",
	}, []string{"handler"})
)
```

3. Modify [makeHandler()](../source/wiki.go#L75) and add the following line for incresing the above metric when accessing a page:

```
		totalPages.WithLabelValues(m[1]).Inc()
```

4. Modify [main()](../source/wiki.go#L86) to send the default metrics for our Go server. To do so: add the following lines before 

```
	//send the default metrics for our Go server
	http.Handle("/metrics", promhttp.Handler())
```

At the end your wiki.go should look like [this](../../gowiki/wiki_metrics.go)

5. Modify [otel configuration](../../otel/otel.yml) so that it will scrape metrics also from our wiki application. To do so, add a new block under `config/scrape_configs` like:

```
        # Scrape from the gowiki service.
        - job_name: 'gowiki'
          scrape_interval: 2s
          static_configs:
            - targets: ['gowiki:8080']
              labels:
                service: 'gowiki'
                group: 'fosscomm'
```

6. Rebuild and restart the services. Switch to the [root](../..) and run:

```bash
docker-compose -f docker-compose-otel.yml up --force-recreate --build -d
```

7. Access [the web application](http://localhost:8080/view/fosscomm2023)

8. Access the [metrics](http://localhost:8080/metrics) that the wiki application exposes. Make sure that `wiki_pages_total` is there. Note the different labels.

6. Open Grafana [Explore](http://localhost:3000/explore).

7. Change the time range to `Last 5 minutes`.

8. Select `Mimir` data source.

9. Search for the `wiki_pages_total` metric. Narrow down results by selecting label `service` to equal `gowiki` and summarize by `handler`. 
Also click on operations and select `Range functions` -> `Rage`.
If you switch to the Code view the [PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/) query should look like:

```
sum by(handler) (rate(wiki_pages_total{service="gowiki"}[$__rate_interval]))
```

10. Toggle `Explain` and try to understand what this query measures.

11. Select `Add` -> `Add to dashboard` to store it into a new Dashboard. Save it.

12. Switch to `Alert` tab and click `Create alert rule from this panel`. You will be redirected to the page for creating an alert rule. Follow instructions for [creating an Grafana-managed alert rule](https://grafana.com/docs/grafana/latest/alerting/alerting-rules/create-grafana-managed-rule/?plcmt=footer).

13. Navifate to [Alert rules apge](http://localhost:3000/alerting/list) to observe the alert rule state.

Resources

- [Instrumenting a Go application for Prometheus](https://prometheus.io/docs/guides/go-application/)