package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/trace"
	"time"
)

func main() {
	trace.Start(os.Stderr)

	c := gen("AAPL")
	outBound := square(c)

	for val := range outBound {
		fmt.Printf("received val: %s\n", val)
	}

	// small sleeps are required before stopping the trace
	// to ensure the output has been collected
	time.Sleep(time.Millisecond * 100)
	trace.Stop()
}

// copied from https://blog.golang.org/pipelines
// converts variadic or slice of ints to distinct values in a pipeline
func gen(nums ...string) <-chan string {
	out := make(chan string)
	go func() {
		for _, n := range nums {
			time.Sleep(time.Millisecond * 100)
			out <- n
		}
		// small sleep before terminating the goroutine to ensure
		// trace output is collected
		time.Sleep(time.Millisecond * 50)
		// inform the pipeline that's all we've got
		close(out)
	}()

	return out
}

func square(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for n := range in {
			out <- fetchData(n)[0].Date
		}
		// small sleep before terminating the goroutine to ensure
		// trace output is collected
		time.Sleep(time.Millisecond * 50)
		close(out)
	}()
	return out
}

type day struct {
	Date   string  `json:"date"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
	VWAP   float64 `json:"vwap"`
}

// TODO: return an error?
func fetchData(symbol string) []day {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.iextrading.com/1.0/stock/%s/chart/1y", symbol), nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// prevent memory leak
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var data []day

	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}

	return data
}

// function to perform some bogus (but CPU intensive) calculations on the data
func crunchNumbers(tradingData []day) string {
	var total float64
	for i := 0; i < 1000; i++ {
		for _, singleDay := range tradingData {
			// total = total + ((singleDay.Close - singleDay.Volume*singleDay.VWAP) / (singleDay.Close * singleDay.Close))
			total = singleDay.Close
		}
	}

	if total > 10000000 {
		return tradingData[0].Date
	}
	return ""
}
