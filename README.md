# prometheus-vuls-exporter

`prometheus-vuls-exporter` is a small Go application that allows scraping of Vuls reports by Prometheus.

## Installation

TBD

## Configuration

Configuration is possible via command line flags or environment variable. Possible options:

* `--address` or `ADDRESS` - where the server will listen for HTTP connections, defaults to `:8080`
* `--logFormat` or `LOGFORMAT` - defines whether or not to output the date to log (`LONG` or `SHORT`, respectively) , default to `LONG`

## Example

TBD

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
