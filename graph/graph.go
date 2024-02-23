package graph

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const (
	maxChanLen = 100
	colorsLen  = 12
)

var wg *sync.WaitGroup
var colors []string = []string{"green", "blue", "black", "red", "yellow", "gray", "brown", "purple", "turquoise", "sienna", "orange", "pink"}

type Graph struct {
	Nodes       map[int]*Node
	isClustered bool
}

type Node struct {
	Id        int
	Weight    int
	Neighbors []*Node
	// This channel is for receiving CH messages.
	// This channel is used in neighbors CHMapChan
	CHChan chan CHMessage
	// Connect curr node with CHChan of each neighbor
	CHMapChan  map[int](chan CHMessage)
	CHMessages map[int]CHMessage
	// This channel is for receiving JOIN messages.
	// This channel is used in neighbors JOINMapChan
	JOINChan chan JOINMessage
	// Connect curr node with JOINChan of each neighbor
	JOINMapChan     map[int](chan JOINMessage)
	JOINMessages    map[int]JOINMessage
	IsClusterhead   bool
	CandidateId     int
	CandidateWeight int
	BelongsTo       int
	terminate       bool
}

type CHMessage struct {
	from   int
	weight int
}
type JOINMessage struct {
	from int
	to   int
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[int]*Node),
	}
}

func (g *Graph) AddNode(weight int) error {
	if weight < 0 {
		return fmt.Errorf("weight should be positive")
	}
	for _, v := range g.Nodes {
		if v.Weight == weight {
			return fmt.Errorf("weight should be unique. Node with same weight found (%+v)", v)
		}
	}
	id := len(g.Nodes) + 1
	n := &Node{
		Id:           id,
		Weight:       weight,
		Neighbors:    make([]*Node, 0),
		CHMapChan:    make(map[int]chan CHMessage),
		CHChan:       make(chan CHMessage, maxChanLen),
		JOINChan:     make(chan JOINMessage, maxChanLen),
		JOINMapChan:  make(map[int]chan JOINMessage),
		JOINMessages: make(map[int]JOINMessage),
		CHMessages:   make(map[int]CHMessage),
	}
	g.Nodes[id] = n
	return nil
}

func (g *Graph) ConnectNodes(id1, id2 int) error {
	if id2 <= id1 {
		return fmt.Errorf("id1 should be bigger than id2")
	}
	n1, ok := g.Nodes[id1]
	if !ok {
		return fmt.Errorf("node with id (%d) doesnt exist", id1)
	}
	n2, ok := g.Nodes[id2]
	if !ok {
		return fmt.Errorf("node with id (%d) doesnt exist", id2)
	}
	n1.Neighbors = append(n1.Neighbors, n2)
	n2.Neighbors = append(n2.Neighbors, n1)
	n1.CHMapChan[n2.Id] = n2.CHChan
	n2.CHMapChan[n1.Id] = n1.CHChan
	n1.JOINMapChan[n2.Id] = n2.JOINChan
	n2.JOINMapChan[n1.Id] = n1.JOINChan
	return nil
}
func GraphFromFile(filename string) (*Graph, error) {
	g := NewGraph()
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	content := strings.Split(string(file), "\n")
	for _, line := range content {
		if strings.Contains(line, ",") {
			sWeight := strings.Split(line, ",")[1]
			iWeight, err := strconv.Atoi(sWeight)
			if err != nil {
				return nil, err
			}
			err = g.AddNode(iWeight)
			if err != nil {
				return nil, err
			}
		} else if strings.Contains(line, "-") {
			temp := strings.Split(line, "-")
			id1, err := strconv.Atoi(temp[0])
			if err != nil {
				return nil, err
			}
			id2, err := strconv.Atoi(temp[1])
			if err != nil {
				return nil, err
			}
			err = g.ConnectNodes(id1, id2)
		} else {
			return nil, fmt.Errorf("invalid file format")
		}
	}
	return g, nil
}

