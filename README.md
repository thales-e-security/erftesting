# erftesting

Contains example and benchmark code for the two ephemeral random fingerprint (ERF) projects, [erfserver](https://github.com/thales-e-security/erfserver) and [erfclient](https://github.com/thales-e-security/erfclient).

An example client application is provided in `client/main.go`, which periodically pushes requests to the example server application in `server/main.go`. Each request includes an ERF, which allows the server to identify the client and count operations.

The server offers a simple web interface to view the count of current clients.

## Building the Tests

You will need to [install Go](https://golang.org/doc/install) and [install dep](https://github.com/golang/dep#installation) if you haven't already got them.


```bash
go get github.com/thales-e-security/erftesting
cd $GOPATH/src/github.com/thales-e-security/erftesting
dep ensure --vendor-only
```


If you have make, then run:

```bash
make
```
Otherwise just build the executables manually:

```bash
go build -o testclient ./client/...
go build -o testserver ./server/...
```


## Running the Tests

The test server listens on a port for incoming requests from the client(s). The default port is 8080, but you can override this by setting an environment variable:

```bash
PORT=5678 ./testserver 
  2018/08/08 12:31:36 Starting server on port 5678
```

Fire up the test client(s) in separate terminals. Each client will POST a message to the server on a regular basis, refreshing its ephemeral random fingerprint (ERF) every now and then. You can adjust behaviour with these settings:

| Environment Variable | Meaning | Default |
|----------------------|---------|---------|
| `TEST_URI`             | The URI that the server is listening on. | `http://localhost:8080` |
| `TEST_REFRESH`             | How often the ERF is refreshed (seconds). | `5` |
| `SLEEP`             | How long to sleep between sending messages (seconds). | `1` |
| `TOKEN_FILE`             | Where to store the ERF token. | `tokenfile.txt` (in working directory) |
| `LOG_FILE`             | Where to store log output. | `client.log` (in working directory) |

You can view the stats collected by the server by visiting `/results` (e.g. [http://localhost:8080/results](http://localhost:8080/results)).

## Benchmarking

A benchmarking test is also included, which calculates the average time needed to run the `OperationsByClient` method (see https://godoc.org/github.com/thales-e-security/erfserver). This is the method that traverses the entire graph of recorded operations and summarises how many unique clients there are and how many operations each client performed.

The test starts with 30 (virtual) clients and generates 10,000 entries per client (by default). This figure can be changed with the `TESTNUM` environment variable. During the test, clients are randomly spawned and cloned up to a maximum of 30% of the initial total clients. This means most runs will conclude with 39 active clients. The figure of 30 can be changed quite easily in the source.

If you would like a great primer on benchmarking in Go, read [this page](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go) by Dave Cheney.

To run the test:

```bash
cd $GOPATH/src/github.com/thales-e-security/erftesting
go test -bench=.
```

The output will look something like this:

```
TESTNUM env var not found, defaulting to  10000
Clone threshold: 5000
Add Threshold: 1000
Adding another client
Adding another client
Making clone
Making clone
Making clone
Making clone
Making clone
Making clone
Making clone
Total clients at end:  39
39
goos: linux
goarch: amd64
pkg: github.com/thales-e-security/erftesting
BenchmarkOpsPerClient-8   	39
      10	 165981187 ns/op
PASS
ok  	github.com/thales-e-security/erftesting	7.565s
```

So on my machine, it takes 0.16s for the server to process the entire graph for 39 clients, with 10,000 operations each.