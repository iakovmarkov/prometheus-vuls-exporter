package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func CreateMetrics(reportsDir string) {
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "vuln_count",
		Help: "Total count of vulnerabilities, across all servers",
	}, func() float64 {
		count := 0
		for _, report := range reports {
			count = count + len(report.cves)
		}
		return float64(count)
	})

	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "server_count",
		Help: "Total count of servers reported",
	}, func() float64 {
		count := len(reports)
		return float64(count)
	})
}
