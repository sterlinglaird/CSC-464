package main

import (
	"fmt"
	"math/rand"
	"time"
)

const NUM_PRODUCERS = 2
const NUM_CONSUMERS = 10
const BUFFER_SIZE = 20
const MAX_WAIT_TIME = 50

var buf = make(chan int, BUFFER_SIZE)

func waitForEvent() int {
	var randomNum = rand.Int() % 10000
	time.Sleep(time.Duration(randomNum%MAX_WAIT_TIME) * time.Millisecond)
	return randomNum
}

func consumeEvent(event int) {
	var randomNum = rand.Int()
	time.Sleep(time.Duration(randomNum%MAX_WAIT_TIME) * time.Millisecond)
}

func producer(id int) {
	for {
		var event = waitForEvent()
		fmt.Printf("produced: %d by %d\n", event, id)
		buf <- event
	}
}

func consumer(id int) {
	for {
		event := <-buf
		consumeEvent(event)
		fmt.Printf("consumed: %d by %d\n", event, id)
	}
}

func main() {
	for i := 0; i < NUM_PRODUCERS; i++ {
		go producer(i)
	}

	for i := 0; i < NUM_CONSUMERS; i++ {
		go consumer(i)
	}

	for {
	}
}
