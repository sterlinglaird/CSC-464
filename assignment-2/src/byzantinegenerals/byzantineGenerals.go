package byzantineGenerals

type Order int

const (
	Attack Order = iota
	Retreat
	Undetermined
)

func (o Order) ToString() string {
	switch o {
	case Attack:
		return "ATTACK"
	case Retreat:
		return "RETREAT"
	case Undetermined:
		return "undetermined order"
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

type General struct {
	Name     string
	Affinity Affinity
	Rank     Rank
}

type Result struct {
	OrderTaken []Order //Orders taken by each general, ordered same way as input general list. [0] is the final consensus
}

type node struct {
	Id        int
	Received  Order
	Reply     Order
	Children  []*node
	Parent    *node
	PathTaken map[int]bool //whether a node id has been seen in the path of this order
	Depth     int          //How deep the node is in the message tree. Commander is at 0
}

func majority(orders map[Order]int) (order Order) {
	order = Retreat
	if orders[Attack] > orders[Retreat] {
		order = Attack
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

//Sends a message to a node by constructing a new node with a higher depth
func (this *node) sendMessage(destId int, traitors map[int]bool) *node {
	if this.PathTaken[destId] {
		return nil
	}

	order := this.Received

	//If you are a traitor, then swap the order for certain generals. Arbitrary choice to make it 50/50
	if traitors[this.Id] {
		if destId%2 == 0 {
			order = swapOrder(order)
		}
	}

	//Copy the paths taken into a new map and mark it as seen by us
	pathTaken := make(map[int]bool)
	for k, v := range this.PathTaken {
		pathTaken[k] = v
	}
	pathTaken[destId] = true

	return &node{destId, order, Undetermined, nil, this, pathTaken, this.Depth + 1}
}

func (this *node) decide() Order {
	//If there is no more children then we just take our own recieved message as the reply
	if len(this.Children) == 0 {
		this.Reply = this.Received
		return this.Reply
	}

	orders := map[Order]int{Attack: 0, Retreat: 0}
	orders[this.Received]++

	//Decide on the order to take at this level by looking at the majority of our siblings children.
	//Easy to understand graphically, this is basically getting the messages sent to us one level down
	//@TODO I think this step can be made faster by caching the results since the will be computed multiple times, really doesnt matter for this though..
	for _, sibling := range this.Parent.Children {
		if sibling.Id != this.Id {
			for _, child := range sibling.Children {
				if child.Id == this.Id {
					orders[child.decide()]++
				}
			}
		}

	}

	this.Reply = majority(orders)
	return this.Reply
}

//Constructs the tree with all the simulated messages passed between generals
func constructMessageTree(numGenerals, numTraitors int, order Order, traitors map[int]bool) *node {
	//Root node for the tree. Represents the commander
	commanderNode := &node{0, order, Undetermined, nil, nil, map[int]bool{0: true}, 0}

	nodeQueue := []*node{commanderNode}

	for len(nodeQueue) > 0 {
		//Take first element of queue
		currNode := nodeQueue[0]
		nodeQueue = nodeQueue[1:]

		//Only go numTraitor times deep since as lamport described that is all that is needed inn OM(m)
		if currNode.Depth < numTraitors {
			for destID := 1; destID < numGenerals; destID++ {
				child := currNode.sendMessage(destID, traitors)
				if child != nil {
					currNode.Children = append(currNode.Children, child)
					child.Parent = currNode

					//Keep appending to the queue so the whole tree is generated
					nodeQueue = append(nodeQueue, child)
				}
			}
		}
	}
	return commanderNode
}

func ByzantineGenerals(generals []General, order Order) (result Result) {
	traitors := make(map[int]bool)
	var numTraitors int = 0
	for generalIdx := range generals {
		if generals[generalIdx].Affinity == Traitor {
			numTraitors++
			traitors[generalIdx] = true
		}
	}

	tree := constructMessageTree(len(generals), numTraitors, order, traitors)

	result.OrderTaken = make([]Order, len(generals))

	//Calculate the orders for the lieutenants
	orders := map[Order]int{Attack: 0, Retreat: 0}
	var genIdx int = 1
	for _, general := range tree.Children {
		order := general.decide()
		result.OrderTaken[genIdx] = order
		orders[order]++
		genIdx++
	}

	//Calculate the majority for the commander (the consensus)
	result.OrderTaken[0] = majority(orders)

	return
}
