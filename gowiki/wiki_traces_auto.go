// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/grafana/pyroscope-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	tracer = otel.Tracer("gowiki")
)

var (
	totalPages = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "wiki_pages_total",
		Help: "The total number of processed events",
	}, []string{"handler"})
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(ctx, w, "view", p)
}

func editHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(ctx, w, "edit", p)
}

func saveHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(ctx context.Context, w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(ctx context.Context, tracer trace.Tracer, fn func(context.Context, http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	ctx, span := tracer.Start(
		ctx,
		"makeHandler",
		trace.WithAttributes(attribute.String("parentAttributeKey1", "parentAttributeValue1")))

	//span.AddEvent("ParentSpan-Event")
	logger.Debug("In parent span")
	defer span.End()

	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		start := time.Now()
		fn(ctx, w, r, m[2])
		logger.Info("request completed", "handler", m[1], "title", m[2], "duration", time.Now().Sub(start), "traceID", span.SpanContext().TraceID().String())

		totalPages.WithLabelValues(m[1]).Inc()
	}
}

func main() {
	cleanup := initTracerAuto()
	ctx := context.Background()
	defer func() {
		err := cleanup(ctx)
		logger.Error("shutdown failure", "err", err)
	}()

	tracer := otel.Tracer("gowiki")

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

	http.HandleFunc("/view/", makeHandler(ctx, tracer, viewHandler))
	http.HandleFunc("/edit/", makeHandler(ctx, tracer, editHandler))
	http.HandleFunc("/save/", makeHandler(ctx, tracer, saveHandler))

	//send the default metrics for our Go server
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("failed to start server", "err", err)
	}
}

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
