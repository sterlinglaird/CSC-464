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
		//fmt.Printf("%d: %s\n", vc.Id, vc.GetClockString())
	}()
	return
}

//Uses the example from the wikipedia page for vector clocks
func testExample() (err error) {
	var wg sync.WaitGroup

	processes := []int{0, 1, 2}
	numProcess := len(processes)
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
	}

	startProcess(&vectorClocks[0], []Event{Event{Receive, 1}, Event{Send, 1}, Event{Receive, 2}, Event{Receive, 2}}, &wg)
	startProcess(&vectorClocks[1], []Event{Event{Receive, 2}, Event{Send, 0}, Event{Send, 2}, Event{Receive, 0}, Event{Send, 2}}, &wg)
	startProcess(&vectorClocks[2], []Event{Event{Send, 1}, Event{Receive, 1}, Event{Send, 0}, Event{Receive, 1}, Event{Send, 0}}, &wg)

	wg.Wait()

	correctClocks := [][]int{
		{4, 5, 5},
		{2, 5, 1},
		{2, 5, 5},
	}

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

func main() {
	fmt.Printf("Testing example...\n")
	err := testExample()
	if err == nil {
		fmt.Printf("PASSED\n")
	} else {
		fmt.Printf("FAILED:\n%s\n", err.Error())
	}
}
