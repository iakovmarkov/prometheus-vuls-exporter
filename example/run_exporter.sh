#! /bin/bash
set -e

echo Running exporter with Go

go run ./src/ --reports_dir=$REPORTS_DIR