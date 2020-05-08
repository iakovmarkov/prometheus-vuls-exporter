package metrics

import (
	"strings"

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

	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "reported_at",
		Help: "Timestamp of last report time, in ms since Unix",
	}, func() float64 {
		var timestamp = reportedAt.Unix()
		return float64(timestamp)
	})

	metrics = append(metrics, Metric{
		prom: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "vuln_severity",
			Help: "Vulnerability count by severity",
		}, []string{"severity"}),
		record: func(metric Metric) {
			metric.prom.(*prometheus.GaugeVec).Reset()
			for _, report := range reports {
				for _, cve := range report.cves {
					severity := strings.ToLower(cve.severity)
					metric.prom.(*prometheus.GaugeVec).WithLabelValues(severity).Inc()
				}
			}
		},
	})

	metrics = append(metrics, Metric{
		prom: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "vulns",
			Help: "Vulnerability information, value represents total amount of hits",
		}, []string{
			"cveID",
			"severity",
			"packageName",
			"notFixedYet",
			"fixState",
			"title",
			"summary",
			"published",
			"lastModified",
			"mitigation",
		}),
		record: func(metric Metric) {
			metric.prom.(*prometheus.GaugeVec).Reset()
			for _, report := range reports {
				for _, cve := range report.cves {
					isFixed := "false"
					if cve.notFixedYet {
						isFixed = "true"
					}
					metric.prom.(*prometheus.GaugeVec).WithLabelValues(
						cve.id,
						cve.severity,
						cve.packageName,
						isFixed,
						cve.fixState,
						cve.title,
						cve.summary,
						cve.published,
						cve.lastModified,
						cve.mitigation,
					).Inc()
				}
			}
		},
	})

	metrics = append(metrics, Metric{
		prom: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "servers",
			Help: "Server information, value represents amount of vulnerabilitites",
		}, []string{
			"serverName",
			"hostname",
			"cveID",
			"kernel_release",
			"kernel_rebootRequired",
		}),
		record: func(metric Metric) {
			metric.prom.(*prometheus.GaugeVec).Reset()
			for _, report := range reports {
				rebootRequired := "false"
				if report.kernel.rebootRequired {
					rebootRequired = "true"
				}
				for _, cve := range report.cves {
					metric.prom.(*prometheus.GaugeVec).WithLabelValues(
						report.serverName,
						report.hostname,
						cve.id,
						report.kernel.release,
						rebootRequired,
					).Inc()
				}
			}
		},
	})
}
