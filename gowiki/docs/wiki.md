# Example web application

Our example Go application is the wiki application thoroughly described in [Writing Web Applications](https://go.dev/doc/articles/wiki).

You can browse the application [source code](../../gowiki/wiki.go)

## Run wiki in a docker [container](https://www.docker.com/resources/what-container/)

1. Build docker image

The following command builds the `fosscomm/wiki` docker images using the `Dockerfile` located [here](./Dockerfile). To do so, switch to the [gowiki](../../gowiki) directory and run:

```bash
docker build -t fosscomm/wiki .
```

2. Inspect created docker image

After successful execution, one can inspect the details of the generated image using:

```bash
docker images fosscomm/wiki 
```

3. Start the container using the above image

```bash
docker run -dp 127.0.0.1:8080:8080 --name wiki fosscomm/wiki
```

4. Try it out!

Navigate to [http://localhost:8080/view/ANewPage](http://localhost:8080/view/ANewPage) to access the wiki.

### Cleanup

#### Stop the container
```bash
docker stop wiki
```

#### Delete the container

```bash
docker rm wiki
```

#### Delete the image

```bash
docker rmi fosscomm/wiki
```

## Run wki together with the other demo services

Alternatively, one can start the wiki application together with the other services we need for this demo listed in [docker-compose.yml](../docker-compose.yml). To do so, switch to the [root](../) directory and run:

### Start services

```bash
docker compose up -d
```

### Try it out!

Navigate to [http://localhost:8080/view/ANewPage](http://localhost:8080/view/ANewPage) to access the wiki.

Naviagate to [Grafana home](http://localhost:3000)

[Browse](http://localhost:3000/connections/datasources) provisioned datasources.

## Cleanup

```bash
docker compose down
```

# Resources

- [Writing Web Applications](https://go.dev/doc/articles/wiki/)