install:
	go mod vendor

compile:
	rm -rf ../bin/*
	GOOS=freebsd GOARCH=386 go build -o bin/prometheus-vuls-exporter-freebsd-386
	GOOS=linux GOARCH=386 go build -o bin/prometheus-vuls-exporter-linux-386
	GOOS=windows GOARCH=386 go build -o bin/prometheus-vuls-exporter-windows-386

build:
	rm -rf ../bin/*
	go build -o ../bin/prometheus-vuls-exporter-386

run:
	go run

version:
	@echo "Current version is `cat ../VERSION`"

version-bump:
	nano ../VERSION

git-release:
	git add ../VERSION
	git commit -m "Bumped version to `cat ../VERSION`"
	git tag `cat ../VERSION`
	git push
	git push --tags

github-upload:
	gh release create `cat ../VERSION` bin/prometheus-vuls-exporter-freebsd-386 bin/prometheus-vuls-exporter-windows-386 bin/prometheus-vuls-exporter-linux-386

release: version-bump compile git-release github-upload