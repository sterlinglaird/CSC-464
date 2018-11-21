package main

import (
	"flag"
	"fmt"
	"strings"

	byzGen "./byzantineGenerals"
)

func strToGeneral(generalStr string, rank byzGen.Rank) (general byzGen.General, err error) {
	splitGeneral := strings.Split(generalStr, ":")
	if len(splitGeneral) != 2 {
		err = fmt.Errorf("All generals must have both a name and an affinity")
		return
	}

	generalName := strings.Split(generalStr, ":")[0]
	affinityStr := strings.Split(generalStr, ":")[1]

	var aff byzGen.Affinity
	if affinityStr == "L" {
		aff = byzGen.Loyal
	} else if affinityStr == "T" {
		aff = byzGen.Traitor
	} else {
		err = fmt.Errorf("General %s must have an affinity of L or T, got %s", generalName, affinityStr)
		return
	}

	general = byzGen.General{generalName, aff, rank}

	return
}

func validateArgs(recLvl int, generalsStr string, orderStr string) (outRecLvl int, outGenerals []byzGen.General, outOrder byzGen.Order, err error) {
	//Recursion level
	outRecLvl = recLvl
	if outRecLvl <= 0 {
		err = fmt.Errorf("Recursion level must be greater that 0")
		return
	}

	//Generals
	generalsSplitStr := strings.Split(generalsStr, ",")
	if len(generalsSplitStr) < 2 {
		err = fmt.Errorf("There must be more than one general")
		return
	}

	outGenerals = make([]byzGen.General, len(generalsSplitStr))

	outGenerals[0], err = strToGeneral(generalsSplitStr[0], byzGen.Commander)
	if err != nil {
		return
	}

	for generalIdx := 1; generalIdx < len(generalsSplitStr); generalIdx++ {
		outGenerals[generalIdx], err = strToGeneral(generalsSplitStr[generalIdx], byzGen.Lieutenant)
		if err != nil {
			return
		}
	}

	//Order
	if orderStr == "ATTACK" {
		outOrder = byzGen.Attack
	} else if orderStr == "RETREAT" {
		outOrder = byzGen.Retreat
	} else {
		err = fmt.Errorf("Order must be ATTACK or RETREAT, got %s", orderStr)
		return
	}

	return
}

func main() {
	var recLvl int
	var generalsStr string
	var orderStr string

	flag.IntVar(&recLvl, "r", 0, "Level of recursion")
	flag.StringVar(&generalsStr, "g", "", "A list of generals of form G0:L,G1:L,G2:T,... with L=loyal and T=traitor. First general is the commander")
	flag.StringVar(&orderStr, "o", "", "The order that the commander will give to the other generals. Can be ATTACK or RETREAT")

	flag.Parse()

	recLvl, generals, order, err := validateArgs(recLvl, generalsStr, orderStr)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	numLoyal, numTraitor := 0, 0
	for genIdx := range generals {
		if generals[genIdx].Affinity == byzGen.Loyal {
			numLoyal++
		} else {
			numTraitor++
		}
	}

	result := byzGen.ByzantineGenerals(generals, order, 1)

	for idx := range result.OrderTaken {
		traitorStr := ""
		if generals[idx].Affinity == byzGen.Traitor {
			traitorStr = "(traitor) "
		}
		fmt.Printf("%s %stook order %s\n", generals[idx].Name, traitorStr, result.OrderTaken[idx].ToString())
	}
}
