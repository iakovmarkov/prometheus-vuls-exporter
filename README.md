# prometheus-vuls-exporter

`prometheus-vuls-exporter` is a small Go application that allows scraping of Vuls reports by Prometheus.

## Exported metrics

This exporter exposes the following metrics:

```
# HELP reported_at Timestamp of last report time, in ms since Unix
# TYPE reported_at gauge
reported_at 1.58896003e+09
# HELP server_count Total count of servers reported
# TYPE server_count gauge
server_count 1
# HELP servers Server information, value represents amount of vulnerabilitites
# TYPE servers gauge
servers{cveID="CVE-2009-5080",hostname="test1.iakov.local",kernel_rebootRequired="false",kernel_release="4.15.0-91-generic",serverName="testServer1"} 1
servers{cveID="CVE-2009-5155",hostname="test1.iakov.local",kernel_rebootRequired="false",kernel_release="4.15.0-91-generic",serverName="testServer1"} 1
servers{cveID="CVE-2009-5155",hostname="test2.iakov.local",kernel_rebootRequired="false",kernel_release="4.15.0-91-generic",serverName="testServer2"} 1
# HELP vuln_count Total count of vulnerabilities, across all servers
# TYPE vuln_count gauge
vuln_count 2
# HELP vuln_severity Vulnerability count by severity
# TYPE vuln_severity gauge
vuln_severity{severity="high"} 0
vuln_severity{severity="low"} 1
vuln_severity{severity="medium"} 1
# HELP vulns Vulnerability information, value represents total amount of hits
# TYPE vulns gauge
vulns{cveID="CVE-2009-5080",fixState="",lastModified="2013-12-13T04:34:00Z",mitigation="",notFixedYet="false",packageName="",published="2011-06-30T15:55:00Z",severity="LOW",summary="The (1) contrib/eqn2graph/eqn2graph.sh, (2) contrib/grap2graph/grap2graph.sh, and (3) contrib/pic2graph/pic2graph.sh scripts in GNU troff (aka groff) 1.21 and earlier do not properly handle certain failed attempts to create temporary directories, which might allow local users to overwrite arbitrary files via a symlink attack on a file in a temporary directory, a different vulnerability than CVE-2004-1296.",title=""} 1
vulns{cveID="CVE-2009-5155",fixState="",lastModified="2019-03-25T17:29:00Z",mitigation="",notFixedYet="false",packageName="",published="2019-02-26T02:29:00Z",severity="MEDIUM",summary="In the GNU C Library (aka glibc or libc6) before 2.28, parse_reg_exp in posix/regcomp.c misparses alternatives, which allows attackers to cause a denial of service (assertion failure and application exit) or trigger an incorrect result by attempting a regular-expression match.",title=""} 1
```

## Installation

Download the latest release binary from (GitHub Releases page)[https://github.com/iakovmarkov/prometheus-vuls-exporter/releases]. Put it into your `/usr/bin` or anywhere on your `PATH`. Something like this, maybe:

    $ curl -Lo prometheus-vuls-exporter https://github.com/iakovmarkov/prometheus-vuls-exporter/releases/latest/prometheus-vuls-exporter-linux-386
    $ cp prometheus-vuls-exporter /usr/bin/

## Configuration

Configuration is possible via command line flags or environment variable. Possible options:

* `--reports_dir` or `REPORTS_DIR` - *mandatory* - path on a filesystem where the Vulns JSON reports are; must be readable by application
* `--basic_username` or `BASIC_USERNAME`; `--basic_password` or `BASIC_PASSWORD` - if both are set, enables HTTP Basic authorization
* `--address` or `ADDRESS` - where the server will listen for HTTP connections, defaults to `:8080`
* `--log_format` or `LOGFORMAT` - defines whether or not to output the date to log (`LONG` or `SHORT`, respectively) , default to `LONG`
* `--version` - print version and exit

## Example

Install `vuls` as described in (Vuls docs)[https://vuls.io/docs/en/install-with-vulsctl.html]. We will use `vulsctl` in the example.

Set environment variables:

    $ export REPORTS_DIR=/tmp/vuls_reports
    $ export SSH_DIR=/home/root/.ssh
    $ export VULS_DIR=/home/root/vulsctl

Once you have a config file prepared, run the script from `example/run_vuls.sh`. It will run Vuls scan, create JSON reports and set correct folder permissions. Ideally, this should be ran as a CRON job on your Vuls server.

Then you can run `prometheus-vuls-exporter` by running `example/run_exporter.sh`. This script will just run the `prometheus vuls-exporter` and make it collect reports from configured `REPORTS_DIR`. Point your browser to (localhost:8080/metrics)[http://localhost:8080/metrics] and observe.

New metrics will appear when you run the `example/run_vuls.sh` script again. It will generate new JSON metrics, and `prometheus-vuls-exporter` will pick them up next time it will gets asked to report metrics.

## Developing

### Go environment

First, validate that your `GOPATH` and `GOBIN` environmnet variables work. If you haven't set them up, you can try something like this:

    $ export GOPATH=/root/go
    $ export GOBIN=$GOPATH/bin
    $ export PATH=$PATH:$GOBIN

### Dependencies

After those are set, you can install dependencies:

    $ make install

### Run

You're set! Make your changes and run the application:

    $ make run

### Build

When you're done with your changes and ready to compile the binary, run this:

    $ make build

Your binary will be in the `bin` folder.