package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

func main() {
	rate := uint64(10000) // per second
	duration := 3 * time.Second
	targeter := NewFooBarTargeter(vegeta.Target{
		Method: "GET",
		URL:    "https://foo-bar.com/ping",
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration) {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("total number of requests: %d\n", metrics.Requests)
	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf("percentage of non-error responses: %g\n", metrics.Success)
}

// NewFooBarTargeter returns a Targeter where we calculate a custom header
//  for each request.
func NewFooBarTargeter(target vegeta.Target) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		*tgt = target
		tgt.Header = http.Header{}

		tgt.Header.Add("rand", strconv.Itoa(rand.Intn(100)))
		log.Printf("%+v", tgt)
		return nil
	}
}
