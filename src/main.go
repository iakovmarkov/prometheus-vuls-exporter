package main

import (
	"flag"
	"log"
	"net/http"

	Util "./util"

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
	flag.String("log_format", "LONG", "Log format - LONG or SHORT.")
	flag.String("basic_username", "", "Log format - LONG or SHORT.")
	flag.String("basic_password", "", "Log format - LONG or SHORT.")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()

	log.SetPrefix("prometheus-vuls-exporter ")
	if viper.GetString("log_format") == "SHORT" {
		log.SetFlags(log.Lmsgprefix)
	} else {
		log.SetFlags(log.Ldate + log.Ltime + log.Lmsgprefix)
	}
}

func main() {
	dummyMetric.Set(1)

	var authHandler = Util.HTTPBasicAuthHandler(viper.GetString("basic_username"), viper.GetString("basic_password"))
	var promHandler = promhttp.Handler().(http.HandlerFunc)
	var handler = Util.Use(
		promHandler,
		authHandler,
	)

	http.Handle("/metrics", handler)

	log.Printf("listening on %s\n", viper.GetString("address"))
	log.Fatal(http.ListenAndServe(viper.GetString("address"), nil))
}
