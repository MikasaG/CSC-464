// Vector_Clock.go   Chang Gong  V00898803
package main

import (
	"fmt"
	"sync"
)

//The vector clock struct contains three components.
//Values:
//thread(string): the name of the thread which own this vector clock
//vector(map[string]int): the key part, vector holds its last record of the counter informations of all other threads.
//mutex: ensure only one thread modify the vector clock at one time.
type Vector_Clock struct {
	thread string
	vector map[string]int
	mutex  sync.Mutex
}

//This method will take three parameters and return a new Vector_Clock
//Parameters:
//index(int): the index of the thread
//thread_name(string):name of the thread
//threads([]string):an string array, containing all thread names.
func newVectorClock(index int, thread_name string, threads []string) (vc Vector_Clock, e error) {
	if index < 0 || index >= len(threads) {
		e = fmt.Errorf("The index value of the new vector clock is invalid\n")
		return
	}
	vector := make(map[string]int)
	for _, thread := range threads {
		vector[thread] = 0
	}
	vc = Vector_Clock{thread_name, vector, sync.Mutex{}}
	return
}

//just add 1 to this thread's counter value.
func (self Vector_Clock) increaseOne() {
	self.vector[self.thread]++
	return
}

//when a message is sent from one thread to another, the receiver call this method.
//to update its vector clock.
func (self Vector_Clock) update(v map[string]int) (e error) {
	//self.mutex.Lock()
	//defer self.mutex.Unlock()
	//two vector_clock should have same vector.
	if len(self.vector) != len(v) {
		e = fmt.Errorf("Update %s's vector clock failed. It has a different length with message sender\n", self.thread)
	}
	for thread, value := range v {
		if _, ok := self.vector[thread]; !ok {
			e = fmt.Errorf("Thread Information doesn't match, cannot find %s in %s's vector clock\n", thread, self.thread)
		}
		// take the maximum one of every value in the vector.
		self.vector[thread] = max(value, self.vector[thread])
	}
	return
}

func (self Vector_Clock) sendMessage(receiver Vector_Clock, ch chan map[string]int) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.increaseOne()
	ch <- self.vector
	fmt.Printf("Thread %s has sent a message to Thread %s\n", self.thread, receiver.thread)
}

func (self Vector_Clock) receiveMessage(sender Vector_Clock, ch chan map[string]int) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.increaseOne()

	vector := <-ch
	err := self.update(vector)
	if err == nil {
		fmt.Printf("Thread %s has received a message from Thread %s\n", self.thread, sender.thread)
	} else {
		fmt.Print(err.Error())
	}

}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
