package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	dummyMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "Dummy",
		Help: "Dummy Metric",
	})
)

func init() {
	flag.String("address", ":8080", "The address to listen on for HTTP requests.")
	flag.String("logFormat", "LONG", "Log format - LONG or SHORT.")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()

	log.SetPrefix("prometheus-vuls-exporter ")
	if viper.GetString("logFormat") == "SHORT" {
		log.SetFlags(log.Lmsgprefix)
	} else {
		log.SetFlags(log.Ldate + log.Ltime + log.Lmsgprefix)
	}
}

func main() {
	dummyMetric.Set(1)

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("listening on %s", viper.GetString("address"))
	log.Fatal(http.ListenAndServe(viper.GetString("address"), nil))
}
