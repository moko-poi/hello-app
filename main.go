package main

import (
	"io"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	helloPath      = "/hello"
	wrongPath      = "/wrong"
	monitoringPath = "/metrics"
)

var (
	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"handler", "code", "method"},
	)
)

func init() {
	prometheus.MustRegister(counter)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hello, world!")
}

func wrongHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, "Error!")
}

func instrumentHandlerMiddleware(handlerPath string, handler http.HandlerFunc) http.Handler {
	return promhttp.InstrumentHandlerCounter(
		counter.MustCurryWith(prometheus.Labels{"handler": handlerPath}),
		handler,
	)
}

func main() {
	http.Handle(helloPath, instrumentHandlerMiddleware(helloPath, helloHandler))
	http.Handle(wrongPath, instrumentHandlerMiddleware(wrongPath, wrongHandler))
	http.Handle(monitoringPath, promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
