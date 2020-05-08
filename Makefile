install:
	go get ./src/

compile:
	rm -rf ./bin/*
	GOOS=freebsd GOARCH=386 go build -o bin/prometheus-vuls-exporter-freebsd-386 src/main.go
	GOOS=linux GOARCH=386 go build -o bin/prometheus-vuls-exporter-linux-386 src/main.go
	GOOS=windows GOARCH=386 go build -o bin/prometheus-vuls-exporter-windows-386 src/main.go

build:
	rm -rf ./bin/*
	go build -o bin/prometheus-vuls-exporter-386 src/main.go

run:
	go run ./src/main.go

version:
	@echo "Current version is `cat VERSION`"

version-bump:
	nano VERSION

git-release:
	git add VERSION
	git commit -m "Bumped version to `cat VERSION`"
	git tag `cat VERSION`
	git push
	git push --tags

github-upload:
	hub release create `cat VERSION` -a bin/prometheus-vuls-exporter-freebsd-386 -a bin/prometheus-vuls-exporter-windows-386 -a bin/prometheus-vuls-exporter-linux-386

release: version-bump compile git-release github-upload