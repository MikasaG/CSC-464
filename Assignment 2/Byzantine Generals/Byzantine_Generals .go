package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

const (
	ATTACK  = "ATTACK"
	RETREAT = "RETREAT"
)

var G = flag.Int("G", 10, "The number of generals in total, the first one (General 0) is a commander")

var M = flag.Int("M", 3, "The number of traitors, this number should less than 1/3 of G. Otherwise generals cannot come to consensus")

var O = flag.String("O", "ATTACK", "Order,Please input only 'ATTACK' or 'RETREAT'")

//Every Node stands for a order sent from one general to another,
//  Components:
//	General(int): The id of the order receiver.
//	Receive(string): The order itself.
//  Responce(string): The receiver received the order, then he becomes a commander and send it to other generals
//					  who haven't seen this order before, and make a decision based on these generals' responces.
//					  then, he store his decision in this field.
//  children([]*Node): receiver becomes a commander and send this order to childeren.
//	path(map[int]int): Generals id who have seen this order before (i.e, the receiver will not send order to them again)
//	depth(int): the depth of this node in the whole tree. Root(General 0, first commander) is 0
//	parent(*Node): Message sender.
type Node struct {
	General  int
	Receive  string
	Responce string
	children []*Node
	path     map[int]int
	depth    int
	parent   *Node
}

//This method gives a random traitors-map which indicates who are the traitors.
//  Parameters:
//	GeneralsNum(int): The number of generals
//	TraitorsNum(int): The total number of traitors
func getTraitors(GeneralsNum, TraitorsNum int) map[int]int {
	TraitorsMap := make(map[int]int)
	rand.Seed(time.Now().UnixNano())
	p := rand.Perm(GeneralsNum)
	fmt.Printf("Randomly Chose %d Traitors:", TraitorsNum)
	for _, r := range p[0:TraitorsNum] {
		fmt.Printf("General %d ", r)
		TraitorsMap[r] = 1
	}
	fmt.Print("\n\n")
	return TraitorsMap
}

func opposite(order string) string {
	if order == RETREAT {
		return ATTACK
	} else {
		return RETREAT
	}
}

// This method will simulate sending a order by constructing a node.
// Parameters:
//		ReceiverId(int): The id of the receiver
//		TraitorsMap(map[int]int): A map of traitors.
func (node *Node) sendMessage(ReceiverId int, TraitorsMap map[int]int) *Node {
	order := node.Receive
	// If this general has already seen this order before, don't send.
	if _, ok := node.path[ReceiverId]; ok {
		return nil
	}
	// Check if the sender is a traitor and the receiver id is even.
	if _, ok := TraitorsMap[node.General]; ok && ReceiverId%2 == 0 {
		order = opposite(order)
	}
	path := make(map[int]int)
	// copy 'path' field.
	for k, v := range node.path {
		path[k] = v
	}
	// add receiver to 'path' field
	path[ReceiverId] = 1
	//constructing new node, depth+1,
	return &Node{ReceiverId, order, "", nil, path, node.depth + 1, node}
}

//This method build the whole message sending tree by Completing all order sending processes
func buildTree(GeneralsNum, TraitorsNum int, Order string, TraitorsMap map[int]int) *Node {
	// Construct the root (Commander).
	commander := &Node{0, Order, "", nil, map[int]int{0: 1}, 0, nil}
	if TraitorsMap[0] == 1 {
		fmt.Println("Commander is a traitor")
	}

	queue := []*Node{commander}

	// This is a depth-limited BFS to create nodes,
	// only up to the depth defined as the number of traitors.
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node.depth <= TraitorsNum {
			// Iterate through each general and send the message, add to end of BFS queue.
			for ReveiverId := 1; ReveiverId < GeneralsNum; ReveiverId++ {
				child := node.sendMessage(ReveiverId, TraitorsMap)
				if child == nil {
					continue
				}
				// add new node to childern field. record parent.
				node.children = append(node.children, child)
				child.parent = node
				queue = append(queue, child)
			}
		}
	}
	return commander
}

// This is responsible for deciding the decision of each of the nodes within our created tree.
func (node *Node) decide() string {
	// Base case: If this is a leaf node, then it just follow the order it received
	if len(node.children) == 0 {
		node.Responce = node.Receive
		return node.Responce
	}
	// If it isn't a leaf node, the general make a decision based on how he make decision in
	// the next level(i.e. where he appears as a child of this node's sibling)
	decisions := map[string]int{node.Receive: 1}
	for _, sibling := range node.parent.children {
		if sibling.General != node.General {
			for _, child := range sibling.children {
				if child.General == node.General {
					decision := child.decide()
					if _, ok := decisions[decision]; ok {
						decisions[decision]++
					} else {
						decisions[decision] = 1
					}
				}
			}
		}

	}
	// Get the majority vote
	decision := RETREAT
	if decisions[ATTACK] > decisions[RETREAT] {
		decision = ATTACK
	}
	node.Responce = decision
	return node.Responce
}

func main() {
	flag.Parse()
	// Parse flags in to better variable names
	var GeneralsNum int = *G
	var TraitorsNum int = *M
	var order string = *O
	// Make sure that the number of faulty generals is less than 1/3 of the number of generals
	if 3*TraitorsNum > GeneralsNum {
		fmt.Println("Too many traitors, can only have a 1/3 of the total generals being traitors")
		return
	}
	fmt.Printf("\nThe Order Given is %s\n", order)
	TraitorsTable := getTraitors(GeneralsNum, TraitorsNum)

	// build tree
	tree := buildTree(GeneralsNum, TraitorsNum, order, TraitorsTable)
	// let root make decision(it will recursively go through the whold tree)
	decisions := make(map[string]int)
	for _, General := range tree.children {
		decision := General.decide()
		if _, ok := decisions[decision]; ok {
			decisions[decision]++
		} else {
			decisions[decision] = 1
		}

	}
	// give root's decision
	decision := RETREAT
	if decisions[ATTACK] > decisions[RETREAT] {
		decision = ATTACK
	}
	tree.Responce = decision

	//result printing
	for _, generalNode := range tree.children {
		if _, ok := TraitorsTable[generalNode.General]; ok {
			fmt.Print("Traitor ")
		}
		fmt.Printf("General %d decides on %s\n", generalNode.General, generalNode.Responce)
	}
	fmt.Printf("Consensus decision is %s\n\n", tree.Responce)
}
