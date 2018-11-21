package main

import (
	"fmt"

	byzGen "./byzantineGenerals"
)

type General byzGen.General
type Order byzGen.Order
type Affinity byzGen.Affinity
type Rank byzGen.Rank

func main() {
	// var recLvl int
	// var generalsStr string
	// var orderStr string

	// flag.IntVar(&recLvl, "r", 0, "Level of recursion")
	// flag.StringVar(&generalsStr, "g", "", "A list of generals of form G0:L,G1:L,G2:T,... with L=loyal and T=traitor. First general is the commander")
	// flag.StringVar(&orderStr, "o", "", "The order that the commander will give to the other generals. Can be ATTACK or RETREAT")

	// flag.Parse()

	// recLvl, generals, order, err := validateArgs(recLvl, generalsStr, orderStr)
	// if err != nil {
	// 	fmt.Printf("Error: %s\n", err.Error())
	// 	return
	// }

	// var wg sync.WaitGroup
	// wg.Add(len(generals))

	// //First general is the commander
	// go commander(generals[0], order, &wg)
	// for generalIdx := 1; generalIdx < len(generals); generalIdx++ {
	// 	go lieutenant(generals[generalIdx], &wg)
	// }

	// wg.Wait()

	generals := []byzGen.General{
		byzGen.General{"C", byzGen.Loyal, byzGen.Commander},
		byzGen.General{"L1", byzGen.Loyal, byzGen.Lieutenant},
		byzGen.General{"L2", byzGen.Loyal, byzGen.Lieutenant},
		byzGen.General{"L3", byzGen.Loyal, byzGen.Lieutenant},
		byzGen.General{"L4", byzGen.Loyal, byzGen.Lieutenant},
		byzGen.General{"T5", byzGen.Traitor, byzGen.Lieutenant},
		byzGen.General{"T6", byzGen.Traitor, byzGen.Lieutenant},
	}

	numLoyal, numTraitor := 0, 0
	for genIdx := range generals {
		if generals[genIdx].Affinity == byzGen.Loyal {
			numLoyal++
		} else {
			numTraitor++
		}
	}

	result := byzGen.ByzantineGenerals(generals, byzGen.Attack, 1)

	for idx := range result.OrderTaken {
		traitorStr := ""
		if generals[idx].Affinity == byzGen.Traitor {
			traitorStr = "(traitor) "
		}
		fmt.Printf("%s %stook order %s\n", generals[idx].Name, traitorStr, result.OrderTaken[idx].ToString())
	}
}
