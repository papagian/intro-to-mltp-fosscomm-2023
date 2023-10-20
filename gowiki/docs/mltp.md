# Integrate wiki application intro Grafana Metrics, Logs, Traces and Profiles demo

One can also start the wiki application together with the other services we need for this demo listed in [docker-compose.yml](../../docker-compose-otel.yml).

1. Clone [https://github.com/grafana/intro-to-mltp](https://github.com/grafana/intro-to-mltp)

```bash
git clone https://github.com/papagian/intro-to-mltp.git
```

2. Add the following lines at the end of the `docker-compose-otel.yml` file:

```
  gowiki:
    build:
      context: ./gowiki
      dockerfile: Dockerfile
    ports:
      - 8080:8080
```

3. Start services using docker compose:
```bash
docker compose -f docker-compose-otel.yml up
```

4. Access wiki application at [http://localhost:8080/view/ANewPage](http://localhost:8080/view/ANewPage).

