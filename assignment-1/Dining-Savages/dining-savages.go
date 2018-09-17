package main

import (
	"fmt"
	"math/rand"
	"time"
)

const NUM_SAVAGES = 10
const MAX_WAIT_TIME = 50
const MAX_SERVINGS = 5

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

type Pot struct {
	filledAlert chan bool
	servings    chan int
	waiting     chan bool
	alertNum    int
}

func newPot(numServings int) *Pot {
	pot := new(Pot)

	pot.servings = make(chan int, numServings)
	pot.filledAlert = make(chan bool)
	pot.waiting = make(chan bool)

	return pot
}

func (pot *Pot) getServing() int {
	var serving int

	//If there is no available serving, then wait until the cook fills the pot
	select {
	case serving = <-pot.servings:
		return serving
	default:
		pot.waiting <- true
		<-pot.filledAlert
		serving = pot.getServing()
	}

	return serving
}

func (pot *Pot) fillPot() {
	//Make sure all savages are waiting
	for i := 0; i < NUM_SAVAGES; i++ {
		<-pot.waiting
	}

	//Number all of the servings
	for i := 0; i < MAX_SERVINGS; i++ {
		pot.servings <- i
	}

	//Alert the tribe that the pot is full
	for i := 0; i < NUM_SAVAGES; i++ {
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
	}
}

func savage(id int) {
	for {
		fmt.Printf("savage %d looking for food\n", id)
		serving := pot.getServing()

		//Alert the cook that the pot is empty if we got the last serving
		if serving == MAX_SERVINGS-1 {
			wakeCook <- true
		}

		fmt.Printf("savage %d got serving #%d\n", id, serving)
		eat()
	}
}

func main() {
	pot = newPot(MAX_SERVINGS)

	go cook()

	//Get cook to fill pot initially
	wakeCook <- true

	for i := 0; i < NUM_SAVAGES; i++ {
		go savage(i)
	}

	for {
	}
}
