// smoker project main.go
package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const Savage_Num = 1000000

var ch = make(chan string)
var servings int = 10

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	startTime := time.Now()

	go Cook(wg)
	go Savage(wg)

	wg.Wait()
	fmt.Printf("Time Cost: %s.\n", time.Since(startTime))
}

func Cook(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < Savage_Num/10-1; i++ {
		message := <-ch
		if message == "empty" {
			putServingsInPot()
		}
		ch <- "full"
	}
}

func Savage(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i < Savage_Num+1; i++ {
		if servings == 0 {
			ch <- "empty"
			message := <-ch
			if message == "full" {
				servings = 10
			} else {
				fmt.Printf("error in Cook!\n")
			}
		}
		servings--
		eat("Savege " + strconv.Itoa(i))
	}
}

func eat(s string) {
	fmt.Printf("%s is eating!\n", s)
}

func putServingsInPot() {
	fmt.Printf("The Cook has put 10 servings in pot.\n")
}
