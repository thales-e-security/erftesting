// Copyright 2018 Thales UK Limited
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
// documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
// Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/thales-e-security/erfserver"
)

const (
	portEnvKey = "PORT"
)

var erfServer erfserver.ERFServer

func writeError(w http.ResponseWriter, code int, error string) {
	w.WriteHeader(code)
	w.Write([]byte(error))
}

func msgHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusBadRequest, "Only POST supported")
		return
	}

	erf := r.Header.Get("ERF")
	if erf == "" {
		writeError(w, http.StatusBadRequest, "Missing ERF header")
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	erfServer.Append([]byte(erf), string(content), time.Now())
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusBadRequest, "Only GET supported")
		return
	}

	results := erfServer.OperationsByClient()

	w.Header().Add("content-type", "text/html")

	w.Write([]byte("<html><table><tr><th>Client ID</th><th>Operation</th><th>Count</th></tr>"))

	for id, opMap := range results {
		for op, count := range opMap {
			w.Write([]byte(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%d</td></tr>", id, op, count)))
		}
	}

	w.Write([]byte("</table></html>"))
}

func main() {
	var port int
	var err error

	portString, found := os.LookupEnv(portEnvKey)
	if !found {
		port = 8080
	} else {
		port, err = strconv.Atoi(portString)
		if err != nil {
			panic(err)
		}
	}

	erfServer = erfserver.NewInMemory()

	http.HandleFunc("/", msgHandler)
	http.HandleFunc("/results", resultsHandler)
	log.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
