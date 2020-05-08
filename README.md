# prometheus-vuls-exporter

`prometheus-vuls-exporter` is a small Go application that allows scraping of Vuls reports by Prometheus.

## Installation

TBD

## Configuration

Configuration is possible via command line flags or environment variable. Possible options:

* `--reports_dir` or `REPORTS_DIR` - *mandatory* - path on a filesystem where the Vulns JSON reports are; must be readable by application
* `--basic_username` or `BASIC_USERNAME`; `--basic_password` or `BASIC_PASSWORD` - if both are set, enables HTTP Basic authorization
* `--address` or `ADDRESS` - where the server will listen for HTTP connections, defaults to `:8080`
* `--log_format` or `LOGFORMAT` - defines whether or not to output the date to log (`LONG` or `SHORT`, respectively) , default to `LONG`

## Example

Install `vuls` as described in (Vuls docs)[https://vuls.io/docs/en/install-with-vulsctl.html]. We will use `vulsctl` in the example.

Set environment variables:

    $ export REPORTS_DIR=/tmp/vuls_reports
    $ export SSH_DIR=/home/root/.ssh
    $ export VULS_DIR=/home/root/vulsctl

Once you have a config file prepared, run the script from `example/run_vuls.sh`. It will run Vuls scan, create JSON reports and set correct folder permissions. Ideally, this should be ran as a CRON job on your Vuls server.

Then you can run `prometheus-vuls-exporter` by running `example/run_exporter.sh`. This script will just run the `prometheus vuls-exporter` and make it collect reports from configured `REPORTS_DIR`. Point your browser to (localhost:8080/metrics)[http://localhost:8080/metrics], use `admin/s3cr3t` credentials and observe.

New metrics will appear when you run the `example/run_vuls.sh` script again. It will generate new JSON metrics, and `prometheus-vuls-exporter` will pick them up next time it will gets asked to report metrics.

## Developing

### Go environment

First, validate that your `GOPATH` and `GOBIN` environmnet variables work. If you haven't set them up, you can try something like this:

    $ export GOPATH=~/markov/go
    $ export GOBIN=$GOPATH/bin
    $ export PATH=$PATH:$GOBIN

### Dependencies

After those are set, you can install dependencies:

    $ go get ./src/

### Run

You're set! Make your changes and run the application:

    $ go run ./src/

### Build

When you're done with your changes and ready to compile and distribute the binary, run this:

    $ go build ./src/
