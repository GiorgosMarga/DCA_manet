package main

import (
	"flag"
	"log"
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
	g.MakeGraphViz(file)
	g.DCA()
	g.MakeGraphVizClustered(file)
	g.PrintClusters()
}
