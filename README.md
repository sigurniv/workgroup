workgroup extends golang.org/x/sync/errgroup with possibility to retrieve result and error for each working goroutine
 
 ```go
package main

import (
	"net/http"
	"github.com/sigurniv/workgroup"
	"fmt"
)

func main() {
	g := workgroup.New()
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.213sadasd.com/",
	}
	for _, url := range urls {
		// Launch a goroutine to fetch the URL.
		url := url // https://golang.org/doc/faq#closures_and_goroutines

		// label is a key in returning maps
		g.Go(url, func() (interface{}, error) {
			// Fetch the URL.
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
			}
			return nil, err
		})
	}

	// Wait for all HTTP fetches to complete.
	// result and error for each goroutine can be retrieved using label provided to g.Go
	results, errors := g.Wait()

	fmt.Println("Results")
	for k, v := range results {
		fmt.Printf("%s : %v\n", k, v)
	}

	fmt.Println("Errors")
	for k, v := range errors {
		fmt.Printf("%s : %v\n", k, v)
	}
}
 ```