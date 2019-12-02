package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	exporter := NewExporter("myapp")
	prometheus.MustRegister(exporter)

	mux := http.NewServeMux()
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		exporter.IncreaseCounter()
		exporter.IncreaseCounterByEndpointAndMethod(r.URL.Path, r.Method)

		w.Write([]byte("hello world!"))
	})

	mux.Handle("/metrics", promhttp.Handler())

	s := http.Server{Addr:":8089", Handler: mux}

	log.Println("starting server", s.Addr)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal("shutting down", err)
	}
}

type Exporter struct {
	counter    prometheus.Counter
	counterVec prometheus.CounterVec
}

func NewExporter(metricsPrefix string) *Exporter {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: metricsPrefix,
		Name:      "counter_metric",
		Help:      "This is a counter for number of total api calls"})

	counterVec := *prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricsPrefix,
		Name:      "counter_vec_metric",
		Help:      "This is a counter vec for number of all api calls"},
		[]string{"endpoint", "method"})

	return &Exporter{
		counter:    counter,
		counterVec: counterVec,
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.counter.Collect(ch)
	e.counterVec.Collect(ch)
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.counter.Describe(ch)
	e.counterVec.Describe(ch)
}

func (e *Exporter) IncreaseCounter() {
	e.counter.Inc()
}

func (e *Exporter) IncreaseCounterByEndpointAndMethod(endpoint string, method string) {
	e.counterVec.With(prometheus.Labels{"endpoint": endpoint, "method": method}).Inc()
}