func (g *Graph) MakeGraphViz(filename string) error {
	filename = strings.Split(filename, ".")[0]
	graphName := fmt.Sprintf("graph_%s.dot", filename)
	f, err := os.Create(graphName)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("graph {\n")
	for _, v := range g.Nodes {
		for _, neighbor := range v.Neighbors {
			if v.Id < neighbor.Id {
				line := fmt.Sprintf("	%d--%d;\n", v.Id, neighbor.Id)
				_, err := f.WriteString(line)
				if err != nil {
					return err
				}
			}
		}
	}
	f.WriteString("}")
	output := fmt.Sprintf("%s.png", filename)
	cmd := exec.Command("neato", "-Tpng", filename, "-o", output)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (g *Graph) DCA() {
	g.isClustered = false
	wg = &sync.WaitGroup{}
	for _, v := range g.Nodes {
		wg.Add(1)
		go func(v *Node) {
			v.Init()
			wg.Done()
		}(v)
	}
	wg.Wait()
	g.isClustered = true
}

func (g *Graph) MakeGraphVizClustered(filename string) error {
	if !g.isClustered {
		return fmt.Errorf("graph is not clustered. Call g.DCA() first")
	}
	filename = strings.Split(filename, ".")[0]

	clusteredFile := fmt.Sprintf("clustered_%s.dot", filename)
	f, err := os.Create(clusteredFile)
	if err != nil {
		return err
	}
	f.WriteString("graph {\nlayout=\"fdp\" sep=\"10\"\n")
	for _, v := range g.Nodes {
		var color string
		if v.IsClusterhead {
			color = colors[v.Id%colorsLen]
		} else {
			color = colors[v.BelongsTo%colorsLen]
		}
		attr := fmt.Sprintf("%d [label=\"%d,%d\",fillcolor=%s,fontcolor=white,style=filled];\n", v.Id, v.Id, v.Weight, color)
		f.WriteString(attr)

	}
	for _, v := range g.Nodes {
		for _, neighbor := range v.Neighbors {
			if v.Id < neighbor.Id {
				line := fmt.Sprintf("	%d--%d;\n", v.Id, neighbor.Id)
				_, err := f.WriteString(line)
				if err != nil {
					return err
				}
			}
		}
	}
	f.WriteString("}")
	output := fmt.Sprintf("clustered_%s.png", filename)
	fmt.Println(output)
	cmd := exec.Command("neato", "-Tpng", filename, "-o", output)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}
	return nil
}
func (g *Graph) Print() {
	for k, v := range g.Nodes {
		fmt.Printf("%d -> ", k)
		for _, neighbor := range v.Neighbors {
			fmt.Printf("(%d,%d) ", neighbor.Id, neighbor.Weight)
		}
		fmt.Println()
	}
}
func (g *Graph) PrintClusters() {
	clusters := make(map[int][]*Node)
	clusters[0] = make([]*Node, 0)
	for _, v := range g.Nodes {
		if v.IsClusterhead {
			clusters[0] = append(clusters[0], v)
			continue
		}
		if _, ok := clusters[v.BelongsTo]; !ok {
			clusters[v.BelongsTo] = make([]*Node, 0)
		}
		clusters[v.BelongsTo] = append(clusters[v.BelongsTo], v)
	}
	for k, v := range clusters {
		fmt.Printf("Cluster %d: ", k)
		for _, node := range v {
			fmt.Printf("(%d,%d) ", node.Id, node.Weight)
		}
		fmt.Println()
	}
}
func (n *Node) Init() {
	fmt.Printf("[%d] init\n", n.Id)
	if n.isClusterhead() {
		n.broadcastCH()
	}
	for {
		if n.terminate {
			fmt.Printf("[%d] has stopped\n", n.Id)
			return
		}
		select {
		case chmsg := <-n.CHChan:
			n.onReceivingCH(chmsg)
		case joinmsg := <-n.JOINChan:
			n.onReceivingJOIN(joinmsg)
		}
	}

}

