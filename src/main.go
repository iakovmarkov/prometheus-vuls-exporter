package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"./metrics"
	"./utils"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	flag.String("reports_dir", "", "The folder where Vulns stores JSON reports.")
	flag.String("address", ":8080", "The address to listen on for HTTP requests.")
	flag.String("log_format", "LONG", "Log format - LONG or SHORT.")
	flag.String("basic_username", "", "Log format - LONG or SHORT.")
	flag.String("basic_password", "", "Log format - LONG or SHORT.")
	flag.Bool("version", false, "Print version and exit")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()

	if viper.GetBool("version") {
		version := utils.ReadFile("./VERSION")
		log.Printf("prometheus-vuls-exporter v%s", version)
		os.Exit(0)
	}

	log.SetPrefix("prometheus-vuls-exporter ")
	if viper.GetString("log_format") == "SHORT" {
		log.SetFlags(log.Lmsgprefix)
	} else {
		log.SetFlags(log.Ldate + log.Ltime + log.Lmsgprefix)
	}
}

func main() {
	if viper.GetString("reports_dir") == "" {
		log.Fatalln("reports_dir is not configured, exiting...")
	}

	metrics.CreateMetrics(viper.GetString("reports_dir"))

	authHandler := utils.HTTPBasicAuthHandler(viper.GetString("basic_username"), viper.GetString("basic_password"))
	metricCollectionHandler := metrics.MetricCollectionHandler(viper.GetString("reports_dir"))
	promHandler := promhttp.Handler().(http.HandlerFunc)
	handler := utils.Use(
		promHandler,
		metricCollectionHandler,
		authHandler,
	)

	http.Handle("/metrics", handler)

	log.Printf("listening on %s\n", viper.GetString("address"))
	log.Fatal(http.ListenAndServe(viper.GetString("address"), nil))
}
