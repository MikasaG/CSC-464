package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Oxygen struct{}
type Hydrogen struct{}
type Water struct {
	id int
	o  Oxygen
	h1 Hydrogen
	h2 Hydrogen
}

const Water_num = 10000

var oxygen_channel = make(chan Oxygen, Water_num)
var hydrogen_channel = make(chan Hydrogen, Water_num*2)

func main() {
	Oxygen_num := 0
	Hydrogen_num := 0
	startTime := time.Now()
	r := rand.New(rand.NewSource(startTime.UnixNano()))

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go bonder(wg)

	for i := 0; i < Water_num*3; i++ {
		random := r.Intn(2)
		if (random != 0 && Hydrogen_num != Water_num*2) || Oxygen_num == Water_num {
			hydrogen_channel <- Hydrogen{}
			Hydrogen_num++
			//fmt.Printf("No. %d Hydrogen arrive\n", Hydrogen_num)
		} else {
			oxygen_channel <- Oxygen{}
			Oxygen_num++
			//fmt.Printf("No. %d Oxygen arrive\n", Oxygen_num)
		}
	}
	wg.Wait()
	fmt.Printf("--- %s ---\n", time.Since(startTime))
}

func bonder(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < Water_num; i++ {
		o := <-oxygen_channel
		h1 := <-hydrogen_channel
		h2 := <-hydrogen_channel
		water := Water{
			i + 1,
			o,
			h1,
			h2,
		}
		fmt.Printf("No.%d water bonded.\n", water.id)
	}
}
