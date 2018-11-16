package main

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Order int
type Affinity int
type Rank int

const (
	attack Order = iota
	retreat
)

func (o Order) ToString() string {
	switch o {
	case attack:
		return "ATTACK"
	case retreat:
		return "RETREAT"
	default:
		return "unimplemented order type"
	}
}

const (
	loyal Affinity = iota
	traitor
)

const (
	commander Rank = iota
	lieutenant
)

type General struct {
	name     string
	affinity Affinity
	rank     Rank
}

type Result struct {
	orderTaken []Order //Orders taken by each general, ordered same way as input general list
}

func validateArgs(recLvl int, generalsStr string, orderStr string) (outRecLvl int, outGenerals []General, outOrder Order, err error) {
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

	outGenerals = make([]General, len(generalsSplitStr))
	for generalIdx := range generalsSplitStr {
		splitGeneral := strings.Split(generalsSplitStr[generalIdx], ":")
		if len(splitGeneral) != 2 {
			err = fmt.Errorf("All generals must have both a name and an affinity")
			return
		}

		generalName := strings.Split(generalsSplitStr[generalIdx], ":")[0]
		affinityStr := strings.Split(generalsSplitStr[generalIdx], ":")[1]

		var aff Affinity
		if affinityStr == "L" {
			aff = loyal
		} else if affinityStr == "T" {
			aff = traitor
		} else {
			err = fmt.Errorf("General %s must have an affinity of L or T, got %s", generalName, affinityStr)
			return
		}

		//@TODO needs to make a commander
		outGenerals[generalIdx] = General{generalName, aff, lieutenant}
	}

	//Order
	if orderStr == "ATTACK" {
		outOrder = attack
	} else if orderStr == "RETREAT" {
		outOrder = retreat
	} else {
		err = fmt.Errorf("Order must be ATTACK or RETREAT, got %s", orderStr)
		return
	}

	return
}

func majority(orders map[string]Order, ignoredGenerals []string) (order Order) {
	var diff int = 0 //#attack - #retreat
	for genName := range orders {
		for igName := range ignoredGenerals {
			if ignoredGenerals[igName] == genName {
				continue
			}
		}
		if orders[genName] == attack {
			diff++
		} else {
			diff--
		}
	}

	if diff > 0 {
		order = attack
	} else {
		order = retreat
	}
	return
}

func swapOrder(ord Order) Order {
	switch ord {
	case attack:
		return retreat
	case retreat:
		return attack
	default:
		return -1
	}
}

func orderToSend(recievedOrder Order, aff Affinity) Order {
	if aff == loyal {
		return recievedOrder
	} else {
		return swapOrder(recievedOrder)
	}
}

func recieveMessages(thisGen General, orders map[string]Order, messengers map[string]chan Order, wg *sync.WaitGroup) {
	defer wg.Done()

	cases := make([]reflect.SelectCase, len(messengers))

	idxToMessenger := make(map[int]string) //So we can map the result back to a source

	var midx int = 0
	for messengerName, ch := range messengers {
		cases[midx] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		idxToMessenger[midx] = messengerName
		midx++
	}

	//@TODO this should really get ALL messages and then notify when commanders has been gotten
	//Collect N-2 messages (from all generals except commander and ourself)
	for idx := 0; idx < len(messengers)-2; idx++ {
		chosen, value, _ := reflect.Select(cases)
		fmt.Printf("%s recieved %s from %s\n", thisGen.name, value.Interface().(Order).ToString(), idxToMessenger[chosen])
		orders[idxToMessenger[chosen]] = value.Interface().(Order)
	}
}

func runCommander(thisGen General, order Order, messengers map[string]chan Order, orderTaken *Order, wg *sync.WaitGroup) {
	defer wg.Done()

	*orderTaken = order

	//Send order to all lietenants (not itself)
	for genName, messenger := range messengers {
		if genName != thisGen.name {
			//fmt.Printf("%s sent %s to %s\n", thisGen.name, order.ToString(), genName)
			messenger <- order
		}
	}
}

func runLieutenant(thisGen General, commanderName string, orderTaken *Order, messengers map[string]chan Order, wg *sync.WaitGroup) {
	defer wg.Done()

	orders := make(map[string]Order)

	//Recieve the order from the commander
	orders[commanderName] = <-messengers[thisGen.name]
	//fmt.Printf("%s recieved %s from %s\n", thisGen.name, orders[commanderName].ToString(), commanderName)

	//Recieve all the messages from the other lietenants
	var recieveWg sync.WaitGroup
	recieveWg.Add(1)
	go recieveMessages(thisGen, orders, messengers, &recieveWg)

	//Send order to all other lietenants
	for genName, messenger := range messengers {
		if genName != thisGen.name && genName != commanderName {
			//fmt.Printf("%s sent %s to %s\n", thisGen.name, orderToSend(orders[commanderName], thisGen.affinity).ToString(), genName)
			messenger <- orderToSend(orders[commanderName], thisGen.affinity)
		}
	}

	recieveWg.Wait()

	var numAtt, numRet int = 0, 0
	for k, v := range orders {
		if k == thisGen.name || k == commanderName {
			continue
		}
		if v == attack {
			numAtt++
		} else {
			numRet++
		}
	}

	fmt.Printf("%s recieved %d attack and %d retreat\n", thisGen.name, numAtt, numRet)
	//We take the majority as the order we will take (ignoring our results and the commanders results)
	*orderTaken = majority(orders, []string{commanderName, thisGen.name})
}

func byzantineGenerals(generals []General, order Order, recLvl int) (result Result) {
	result.orderTaken = make([]Order, len(generals)) //Each general will add their enty when they compute it

	var wg sync.WaitGroup
	wg.Add(len(generals))

	messengers := make(map[string]chan Order)
	for generalIdx := range generals {
		messengers[generals[generalIdx].name] = make(chan Order)
	}

	//First general is the commander, rest are lietenants
	go runCommander(generals[0], order, messengers, &result.orderTaken[0], &wg)
	for generalIdx := 1; generalIdx < len(generals); generalIdx++ {
		go runLieutenant(generals[generalIdx], generals[0].name, &result.orderTaken[generalIdx], messengers, &wg)
	}

	wg.Wait()

	return
}

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

	generals := []General{
		General{"C", loyal, commander},
		General{"L1", loyal, lieutenant},
		General{"L2", loyal, lieutenant},
		//General{"L3", loyal, lieutenant},
		General{"L4", traitor, lieutenant},
	}

	result := byzantineGenerals(generals, attack, 1)

	fmt.Println()
	for idx := range result.orderTaken {
		fmt.Printf("%s took order %s\n", generals[idx].name, result.orderTaken[idx].ToString())
	}

}
