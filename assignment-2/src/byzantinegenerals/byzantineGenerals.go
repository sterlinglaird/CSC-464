package byzantineGenerals

import (
	"reflect"
	"sync"
)

type Order int

const (
	Attack Order = iota
	Retreat
)

func (o Order) ToString() string {
	switch o {
	case Attack:
		return "ATTACK"
	case Retreat:
		return "RETREAT"
	default:
		return "unimplemented order type"
	}
}

type Affinity int

const (
	Loyal Affinity = iota
	Traitor
)

type Rank int

const (
	Commander Rank = iota
	Lieutenant
)

type Message struct {
	order Order
	depth int
}

type General struct {
	Name     string
	Affinity Affinity
	Rank     Rank
}

type Result struct {
	OrderTaken []Order //Orders taken by each general, ordered same way as input general list
}

func majority(orders map[string]Order, ignoredGenerals []string) (order Order) {
	var diff int = 0 //#attack - #retreat
	for genName := range orders {
		for igName := range ignoredGenerals {
			if ignoredGenerals[igName] == genName {
				continue
			}
		}
		if orders[genName] == Attack {
			diff++
		} else {
			diff--
		}
	}

	if diff > 0 {
		order = Attack
	} else {
		order = Retreat
	}

	return
}

func swapOrder(ord Order) Order {
	switch ord {
	case Attack:
		return Retreat
	case Retreat:
		return Attack
	default:
		return -1
	}
}

func detOrderToSend(recievedOrder Order, aff Affinity) Order {
	if aff == Loyal {
		return recievedOrder
	} else {
		return swapOrder(recievedOrder)
	}
}

func runCommander(thisGen General, order Order, messengersTo map[string]chan Message, orderTaken *Order, wg *sync.WaitGroup) {
	defer wg.Done()

	*orderTaken = order

	//Send order to all lietenants (not itself)
	for genName, messenger := range messengersTo {
		if genName != thisGen.Name {
			//fmt.Printf("%s sent %s to %s\n", thisGen.Name, order.ToString(), genName)
			messenger <- Message{order, 0}
		}
	}
}

func recieveMessages(thisGen General, orders map[string]Order, messengers map[string]chan Message, wg *sync.WaitGroup) {
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
		//fmt.Printf("%s recieved %s from %s\n", thisGen.Name, value.Interface().(Order).ToString(), idxToMessenger[chosen])
		orders[idxToMessenger[chosen]] = value.Interface().(Message).order
	}
}

//From this general to _, to this genereal from _
func runLieutenant(thisGen General, commanderName string, orderTaken *Order, toMessengers map[string]chan Message, fromMessengers map[string]chan Message, wg *sync.WaitGroup) {
	defer wg.Done()

	orders := make(map[string]Order)

	//Recieve the order from the commander
	orders[commanderName] = (<-toMessengers[commanderName]).order
	//fmt.Printf("%s recieved %s from %s\n", thisGen.Name, orders[commanderName].ToString(), commanderName)

	orderToSend := orders[commanderName]

	//Recieve all the messages from the other lietenants
	var recieveWg sync.WaitGroup
	recieveWg.Add(1)
	go recieveMessages(thisGen, orders, toMessengers, &recieveWg)

	//Send order to all other lietenants
	for genName, messenger := range fromMessengers {
		if genName != thisGen.Name && genName != commanderName {
			//fmt.Printf("%s sent %s to %s\n", thisGen.Name, orderToSend(orders[commanderName], thisGen.Affinity).ToString(), genName)
			messenger <- Message{detOrderToSend(orderToSend, thisGen.Affinity), 0}

			//Keep changing relayed order if a traitor
			// if thisGen.Affinity == traitor {
			// 	orderToSend = swapOrder(orderToSend)
			// }
		}
	}

	recieveWg.Wait()

	// var numAtt, numRet int = 0, 0
	// for k, v := range messages {
	// 	if k == thisGen.Name || k == commanderName {
	// 		continue
	// 	}
	// 	if v.order == Attack {
	// 		numAtt++
	// 	} else {
	// 		numRet++
	// 	}
	// }

	//fmt.Printf("%s recieved %d attack and %d retreat\n", thisGen.Name, numAtt, numRet)
	//We take the majority as the order we will take (ignoring our results and the commanders results)
	*orderTaken = majority(orders, []string{commanderName})
}

func ByzantineGenerals(generals []General, order Order, recLvl int) (result Result) {
	result.OrderTaken = make([]Order, len(generals)) //Each general will add their enty when they compute it

	var wg sync.WaitGroup
	wg.Add(len(generals))

	var commanderName string
	var commanderIdx int = 0
	for genIdx := range generals {
		if generals[genIdx].Rank == Commander {
			commanderName = generals[genIdx].Name
			commanderIdx = genIdx
		}
	}

	//From -> To -> channel
	messengers := make(map[string]map[string]chan Message)
	for fromGen := range generals {
		messengers[generals[fromGen].Name] = make(map[string]chan Message)
		for toGen := range generals {
			messengers[generals[fromGen].Name][generals[toGen].Name] = make(chan Message)
		}
	}

	//First general is the commander, rest are lietenants
	go runCommander(generals[commanderIdx], order, messengers[commanderName], &result.OrderTaken[0], &wg)
	for generalIdx := 1; generalIdx < len(generals); generalIdx++ {
		toMessengers := make(map[string]chan Message)
		fromMessengers := make(map[string]chan Message)

		for fromKey, _ := range messengers {
			for toKey, ch := range messengers[fromKey] {
				if toKey == generals[generalIdx].Name {
					toMessengers[fromKey] = ch
				}
			}
		}

		for toKey, ch := range messengers[generals[generalIdx].Name] {
			fromMessengers[toKey] = ch
		}

		go runLieutenant(generals[generalIdx], commanderName, &result.OrderTaken[generalIdx], toMessengers, fromMessengers, &wg)
	}

	wg.Wait()

	return
}
