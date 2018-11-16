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

var G = flag.Int("G", 10, "The number of generals. First is commander, rest are lieutenants")

var M = flag.Int("M", 3, "The number of faulty generals")

var O = flag.String("O", "ATTACK", "The order by the general(ATTACK or RETREAT)")

// Values:
//	inputValue(string): The order the node has been given
//	outputValue(string): The order the node decides on
//	processIDs(map[int]int): Keeps track of which nodes have seen this message before
//	children([]*Node): Keeps track of the children of the node
//	id(int): The id of the corresponding general
type Node struct {
	General  int
	Receive  string
	Responce string
	children []*Node
	path     map[int]int
	depth    int
	parent   *Node
}

// Parameters:
//		numFaultyGenerals(int): The number of faulty generals
//		numGenerals(int): The total number of generals

// This function will randomly select numFaultyGenerals generals to be the faulty ones.
// The commander can also be chosen.
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

// Parameters:
//		id(int): The id of the node you want to message
//		faultyGenerals(map[int]int): A map of all the generals that are known to be faulty.

// This function will simulate sending a message by updating the information of the recieving node.
// It will then add the recieivng node as one of its children which will be used later in the decision.
func (node *Node) sendMessage(ReveiverId int, TraitorsMap map[int]int) *Node {
	order := node.Receive
	// If this general has already recieved a message from us, don't send.
	if _, ok := node.path[ReveiverId]; ok {
		return nil
	}
	// Check if the general sending the node is faulty
	if _, ok := TraitorsMap[node.General]; ok && ReveiverId%2 == 0 {
		order = opposite(order)
	}
	path := make(map[int]int)
	// send them all the process id's we have too
	for k, v := range node.path {
		path[k] = v
	}
	// add this time
	path[ReveiverId] = 1
	return &Node{ReveiverId, order, "", nil, path, node.depth + 1, node}
}

func buildTree(GeneralsNum, TraitorsNum int, Order string, TraitorsMap map[int]int) *Node {
	// This will ensure the commander is never messaged again.
	commander := &Node{0, Order, "", nil, map[int]int{0: 1}, 0, nil}
	if TraitorsMap[0] == 1 {
		fmt.Println("Commander is a traitor")
	}

	queue := []*Node{commander}

	// This is a depth-limited breadth first search so that we sendMessages and create nodes for all generals,
	// only up to the depth defined as the number of faulty generals.
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
	// Base case: If this is a leaf node, then the output value is the input value
	if len(node.children) == 0 {
		node.Responce = node.Receive
		return node.Responce
	}
	// If it isn't a leaf node it means we have to rely on the children of siblings.
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
	// Get the majority vote from this nodes children and assign it to the output value.
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
	// Make sure that the number of faulty generals is less than a third of the number of generals
	if 3*TraitorsNum > GeneralsNum {
		fmt.Println("Too many faulty generals, can only have a third of the total generals being faulty")
		return
	}
	fmt.Printf("The Order Given is %s\n", order)
	TraitorsTable := getTraitors(GeneralsNum, TraitorsNum)

	// Call the orchestrator
	tree := buildTree(GeneralsNum, TraitorsNum, order, TraitorsTable)
	decisions := make(map[string]int)
	for _, General := range tree.children {
		decision := General.decide()
		if _, ok := decisions[decision]; ok {
			decisions[decision]++
		} else {
			decisions[decision] = 1
		}

	}

	decision := RETREAT
	if decisions[ATTACK] > decisions[RETREAT] {
		decision = ATTACK
	}
	tree.Responce = decision

	for _, generalNode := range tree.children {
		if _, ok := TraitorsTable[generalNode.General]; ok {
			fmt.Print("Traitor ")
		}
		fmt.Printf("General %d decides on %s\n", generalNode.General, generalNode.Responce)
	}
	fmt.Printf("Consensus decision is %s", tree.Responce)
}
