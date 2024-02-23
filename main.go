package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"github.com/GiorgosMarga/CDA/graph"
)

var wg *sync.WaitGroup

func main() {
	var file string
	flag.StringVar(&file, "f", "", "file")
	flag.Parse()
	g, err := graph.GraphFromFile(file)
	if err != nil {
		log.Fatal(err)
	}
	g.MakeGraphViz("graph_test.dot")
	g.DCA()
	g.MakeGraphVizClustered("clustered_2.dot")
	g.PrintClusters()
}
