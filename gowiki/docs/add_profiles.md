# Add continuous profiling into wiki application

1. Modify [main()](../../gowiki/wiki.go) to look like this:

```
func main() {
	// These 2 lines are only required if you're using mutex or block profiling
	// Read the explanation below for how to set these rates:
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	pyroscope.Start(pyroscope.Config{
		ApplicationName: "gowiki",

		// replace this with the address of pyroscope server
		ServerAddress: "http://pyroscope:4040",

		// you can disable logging by setting this to nil
		Logger: pyroscope.StandardLogger,

		// you can provide static tags via a map:
		Tags: map[string]string{"hostname": os.Getenv("HOSTNAME")},

		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	//send the default metrics for our Go server
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

2. Rebuild and restart the services. Switch to the [root](../..) and run:

```bash
docker-compose -f docker-compose-otel.yml up --force-recreate --build -d
```

3. Access [the web application](http://localhost:8080/view/fosscomm2023)

4. Open Grafana [Explore](http://localhost:3000/explore).

5. Select `Pyroscope` as data source.

6. Select `process-cpu_cpu` as profile type.

7. Filter out by `{service_name="gowiki"}`