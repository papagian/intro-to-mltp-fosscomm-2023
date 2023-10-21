# Add tracing into the gowiki web application

1. Modify [wiki.go](..../gowiki/wiki.go#L89) to look like:

```
func initTracerAuto() func(context.Context) error {
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint("opentelemetry-collector:4317"),
		),
	)

	if err != nil {
		logger.Error("Could not set exporter", "err", err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "gowiki"),
			attribute.String("application", "fosscomm2023"),
		),
	)
	if err != nil {
		logger.Error("Could not set resources", "err", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
			sdktrace.WithSyncer(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}
```

2. Add the following lines at the beginning of `main()`:

```
	cleanup := initTracerAuto()
	ctx := context.Background()
	defer func() {
		err := cleanup(ctx)
		logger.Error("shutdown failure", "err", err)
	}()
```

3. Modify `makeHanlder()` and propagate it in all the handlers to expect an addition context argument and store the traceID into the logs.

```
func makeHandler(ctx context.Context, fn func(context.Context, http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	span := trace.SpanFromContext(ctx)
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		start := time.Now()
		fn(ctx, w, r, m[2])
		logger.Info("request completed", "handler", m[1], "title", m[2], "duration", time.Now().Sub(start), span.SpanContext().TraceID().String())

		totalPages.WithLabelValues(m[1]).Inc()
	}
}
```

4. Propagate context down to hanlders

At the end your wiki.go should look like [this](../../gowiki/wiki_traces_auto.go)

5. Update `gowiki` service at [/docker-compose-otel.ym](../../docker-compose-otel.yml) to look like:

```
  gowiki:
    build:
      context: ./gowiki
      dockerfile: Dockerfile
    logging:
      driver: loki
      options:
        loki-url: http://host.docker.internal:3100/loki/api/v1/push
        loki-pipeline-stages: |
          - regex:
              expression: '(level|lvl|severity)=(?P<level>\w+)'
          - labels:
              level:
    depends_on:
      - opentelemetry-collector
    ports:
      - 8080:8080
    environment:
      - OTEL_EXPORTER_OTLP_TRACES_INSECURE=true
      - OTEL_RESOURCE_ATTRIBUTES=ip=1.2.3.4
```

6. Rebuild and restart the services. Switch to the [root](../..) and run:

```bash
docker-compose -f docker-compose-otel.yml up --force-recreate --build -d
```

7. Access [the web application](http://localhost:8080/view/fosscomm2023)

8. Check gowiki logs:

```bash
docker-compose -f docker-compose-otel.yml logs -f gowiki
```

9. Copy `TraceID` value from one of the log entries.


# Resources

- [Instrument for distributed tracing](https://grafana.com/docs/tempo/latest/getting-started/instrumentation/?pg=oss-tempo&plcmt=resources)
- [Getting Started with OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/getting-started/)
- [Get started with Grafana Tempo](https://grafana.com/docs/tempo/latest/getting-started/?pg=oss-tempo&plcmt=resources)
- [Exporters](https://opentelemetry.io/docs/instrumentation/go/exporters/)