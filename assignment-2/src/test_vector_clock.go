package main

import (
	"fmt"
	"sync"

	vc "./vectorclock"
)

type Action int

//Self is a self incrementing action. Only increments itself. Other events send/recieve after doing the increment.
const (
	Self Action = iota
	Send
	Receive
)

type Event struct {
	action   Action
	targetId int
}

func startProcess(vc *vc.VectorClock, eventList []Event, wg *sync.WaitGroup) (err error) {
	go func() {
		for _, event := range eventList {
			switch event.action {
			case Self:
				vc.Inc()
			case Send:
				vc.Inc()
				vc.Send(event.targetId)
			case Receive:
				vc.Inc()
				vc.Receive(event.targetId)
			default:
				err = fmt.Errorf("Error: Unimplemented action %d", event.action)
			}
		}
		wg.Done()
	}()
	return
}

func testExample() (err error) {
	var wg sync.WaitGroup

	processes := []int{0, 1, 2}
	wg.Add(len(processes))
	vectorClocks := make([]vc.VectorClock, len(processes))
	chans := make([]chan map[int]int, len(processes))

	for idx, _ := range chans {
		//Buffering would not cause any drawback, I choose not to do it so everything happens sequencially
		chans[idx] = make(chan map[int]int)
	}

	//Create the vectorclocks
	for idx, _ := range vectorClocks {
		vectorClocks[idx], err = vc.NewVectorClock(idx, processes, chans)

		if err != nil {
			return
		}
	}

	startProcess(&vectorClocks[0], []Event{Event{Receive, 1}, Event{Send, 1}, Event{Receive, 2}, Event{Receive, 2}}, &wg)
	startProcess(&vectorClocks[1], []Event{Event{Receive, 2}, Event{Send, 0}, Event{Send, 2}, Event{Receive, 0}, Event{Send, 2}}, &wg)
	startProcess(&vectorClocks[2], []Event{Event{Send, 1}, Event{Receive, 1}, Event{Send, 0}, Event{Receive, 1}, Event{Send, 0}}, &wg)

	wg.Wait()

	//Correct values from the example on wikipedia page for vector clocks
	correctClocks := [][]int{
		{4, 5, 5},
		{2, 5, 1},
		{2, 5, 5},
	}

	for processIdx := 0; processIdx < 3; processIdx++ {
		for clockIdx := 0; clockIdx < 3; clockIdx++ {
			correctClock := correctClocks[processIdx][clockIdx]
			computedClock, _ := vectorClocks[processIdx].GetClock(clockIdx)
			if correctClock != computedClock {
				err = fmt.Errorf("Error: Process %d clock %d should be %d, got %d", processIdx, clockIdx, correctClock, computedClock)
			}
		}
	}

	return

}

func main() {
	fmt.Printf("Testing example...\n")
	err := testExample()
	if err == nil {
		fmt.Printf("PASSED\n")
	} else {
		fmt.Printf("FAILED: %s\n", err.Error())
	}

}
