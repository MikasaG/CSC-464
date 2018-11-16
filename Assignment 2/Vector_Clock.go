// Vector_Clock.go   Chang Gong  V00898803
package main

import (
	"fmt"
	"sync"
)

type Vector_Clock struct {
	thread string
	vector map[string]int
	mutex  sync.Mutex
}

func newVectorClock(index int, thread_name string, threads []string) (vc Vector_Clock, e error) {
	if index < 0 || index >= len(threads) {
		e = fmt.Errorf("The index value of the new vector clock is invalid")
		return
	}
	vector := make(map[string]int)
	for _, thread = range threads {
		vector[thread] = 0
	}
	vc = Vector_Clock{thread_name, vector, sync.Mutex{}}
	return
}

func (self Vector_Clock) increaseOne() {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.vector[self.thread]++
	return
}

func (self Vector_Clock) update(sender Vector_Clock) (e error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	if len(self.vector) != len(sender.vector) {
		e = fmt.Errorf("Update %s's vector clock failed. It has a different length with message sender %s", self.thread, sender.thread)
	}
	for thread, value = range sender.vector {
		if _, ok = self.vector[thread]; !ok {
			e = fmt.Errorf("Thread Information doesn't match, cannot find %s in recevier's vector clock", thread)
		}
		self.vector[thread] = max(value, self.vector[thread])
	}
	return
}
