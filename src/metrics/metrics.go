package metrics

import (
	"os"
	"time"
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

type Metric struct {
	prom   interface{}
	record func(metric Metric)
}

var (
	reportsPath string
	reportsDir  []os.FileInfo
	latestPath  string
	latestDir   []os.FileInfo
	reports     []Report
	reportedAt  time.Time
	metrics     []Metric
)
