package metrics

import (
	"fmt"
	"os"
	"prometheus-vuls-exporter/utils"
	"strings"

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

//  removeDuplicateValues removes duplicates from a slice of string
func removeDuplicateValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	results := []string{}

	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			results = append(results, entry)
		}
	}
	return results
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
		var severity string
		cvssSeverities := c.Get("cveContents.@values.@flatten.#.cvss2Severity")

		var severitiesSlice []string
		for _, sev := range cvssSeverities.Array() {
			severitiesSlice = append(severitiesSlice, sev.String())
		}

		uniqueSeverities := removeDuplicateValues(severitiesSlice)
		if len(uniqueSeverities) > 0 {
			severity = strings.ToLower(uniqueSeverities[0])
		} else {
			severity = "UNKNOWN"
		}

    // Get affected package or CPE URI
    var packageName string
		var fixState string
		var notFixedYet bool
    if c.Get("affectedPackages").Exists() {
        packageName = c.Get("affectedPackages.0.name").String()
				fixState    = c.Get("affectedPackages.0.fixState").String()
				notFixedYet = c.Get("affectedPackages.0.notFixedYet").Bool()
    } else if c.Get("cpeURIs").Exists() {
        packageName = c.Get("cpeURIs.0").String()
				fixState    = ""
				notFixedYet = false
    } else {
				// Fallback: No affectedPackages or cpeURIs available
				packageName = ""
				fixState    = ""
				notFixedYet = false
		}

		cve := CVEInfo{
			id:           c.Get("cveID").String(),
			packageName:  packageName,
			severity:     severity,
			fixState:     fixState,
			notFixedYet:  notFixedYet,
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
