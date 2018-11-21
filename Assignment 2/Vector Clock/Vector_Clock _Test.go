package main

import (
	"fmt"
	"sync"
)

func TestMessageSending() {
	threads := []string{"t1", "t2", "t3"}
	vc1, _ := newVectorClock(0, "t1", threads)
	vc2, _ := newVectorClock(1, "t2", threads)
	vcBad, _ := newVectorClock(0, "t3", []string{"invalid threads information"})
	vc1.vector[threads[0]] = 5
	vc1.vector[threads[1]] = 3
	vc1.vector[threads[2]] = 71

	vc2.vector[threads[0]] = 1
	vc2.vector[threads[1]] = 23
	vc2.vector[threads[2]] = 100

	var ch_2_to_1 = make(chan map[string]int)
	wg := sync.WaitGroup{}
	wg.Add(2)
	//say thread 2 send a message to thread 1, then thread 1 update it's vector clock.

	//thread 2
	go func() {
		defer wg.Done()
		vc2.sendMessage(vc1, ch_2_to_1)

	}()
	//thread 1
	go func() {
		defer wg.Done()
		vc1.receiveMessage(vc2, ch_2_to_1)

	}()
	wg.Wait()
	if vc1.vector[threads[0]] != 6 && vc1.vector[threads[1]] != 24 && vc1.vector[threads[2]] != 101 {
		fmt.Printf("Expected values of %d,%d,%d for thread1,thread2,thread3 respectively in vc1, but got %d,%d,%d\n", 5, 23, 100, vc1.vector[threads[0]], vc1.vector[threads[1]], vc1.vector[threads[2]])
	} else {
		fmt.Print("values in vc1 are right\n")
	}
	if vc2.vector[threads[0]] != 2 && vc2.vector[threads[1]] != 24 && vc2.vector[threads[2]] != 101 {
		fmt.Printf("Expected values of %d,%d,%d for thread1,thread2,thread3 respectively in vc2, but got %d,%d,%d\n", 2, 23, 100, vc2.vector[threads[0]], vc2.vector[threads[1]], vc2.vector[threads[2]])
	} else {
		fmt.Print("values in vc2 are right\n")
	}
	//say a bad thread send a message to thread 1, then thread 1 try to update it's vector clock, but detect an error
	err := vc1.update(vcBad.vector)
	if err == nil {
		fmt.Printf("Expected error for invalid update but none was returned\n")
	} else {
		fmt.Print("TestMessageSending Passed.\n\n")
	}
}

func TestIncreaseOne() {
	threads := []string{"t1", "t2", "t3"}
	vc, _ := newVectorClock(0, "t1", threads)
	vc.increaseOne()
	vc.increaseOne()
	if vc.vector[threads[0]] != 2 {
		fmt.Printf("Unexpected value for thread 1 in vector clock: Expected %d but got %d", 3, vc.vector[threads[0]])
	} else if vc.vector[threads[1]] != 0 {
		fmt.Printf("Unexpected value for thread 1 in vector clock: Expected %d but got %d", 0, vc.vector[threads[1]])
	} else if vc.vector[threads[2]] != 0 {
		fmt.Printf("Unexpected value for thread 2 in vector clock: Expected %d but got %d", 0, vc.vector[threads[2]])
	} else {
		fmt.Printf("TestIncreaseOne Passed.\n\n")
	}
}

func TestNewVectorClock() {
	threads := []string{"t1", "t2", "t3"}
	vc, err := newVectorClock(0, "t1", threads)
	_, err2 := newVectorClock(3, "t3", threads)
	val := vc.vector[threads[0]]
	if err != nil {
		fmt.Print(err)
	} else if _, ok := vc.vector[threads[0]]; !ok {
		fmt.Print("Failed to find thread 1 in vector clock")
	} else if val != 0 {
		fmt.Printf("Unexpected value for node1 in vector clock: Expected %d but got %d", 0, val)
	} else if err2 == nil {
		fmt.Printf("Expected error for invalid index but none was returned")
	} else {
		fmt.Printf("TestNewVectorClock Passed.\n\n")
	}
}

