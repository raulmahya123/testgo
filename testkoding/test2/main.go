package main

import (
	"sync"
	"testing"
	"time"
)

// Put your logic and code inside this function to optimize the app.
func Runner() map[int]string {
	results := make(map[int]string)
	var wg sync.WaitGroup

	// Use a channel to communicate results back from the Goroutines
	resultChannel := make(chan Result, worker)

	// Loop through workers and execute concurrently
	for i := 1; i <= worker; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			res, err := mockGetData(id)
			if err != nil {
				return
			}
			resultChannel <- *res
		}(i)
	}

	// Wait for all Goroutines to finish
	wg.Wait()
	close(resultChannel)

	// Collect results from the channel
	for res := range resultChannel {
		results[res.ID] = res.Title
	}

	return results
}

// ===== DO NOT EDIT. =====
type (
	Result struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}
)

const worker = 3

var (
	expectedWorker                = 0
	expected       map[int]string = map[int]string{
		1: "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
		2: "qui est esse",
		3: "ea molestias quasi exercitationem repellat qui ipsa sit aut",
	}
)

func TestCaseParalelUniverse(t *testing.T) {
	start := time.Now()
	results := Runner()
	lat := time.Since(start).Milliseconds()

	switch {
	case lat >= 2500:
		t.Fail()

	case assertEqual(results):
		t.Fail()

	case expectedWorker < worker:
		t.Fail()
	}
}

func assertEqual(results map[int]string) bool {
	for k, v := range results {
		if expected[k] != v {
			return true
		}
	}

	return false
}

func mockGetData(id int) (*Result, error) {
	expectedWorker++
	time.Sleep(2 * time.Second)

	result := Result{
		ID:    id,
		Title: expected[id],
	}
	expected[id] = result.Title
	return &result, nil
}

func main() {
	testSuite := []testing.InternalTest{
		{
			Name: "TestCaseParalelUniverse",
			F:    TestCaseParalelUniverse,
		},
	}

	testing.Main(nil, testSuite, nil, nil)
}

// ===== DO NOT EDIT. =====
