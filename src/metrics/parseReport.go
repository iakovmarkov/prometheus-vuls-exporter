package metrics

import (
	"fmt"
	"os"
	"strings"

	"../utils"
	"github.com/tidwall/gjson"
)

func getterFactory(jsonString string) func(path string, a ...interface{}) gjson.Result {
	return func(path string, a ...interface{}) gjson.Result {
		finalPath := fmt.Sprintf(path, a...)
		// log.Printf("Trying to get data at path: %s", finalPath)
		return gjson.Get(jsonString, finalPath)
	}
}

func getServerName(file os.FileInfo) string {
	filename := file.Name()
	lastDot := strings.LastIndex(filename, ".")
  	serverName := filename[0:lastDot]	
	return serverName
}

func parseReport(file os.FileInfo) Report {
	var r Report

	// Get basic file info
	filePath := fmt.Sprintf("%s/%s", latestPath, file.Name())
	r.filename = file.Name()
	r.serverName = getServerName(file)
	r.path = filePath

	// log.Printf("Parsing report: %s", file.Name())

	// Get JSON contents
	jsonString := string(utils.ReadFile(filePath))
	getData := getterFactory(jsonString)

	// Basic host information
	r.hostname = getData("config.report.servers.%s.host", r.serverName).String()

	// Kernel information
	r.kernel = KernelInfo{
		rebootRequired: getData("runningKernel.rebootRequired").Bool(),
		release:        getData("runningKernel.release").String(),
	}

	// Vulnerability information
	var cves []CVEInfo
	for _, c := range getData("scannedCves").Map() {
		// TODO: This should scan through other properties than "nvd"!

		severity := c.Get("cveContents.nvd.cvss2Severity").String()
		if severity == "" {
			severity = "UNKNOWN"
		}

		cve := CVEInfo{
			id:           c.Get("cveID").String(),
			packageName:  c.Get("affectedPackages.0.name").String(),
			severity:     severity,
			fixState:     c.Get("affectedPackages.0.fixState").String(),
			notFixedYet:  c.Get("affectedPackages.0.notFixedYet").Bool(),
			title:        c.Get("cveContents.nvd.title").String(),
			summary:      c.Get("cveContents.nvd.summary").String(),
			published:    c.Get("cveContents.nvd.published").String(),
			lastModified: c.Get("cveContents.nvd.lastModified").String(),
			mitigation:   c.Get("cveContents.nvd.mitigation").String(),
		}
		cves = append(cves, cve)
	}

	r.cves = cves

	// Debug
	// log.Printf("Report:\n")
	// log.Printf("%+v\n\n", r.cves)

	return r
}
