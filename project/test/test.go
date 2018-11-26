package main

import (
	"fmt"
	"math/rand"

	d "../document"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randChar() string {
	return string(letters[rand.Intn(len(letters))])
}

//Random number between x and y inclusive on both ends [x,y]
func randBetween(x int, y int) int {
	//fmt.Printf("%d %d\n", x, y)
	return int(rand.Intn(int(y-x+1)) + int(x))
}

func amountToMove(currIdx int, goalIdx int) int {
	if currIdx > goalIdx {
		return goalIdx - currIdx
	} else {
		return goalIdx - currIdx
	}
}

func addAt(doc []string, s string, at int) []string {
	if len(doc) >= at+1 {
		doc = append(doc, "")
		copy(doc[at+1:], doc[at:])
		doc[at] = s
		return doc
	} else {
		doc = append(doc, s)
		return doc
	}
}

//Fuzzes with 100 random insertions
func fuzzAdd() (err error) {
	site := 1
	doc := d.NewDocument(site)

	firstChar := randChar()

	groundTruth := []string{firstChar}
	pos, err := doc.InsertRight(d.StartPos, firstChar)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	lowIdx := 0
	highIdx := 1
	currIdx := 1
	addIdx := currIdx
	for idx := 0; idx < 100; idx++ {
		char := randChar()

		addIdx = randBetween(lowIdx, highIdx)
		for addIdx == currIdx {
			addIdx = randBetween(lowIdx, highIdx)
		}

		//Ground truth
		groundTruth = addAt(groundTruth, char, addIdx)

		//Actual document
		pos, err = doc.Move(pos, addIdx-currIdx)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
		pos, err = doc.InsertRight(pos, char)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}

		currIdx = addIdx + 1
		highIdx += 1

		testResult, _ := doc.GetContent()

		for idx := range groundTruth {
			if string(testResult[idx]) != groundTruth[idx] {
				err = fmt.Errorf("Ground truth: %s\nResult: %s.\nFull doc: %s", groundTruth, testResult, doc.ToString())
				return
			}
		}
	}

	testResult, err := doc.GetContent()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	for idx := range groundTruth {
		if string(testResult[idx]) != groundTruth[idx] {
			err = fmt.Errorf("Ground truth: %s\nResult: %s.\nFull doc: %s", groundTruth, testResult, doc.ToString())
		}
	}

	return
}

func main() {
	fmt.Printf("Fuzzing insert...\n")
	wasErr := false

	//Fuzz test 100 times with a different random seed
	for idx := 0; idx < 100; idx++ {
		rand.Seed(int64(idx))
		err := fuzzAdd()
		if err != nil {
			fmt.Printf("FAILED: %s\n", err.Error())
			fmt.Printf("Seed: %d\n", idx)
			wasErr = true
		}
	}

	if !wasErr {
		fmt.Printf("PASSED")
	}
}
