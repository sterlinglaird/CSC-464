package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	byzGen "./byzantineGenerals"
)

func runTest(numTraitor int, numGeneral int, order byzGen.Order) (err error) {
	// expectConsensus := numGeneral > (3 * numTraitor)

	// if !expectConsensus {
	// 	err = fmt.Errorf("Not expected to reach consensus so test is not meaningful")
	// 	return
	// }

	generals := make([]byzGen.General, numGeneral)

	//Generate the indices for the traitors.
	traitorIndices := rand.Perm(numGeneral)[0:numTraitor]

	//Generate the generals
	for genIdx := 0; genIdx < numGeneral; genIdx++ {
		affinity := byzGen.Loyal
		prefix := "L"
		rank := byzGen.Lieutenant
		for traitorIdx := 0; traitorIdx < numTraitor; traitorIdx++ {
			if traitorIndices[traitorIdx] == genIdx {
				affinity = byzGen.Traitor
				prefix = "T"
			}
		}

		if genIdx == 0 {
			rank = byzGen.Commander
			prefix = "C"
		}

		generals[genIdx] = byzGen.General{fmt.Sprintf("%s%d", prefix, genIdx), affinity, rank}
	}

	result := byzGen.ByzantineGenerals(generals, order)

	numLoyalLieutenants := numGeneral - numTraitor - 1
	if generals[0].Affinity == byzGen.Traitor {
		numLoyalLieutenants += 1 //Add back the general since it counts towards the traitors
	}

	var errBuff bytes.Buffer
	var numErr int = 0

	//Check if all loyal lieutenants follow the same order.
	var loyalAttack int = 0
	for idx := 1; idx < len(result.OrderTaken); idx++ {
		if generals[idx].Affinity == byzGen.Loyal {
			if result.OrderTaken[idx] == byzGen.Attack {
				loyalAttack++
			}
		}
	}
	if loyalAttack != 0 && loyalAttack != numLoyalLieutenants {
		errBuff.WriteString("Not all loyal lieutenants followed same order.")
		numErr++
	}

	//Check if all loyal lieutenants order from loyal commander.
	if generals[0].Affinity == byzGen.Loyal {
		var loyalAttack int = 0
		for idx := 1; idx < len(result.OrderTaken); idx++ {
			if generals[idx].Affinity == byzGen.Loyal {
				if result.OrderTaken[idx] == byzGen.Attack {
					loyalAttack++
				}
			}
		}
		if loyalAttack != 0 && loyalAttack != numLoyalLieutenants {
			numErr++
			if numErr > 1 {
				errBuff.WriteString("\n\t")
			}
			errBuff.WriteString("Not all loyal lieutenants followed order from loyal commander.")
		}
	}

	if numErr > 1 {
		if generals[0].Affinity == byzGen.Loyal {
			errBuff.WriteString("\n\tCommander was a traitor")
		}
		errBuff.WriteString(fmt.Sprintf("\n\tConsensus was %s", result.OrderTaken[0].ToString()))
		for idx := 1; idx < len(result.OrderTaken); idx++ {
			traitorStr := ""
			if generals[idx].Affinity == byzGen.Traitor {
				traitorStr = "(traitor) "
			}
			errBuff.WriteString(fmt.Sprintf("\n\t%s %stook order %s", generals[idx].Name, traitorStr, result.OrderTaken[idx].ToString()))
		}

		err = fmt.Errorf(errBuff.String())
	}

	return
}

func main() {
	rand.Seed(time.Now().Unix())
	testSizes := []int{4, 7, 10, 13, 16}

	var numFailed int = 0
	var numTotal int = 0
	for sizeIdx := range testSizes {
		size := testSizes[sizeIdx]

		//Only do the tests which we should reach consensus, otherwise we cannot tell a success from a fail
		//Test both attack and retreat orders
		for numTraitor := 0; numTraitor*3 < size; numTraitor++ {
			numTotal++

			fmt.Printf("Running test with Num generals: %d, Num traitors: %d, Order: %s\n", size, numTraitor, byzGen.Attack.ToString())
			err := runTest(numTraitor, size, byzGen.Attack)
			if err == nil {
				fmt.Printf("\tPASSED:\n")
			} else {
				fmt.Printf("\tFAILED:\n\t%s\n", err.Error())
				numFailed++
			}

			numTotal++

			fmt.Printf("Running test with Num generals: %d, Num traitors: %d, Order: %s\n", size, numTraitor, byzGen.Retreat.ToString())
			err = runTest(numTraitor, size, byzGen.Retreat)
			if err == nil {
				fmt.Printf("\tPASSED:\n")
			} else {
				fmt.Printf("\tFAILED:\n\t%s\n", err.Error())
				numFailed++
			}
		}
	}

	fmt.Printf("Number Failed: %d/%d\n", numFailed, numTotal)

}
