package metrics

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type JSONData = map[string]interface{}

type KernelInfo struct {
	rebootRequired bool
	release        string
}

type CVEInfo struct {
	id             string
	packageName    string
	severity       string
	notFixedYet    bool
	fixState       string
	title          string
	summary        string
	referenceLinks []string
	published      string
	lastModified   string
	mitigation     string
}

type Report struct {
	filename   string
	path       string
	serverName string
	hostname   string
	kernel     KernelInfo
	cves       []CVEInfo
}

var (
	reportsPath string
	reportsDir  []os.FileInfo
	latestPath  string
	latestDir   []os.FileInfo
)

func CreateMetrics(reportsDir string) {
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "vulnerability_total",
		Help: "Total count of vulnerabilities, across all hosts",
	}, createMetric(reportsDir))
}
