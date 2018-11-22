package main

import (
	"bytes"
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

func testExample(events [][]Event, correctClocks [][]int) (err error) {
	var wg sync.WaitGroup
	numProcess := len(events)

	processes := make([]int, numProcess)
	for procIdx := range events {
		processes[procIdx] = procIdx
	}

	wg.Add(len(processes))

	vectorClocks := make([]vc.VectorClock, numProcess)
	chans := make([]chan map[int]int, numProcess)

	for idx, _ := range chans {
		//Only works with zero buffering as it needs the waiting to communicate as the method of syncronization between the processes
		chans[idx] = make(chan map[int]int, 0)
	}

	//Create the vectorclocks
	for idx, _ := range vectorClocks {
		vectorClocks[idx], err = vc.NewVectorClock(idx, processes, chans)

		if err != nil {
			return
		}
		startProcess(&vectorClocks[idx], events[idx], &wg)
	}

	wg.Wait()

	var errBuff bytes.Buffer
	var numErr int = 0
	for processIdx := 0; processIdx < numProcess; processIdx++ {
		for clockIdx := 0; clockIdx < numProcess; clockIdx++ {
			correctClock := correctClocks[processIdx][clockIdx]
			computedClock, _ := vectorClocks[processIdx].GetClock(clockIdx)
			if correctClock != computedClock {
				numErr++
				if numErr > 1 {
					errBuff.WriteString("\n")
				}
				errBuff.WriteString(fmt.Sprintf("Error: Process %d clock %d should be %d, got %d", processIdx, clockIdx, correctClock, computedClock))
			}
		}
	}

	if numErr > 1 {
		err = fmt.Errorf(errBuff.String())
	}

	return
}

//Example from the wikipedia page on vector clocks
func testWikiExample() error {
	events := [][]Event{
		[]Event{Event{Receive, 1}, Event{Send, 1}, Event{Receive, 2}, Event{Receive, 2}},
		[]Event{Event{Receive, 2}, Event{Send, 0}, Event{Send, 2}, Event{Receive, 0}, Event{Send, 2}},
		[]Event{Event{Send, 1}, Event{Receive, 1}, Event{Send, 0}, Event{Receive, 1}, Event{Send, 0}},
	}

	correctClocks := [][]int{
		{4, 5, 5},
		{2, 5, 1},
		{2, 5, 5},
	}

	return testExample(events, correctClocks)
}

/*

0   o oo
   / /  \
1 o /    o  o
   /       /
2 o oo    /
   /  \  /
3 o    oo

*/
func testMyExample() error {
	events := [][]Event{
		[]Event{Event{Receive, 1}, Event{Receive, 2}, Event{Send, 1}},
		[]Event{Event{Send, 0}, Event{Receive, 0}, Event{Receive, 3}},
		[]Event{Event{Send, 0}, Event{Receive, 3}, Event{Send, 3}},
		[]Event{Event{Send, 2}, Event{Receive, 2}, Event{Send, 1}},
	}

	correctClocks := [][]int{
		{3, 1, 1, 0},
		{3, 2, 3, 3},
		{0, 0, 3, 1},
		{0, 0, 3, 3},
	}

	return testExample(events, correctClocks)
}

func reportTest(err error) {
	if err == nil {
		fmt.Printf("PASSED\n")
	} else {
		fmt.Printf("FAILED:\n%s\n", err.Error())
	}
}

func main() {
	fmt.Printf("Testing wiki example...\n")
	reportTest(testWikiExample())

	fmt.Printf("Testing my example...\n")
	reportTest(testMyExample())
}
