// barbershop project main.go
package main

import (
	"fmt"
	"sync"
	"time"
)

const n = 4
const Customers_Num = 10000

var customers = 0
var barberToCustomer_ch = make(chan string)
var customerToBarber_ch = make(chan string)
var mutex sync.Mutex

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(Customers_Num)
	startTime := time.Now()

	go Barber()
	for i := 1; i < Customers_Num+1; i++ {
		go Customer(wg, i)

		time.Sleep(10 * time.Millisecond)
	}

	wg.Wait()
	fmt.Printf("Time Cost: %s.\n", time.Since(startTime))
}

func Barber() {
	for {
		barberToCustomer_ch <- "ready for next customer"

		message := <-customerToBarber_ch
		if message != "customer done" {
			fmt.Printf("error in barber")
		}
	}
}

func Customer(wg *sync.WaitGroup, i int) {
	defer wg.Done()
	mutex.Lock()
	if customers >= n {
		balk(i)
		mutex.Unlock()
	} else {
		customers++
		mutex.Unlock()
		message := <-barberToCustomer_ch
		if message == "ready for next customer" {
			getHairCut(i)
			customerToBarber_ch <- "customer done"

			mutex.Lock()
			customers--
			mutex.Unlock()
		} else {
			fmt.Printf("error in trying to contact barber")
			fmt.Printf(message)
		}
	}
}

func balk(i int) {
	fmt.Printf("Customer %d left because barbershop is full\n", i)
}

func getHairCut(i int) {
	time.Sleep(15 * time.Millisecond)
	fmt.Printf("Customer %d is getting hair cut\n", i)
}
