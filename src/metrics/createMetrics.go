package metrics

import (
	"fmt"
	"log"
	"sort"

	"../utils"
)

func createMetric(path string) func() float64 {
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
			count = count + len(report.cves)
		}

		return float64(count)
	}
}
