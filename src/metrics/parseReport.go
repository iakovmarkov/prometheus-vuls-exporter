package metrics

import (
	"fmt"
	"log"
	"os"
	"strings"

	"../utils"
	"github.com/tidwall/gjson"
)

func getterFactory(jsonString string) func(path string, a ...interface{}) gjson.Result {
	return func(path string, a ...interface{}) gjson.Result {
		var finalPath = fmt.Sprintf(path, a...)
		log.Printf("Trying to get data at path: %s", finalPath)
		return gjson.Get(jsonString, finalPath)
	}
}

func getServerName(file os.FileInfo) string {
	var parts = strings.Split(file.Name(), ".")
	var serverName = parts[len(parts)-2]
	return serverName
}

func parseReport(file os.FileInfo) Report {
	var r Report

	// Get basic file info
	var filePath = fmt.Sprintf("%s/%s", latestPath, file.Name())
	r.filename = file.Name()
	r.serverName = getServerName(file)
	r.path = filePath

	log.Printf("Parsing report: %s", file.Name())

	// Get JSON contents
	var jsonString = string(utils.ReadFile(filePath))
	var getData = getterFactory(jsonString)

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
		// var referenceLinks

		var severity = c.Get("cvss2Severity").String()
		var cve = CVEInfo{
			id:          c.Get("cveID").String(),
			packageName: c.Get("affectedPackages.last.name").String(),
			severity:    severity,
			fixState:    c.Get("affectedPackages.last.fixState").String(),
			notFixedYet: c.Get("affectedPackages.last.notFixedYet").Bool(),
			title:       c.Get("cveContents.nvd.title").String(),
			summary:     c.Get("cveContents.nvd.summary").String(),
			// referenceLinks: referenceLinks,
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