func Scenario_1() {
	fmt.Printf("Scenario 1 Test begin:.\n")
	threads := []string{"1", "2", "3"}
	vc1, _ := newVectorClock(0, "1", threads)
	vc2, _ := newVectorClock(1, "2", threads)
	vc3, _ := newVectorClock(2, "3", threads)
	var ch_1_to_2 = make(chan map[string]int)
	var ch_1_to_3 = make(chan map[string]int)
	var ch_2_to_1 = make(chan map[string]int)
	var ch_3_to_2 = make(chan map[string]int)
	wg := sync.WaitGroup{}
	wg.Add(3)
	//Thread 1
	go func() {
		vc1.increaseOne()
		vc1.sendMessage(vc2, ch_1_to_2)
		vc1.increaseOne()
		vc1.receiveMessage(vc2, ch_2_to_1)
		vc1.sendMessage(vc3, ch_1_to_3)
		wg.Done()
	}()

	//Thread 2
	go func() {
		vc2.receiveMessage(vc3, ch_3_to_2)
		vc2.receiveMessage(vc1, ch_1_to_2)
		vc2.sendMessage(vc1, ch_2_to_1)
		wg.Done()
	}()

	//Thread 3
	go func() {
		vc3.sendMessage(vc2, ch_3_to_2)
		vc3.increaseOne()
		vc3.receiveMessage(vc1, ch_1_to_3)
		wg.Done()
	}()
	wg.Wait()
	if vc3.vector[threads[0]] != 5 || vc3.vector[threads[1]] != 3 || vc3.vector[threads[2]] != 3 {
		fmt.Printf("Incorrect values: Expected 5,3,3 but got %d,%d,%d\n\n", vc1.vector[threads[0]], vc1.vector[threads[1]], vc1.vector[threads[2]])
	} else {
		fmt.Printf("Thread 1 has a resulting vector of [5,3,3],which is as expected\n\n")
	}

}

// This scenarion is taken from https://en.wikipedia.org/wiki/Vector_clock
func Scenario_2() {
	fmt.Printf("Scenario 2 Test begin:.\n")
	threads := []string{"A", "B", "C"}
	vc1, _ := newVectorClock(0, "A", threads)
	vc2, _ := newVectorClock(1, "B", threads)
	vc3, _ := newVectorClock(2, "C", threads)
	var ch_1 = make(chan map[string]int)
	var ch_2 = make(chan map[string]int)
	var ch_3 = make(chan map[string]int)
	wg := sync.WaitGroup{}
	wg.Add(3)
	//Thread A
	go func() {
		vc1.receiveMessage(vc2, ch_2)
		vc1.sendMessage(vc2, ch_1)
		vc1.receiveMessage(vc3, ch_3)
		vc1.receiveMessage(vc3, ch_3)
		wg.Done()
	}()

	//Thread B
	go func() {
		vc2.receiveMessage(vc3, ch_3)
		vc2.sendMessage(vc1, ch_2)
		vc2.sendMessage(vc3, ch_2)
		vc2.receiveMessage(vc1, ch_1)
		vc2.sendMessage(vc3, ch_2)
		wg.Done()
	}()

	//Thread C
	go func() {
		vc3.sendMessage(vc2, ch_3)
		vc3.receiveMessage(vc2, ch_2)
		vc3.sendMessage(vc1, ch_3)
		vc3.receiveMessage(vc2, ch_2)
		vc3.sendMessage(vc1, ch_3)
		wg.Done()
	}()
	wg.Wait()
	if vc1.vector[threads[0]] != 4 || vc1.vector[threads[1]] != 5 || vc1.vector[threads[2]] != 5 {
		fmt.Printf("Incorrect values: Expected 4,5,5 but got %d,%d,%d\n", vc1.vector[threads[0]], vc1.vector[threads[1]], vc1.vector[threads[2]])
	} else {
		fmt.Printf("Thread A has a resulting vector of [4,5,5],which is as expected.\n\n")
	}

}

func main() {
	TestMessageSending()
	TestIncreaseOne()
	TestNewVectorClock()
	Scenario_1()
	Scenario_2()
}
