package main

import (
	"fmt"
	"math/rand"
	"time"
)

const NUM_SAVAGES = 5
const MAX_WAIT_TIME = 50
const MAX_SERVINGS = 100

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

type Pot struct {
	filledAlert chan bool
	servings    chan int
	alertNum    int
}

func newPot(numServings int) *Pot {
	pot := new(Pot)

	pot.alertNum = min(MAX_SERVINGS, NUM_SAVAGES)
	pot.servings = make(chan int, numServings)
	pot.filledAlert = make(chan bool, pot.alertNum)

	return pot
}

func (pot *Pot) getServing() int {
	var serving int

	//If there is no available serving, then wait until the cook fills the pot
	select {
	case serving = <-pot.servings:
		return serving
	default:
		<-pot.filledAlert
		serving = pot.getServing()
	}

	return serving
}








Try maintaining a "waiting" channel and alert for each waiting











func (pot *Pot) fillPot() {
	//Number all of the servings
	for i := 0; i < MAX_SERVINGS; i++ {
		fmt.Printf("filled %d\n", i)
		pot.servings <- i
	}

	//Alert savages that pot is full. Only wakes the minimal amount
	for i := 0; i < pot.alertNum; i++ {
		fmt.Printf("alert %d\n", i)
		pot.filledAlert <- true
	}
}

var pot *Pot
var wakeCook = make(chan bool)

func eat() {
	var randomNum = rand.Int()
	time.Sleep(time.Duration(randomNum%MAX_WAIT_TIME) * time.Millisecond)
}

func cook() {
	for {
		<-wakeCook
		fmt.Printf("cook refilling pot\n")
		pot.fillPot()
		fmt.Printf("cook refilled pot\n")
	}
}

func savage(id int) {
	for {
		fmt.Printf("savage %d looking for food\n", id)
		serving := pot.getServing()

		//Alert the cook that the pot is empty if we got the last serving
		if serving == MAX_SERVINGS-1 {
			fmt.Printf("savage %d alerting cook\n", id)
			wakeCook <- true
		}

		fmt.Printf("savage %d got food\n", id)
		eat()
	}
}

func main() {
	pot = newPot(MAX_SERVINGS)

	for i := 0; i < MAX_SERVINGS; i++ {
		fmt.Printf("filled %d\n", i)
		pot.servings <- i
	}

	go cook()

	//Get cook to make fill pot initially
	//wakeCook <- true

	for i := 0; i < NUM_SAVAGES; i++ {
		go savage(i)
	}

	for {
	}
}
