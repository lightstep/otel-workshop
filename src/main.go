package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

)

func main() {

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(rootHandler))
	mux.Handle("/favicon.ico", http.NotFoundHandler())
	mux.Handle("/fib", http.HandlerFunc(fibHandler))
	mux.Handle("/fibinternal", http.HandlerFunc(fibHandler))
	// mux.Handle("/metrics", prom)
	os.Stderr.WriteString("Initializing the server...\n")

	err := http.ListenAndServe("127.0.0.1:3000", mux)
	if err != nil {
		log.Fatalf("Could not start web server: %s", err)
	}
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	_ = dbHandler("foo")

	fmt.Fprintf(w, "Your server is live! Try to navigate to: http://localhost:3000/fib?i=6")
}

func fibHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	var i int
	if len(req.URL.Query()["i"]) != 1 {
		err = fmt.Errorf("Wrong number of arguments.")
	} else {
		i, err = strconv.Atoi(req.URL.Query()["i"][0])
	}
	if err != nil {
		fmt.Fprintf(w, "Couldn't parse index '%s'.", req.URL.Query()["i"])
		w.WriteHeader(503)
		return
	}
	ret := 0
	failed := false

	if i < 2 {
		ret = 1
	} else {
		// Call /fib?i=(n-1) and /fib?i=(n-2) and add them together.
		var mtx sync.Mutex
		var wg sync.WaitGroup
		client := http.DefaultClient
		for offset := 1; offset < 3; offset++ {
			wg.Add(1)
			go func(n int) {
				err := func() error {

					url := fmt.Sprintf("http://127.0.0.1:3000/fibinternal?i=%d", n)
					req, _ := http.NewRequestWithContext(req.Context(), "GET", url, nil)
					res, err := client.Do(req)
					if err != nil {
						return err
					}
					body, err := ioutil.ReadAll(res.Body)
					res.Body.Close()
					if err != nil {
						return err
					}
					resp, err := strconv.Atoi(string(body))
					if err != nil {
						return err
					}
					mtx.Lock()
					defer mtx.Unlock()
					ret += resp
					return err
				}()
				if err != nil {
					if !failed {
						w.WriteHeader(503)
						failed = true
					}
					fmt.Fprintf(w, "Failed to call child index '%s'.\n", n)
				}
				wg.Done()
			}(i - offset)
		}
		wg.Wait()
	}
	fmt.Fprintf(w, "%d", ret)
}

func dbHandler(color string) int {
	// Pretend we talked to a database here.
	return 0
}