func (n *Node) broadcastCH() {
	msg := CHMessage{
		from:   n.Id,
		weight: n.Weight,
	}
	for _, v := range n.Neighbors {
		n.CHMapChan[v.Id] <- msg
	}
}
func (n *Node) broadcastJOIN(toID int) {
	fmt.Printf("[%d] sent JOIN message to (%d)\n", n.Id, toID)

	msg := JOINMessage{
		from: n.Id,
		to:   toID,
	}
	for _, v := range n.Neighbors {
		n.JOINMapChan[v.Id] <- msg
	}
}
func (n *Node) isClusterhead() bool {
	for _, v := range n.Neighbors {
		if v.Weight > n.Weight {
			return false
		}
	}
	n.IsClusterhead = true
	fmt.Printf("[%d] is clusterhead\n", n.Id)
	return true
}

func (n *Node) onReceivingCH(ch CHMessage) {
	fmt.Printf("[%d] received CH message %+v\n", n.Id, ch)
	n.CHMessages[ch.from] = ch
	var (
		fromID     = ch.from
		fromWeight = ch.weight
		sendJoin   = true
	)
	for _, v := range n.Neighbors {
		if v.Id != fromID && v.Weight > fromWeight {
			if _, ok := n.JOINMessages[v.Id]; !ok {
				sendJoin = false
			}
		}
	}
	if sendJoin {
		n.BelongsTo = fromID
		n.broadcastJOIN(fromID)
		n.terminate = true
		return
	}

	for _, v := range n.Neighbors {
		if v.Id != fromID && v.Weight > fromWeight {
			_, ok1 := n.JOINMessages[v.Id]
			_, ok2 := n.CHMessages[v.Id]
			if !ok1 && !ok2 {
				if n.CandidateWeight < v.Weight {
					n.CandidateId = v.Id
					n.CandidateWeight = v.Weight
				}
				return
			}
		}
	}
}

func (n *Node) onReceivingJOIN(join JOINMessage) {
	fmt.Printf("[%d] received JOIN message %+v\n", n.Id, join)
	n.JOINMessages[join.from] = join
	if n.IsClusterhead {
		// n had sent a CH message
		if join.to == n.Id {
			fmt.Printf("node (%d) joined [%d]\n", join.from, join.to)
		}
		canStop := true
		for _, v := range n.Neighbors {
			if v.Weight < n.Weight {
				if _, ok := n.JOINMessages[v.Id]; !ok {
					canStop = false
					break
				}
			}
		}
		if canStop {
			n.terminate = true
		}
		return
	}
	for _, v := range n.Neighbors {
		if v.Weight > n.Weight {
			_, ok1 := n.JOINMessages[v.Id]
			_, ok2 := n.CHMessages[v.Id]
			if !ok1 && !ok2 {
				return
			}
		}
	}
	allJoin := true
	receivedCH := false
	for _, v := range n.Neighbors {
		if v.Weight > n.Weight {
			if _, ok := n.JOINMessages[v.Id]; !ok {
				allJoin = false
			}
			if _, ok := n.CHMessages[v.Id]; ok {
				receivedCH = true
			}
		}
	}
	if allJoin {
		fmt.Printf("[%d] is clusterhead\n", n.Id)
		n.IsClusterhead = true
		for _, v := range n.Neighbors {
			if v.Weight < n.Weight {
				if _, ok := n.JOINMessages[v.Id]; !ok {
					n.broadcastCH()
					return
				}
			}
		}
		n.terminate = true
		return
	}
	if receivedCH {
		c := n.FindBiggestIdCluster()
		if c.Weight == 0 {
			return
		}
		n.broadcastJOIN(c.Id)
		n.terminate = true
	}

}

func (n *Node) FindBiggestIdCluster() *Node {
	biggest := &Node{
		Weight: 0,
	}
	for _, v := range n.Neighbors {
		if v.Weight > biggest.Weight && v.IsClusterhead {
			biggest = v
		}
	}
	return biggest
}
