build:
	GOOS=freebsd GOARCH=386 go build -o bin/prometheus-vuls-exporter-freebsd-386 src/main.go
	GOOS=linux GOARCH=386 go build -o bin/prometheus-vuls-exporter-linux-386 src/main.go
	GOOS=windows GOARCH=386 go build -o bin/prometheus-vuls-exporter-windows-386 src/main.go

run:
	go run ./src/main.go

version:
	@echo "Current version is `cat VERSION`"

git-release:
	git add VERSION
	git commit -m "Bumped version to `cat VERSION`"
	git tag `cat VERSION`
	git push
	git push --tags