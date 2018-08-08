package testing

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/thales-e-security/erfclient"
	"github.com/thales-e-security/erfserver"
)

const (
	initialTestClients = 30
	tokenRefresh       = 1
)

var results map[string]map[string]int
var s erfserver.ERFServer
var testEntries int

func BenchmarkOpsPerClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		results = s.OperationsByClient()
	}

	fmt.Println(len(results))
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	if num, found := os.LookupEnv("TESTNUM"); found {
		var err error
		testEntries, err = strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
	} else {
		testEntries = 10000
		fmt.Println("TESTNUM env var not found, defaulting to ", testEntries)
	}

	dir, err := ioutil.TempDir("", "erftesting")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	s = erfserver.NewInMemory()

	var clients []erfclient.ERFClient

	for i := 0; i < initialTestClients; i++ {
		c, err := erfclient.New(path.Join(dir, fmt.Sprintf("tokenfile%d", len(clients))), tokenRefresh)
		if err != nil {
			panic(err)
		}
		clients = append(clients, c)
	}

	// These thresholds affect how often clones are created and new clients added. They
	// are dependent upon the size of the test entries created.
	cloneThreshold := int(math.Ceil(float64(testEntries) / float64(2)))
	addThreshold := int(math.Ceil(float64(testEntries) / float64(10)))

	// Allow 30% growth in clients
	maxClients := int(math.Ceil(float64(initialTestClients) * 1.3))

	fmt.Printf("Clone threshold: %d\nAdd Threshold: %d\n", cloneThreshold, addThreshold)

	for i := 0; i < testEntries; i++ {
		for n, c := range clients {
			token, err := c.Token()
			if err != nil {
				panic(err)
			}

			s.Append(token, "testop", time.Now())

			// roll the dice, do we clone this?
			if len(clients) < maxClients && rand.Intn(cloneThreshold) == 0 {
				fmt.Println("Making clone")
				tokenFileClone := path.Join(dir, fmt.Sprintf("tokenfile%d", len(clients)))
				err = copy(path.Join(dir, fmt.Sprintf("tokenfile%d", n)), tokenFileClone)
				if err != nil {
					panic(err)
				}

				clone, err := erfclient.New(tokenFileClone, tokenRefresh)
				if err != nil {
					panic(err)
				}
				clients = append(clients, clone)
			}
		}

		// roll the dice, do we add another?
		if len(clients) < maxClients && rand.Intn(addThreshold) == 0 {
			fmt.Println("Adding another client")
			c, err := erfclient.New(path.Join(dir, fmt.Sprintf("tokenfile%d", len(clients))), tokenRefresh)
			if err != nil {
				panic(err)
			}
			clients = append(clients, c)
		}
	}

	fmt.Println("Total clients at end: ", len(clients))
}
