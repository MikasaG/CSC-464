// smoker project main.go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type component struct {
	item string
}

const MAX_Cigarette = 10000

var Smoker_channel = make(chan component, 2)
var Agent_channel = make(chan string)

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	startTime := time.Now()

	go agent(wg)
	go smoker(wg)

	wg.Wait()
	fmt.Printf("Time Cost: %s.\n", time.Since(startTime))
}

func smoker(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < MAX_Cigarette; i++ {
		c1 := <-Smoker_channel
		c2 := <-Smoker_channel
		if c1.item == "paper" && c2.item == "match" {
			fmt.Printf("Smoker A is smoking No.%d Cigarette\n", i+1)
		} else if c1.item == "tobacco" && c2.item == "match" {
			fmt.Printf("Smoker B is smoking No.%d Cigarette\n", i+1)
		} else if c1.item == "paper" && c2.item == "tobacco" {
			fmt.Printf("Smoker C is smoking No.%d Cigarette\n", i+1)
		} else {
			fmt.Printf("error! component don't match!")
		}
		Agent_channel <- "proceed"
	}
}

func agent(wg *sync.WaitGroup) {
	defer wg.Done()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < MAX_Cigarette; i++ {
		random := r.Intn(2)
		if random == 0 {
			Smoker_channel <- component{"paper"}
			Smoker_channel <- component{"match"}
		} else if random == 1 {
			Smoker_channel <- component{"tobacco"}
			Smoker_channel <- component{"match"}
		} else {
			Smoker_channel <- component{"paper"}
			Smoker_channel <- component{"tobacco"}
		}
		proceed := <-Agent_channel
		if proceed != "proceed" {
			fmt.Printf("error in proceed!")
		}
	}
}
