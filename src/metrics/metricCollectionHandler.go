package metrics

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"../utils"
)

func collectReports(reportPath string) []Report {
	var reports []Report

	// Read reports from path, remove `current` link
	reportsDir = utils.ReadDir(reportsPath)
	reportsDir = utils.FilterCurrent(reportsDir)

	// Sort directories by their Modified Date, descending
	sort.Slice(reportsDir, func(i, j int) bool {
		return reportsDir[i].ModTime().Unix()-reportsDir[j].ModTime().Unix() > 0
	})

	// Get very first element in the sorted reports
	latestPath = fmt.Sprintf("%s/%s", reportsPath, reportsDir[0].Name())
	latestDir = utils.ReadDir(latestPath)

	// Get only JSON files
	reportFiles := utils.AcceptJSON(latestDir)
	for _, file := range reportFiles {
		report := parseReport(file)
		reports = append(reports, report)
	}

	return reports
}

func MetricCollectionHandler(path string) func(http.HandlerFunc) http.HandlerFunc {
	reportsPath = path
	reportsDir = utils.ReadDir(reportsPath)
	log.Printf("Reports folder configured: %s", reportsPath)

	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Gather data from report files
			reports = collectReports(reportsPath)

			// Serve next request
			h.ServeHTTP(w, r)
		}
	}
}
