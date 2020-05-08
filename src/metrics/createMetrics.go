package metrics

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"

	"../utils"
)

type JSONUnstruct = map[string]interface{}

type Report struct {
	Filename  string
	Path      string
	Hostname  string
	VulnCount int
}

func parseReport(file os.FileInfo) Report {
	var filePath = fmt.Sprintf("%s/%s", latestPath, file.Name())

	var r Report
	r.Filename = file.Name()
	r.Path = filePath

	log.Printf("Parsing report: %s", file.Name())
	var content = utils.ReadFile(filePath)

	var data JSONUnstruct
	json.Unmarshal([]byte(content), &data)

	r.VulnCount = len(data["scannedCves"].(JSONUnstruct))

	return r
}

var (
	reportsPath string
	reportsDir  []os.FileInfo
	latestPath  string
	latestDir   []os.FileInfo
)

func CreateMetric(path string) func() float64 {
	reportsPath = path
	reportsDir = utils.ReadDir(reportsPath)
	log.Printf("Reports folder configured: %s", reportsPath)

	return func() float64 {
		var count = 0
		reportsDir = utils.ReadDir(reportsPath)
		reportsDir = utils.FilterCurrent(reportsDir)

		// Sort directories by their Modified Date, descending
		sort.Slice(reportsDir, func(i, j int) bool {
			return reportsDir[i].ModTime().Unix()-reportsDir[j].ModTime().Unix() > 0
		})

		// Get very first element in the sorted dir
		latestPath = fmt.Sprintf("%s/%s", reportsPath, reportsDir[0].Name())
		latestDir = utils.ReadDir(latestPath)

		// Get only JSON files
		var reports = utils.AcceptJSON(latestDir)
		for _, file := range reports {
			var report = parseReport(file)
			count = count + report.VulnCount
		}

		return float64(count)
	}
}
