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
	"bytes"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/thales-e-security/erfclient"
)

const (
	uriEnvKey       = "TEST_URI"
	refreshEnvKey   = "TEST_REFRESH"
	tokenFileEnvKey = "TOKEN_FILE"
	logFileEnvKey   = "LOG_FILE"
	sleepEnvKey     = "SLEEP"
)

func main() {
	host := getEnv(uriEnvKey, "http://localhost:8080")

	refreshStr := getEnv(refreshEnvKey, "5")
	refresh, err := strconv.Atoi(refreshStr)
	if err != nil {
		log.Fatal(err.Error())
	}

	sleepStr := getEnv(sleepEnvKey, "1")
	sleep, err := strconv.Atoi(sleepStr)
	if err != nil {
		log.Fatal(err.Error())
	}

	tokenFile := getEnv(tokenFileEnvKey, "tokenfile.txt")
	recordFile := getEnv(logFileEnvKey, "client.log")
	logFile, err := os.Create(recordFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.SetOutput(logFile)
	log.Println("New test run")

	erf, err := erfclient.New(tokenFile, uint(refresh))
	client := &http.Client{}

	for {
		req, err := http.NewRequest(http.MethodPost, host, bytes.NewBufferString("Hello, World!"))
		if err != nil {
			log.Fatal(err)
		}

		token, err := erf.Token()
		if err != nil {
			log.Fatal("Failed to get token: ", err)
		}

		req.Header.Add("ERF", string(token))
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Failed to send to web server: ", err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Println("Warning, bad error code: ", resp.StatusCode)
			continue
		}

		log.Printf("Sent message, ERF was: %s", string(token))
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}

func getEnv(env, defaultValue string) string {
	s, found := os.LookupEnv(env)
	if !found {
		return defaultValue
	}
	return s
}
