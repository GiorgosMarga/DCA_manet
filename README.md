# Go-DCA-MANET: Distributed Clustering for Mobile Ad Hoc Networks

A Go program implementing a Distributed Clustering Algorithm (DCA) for efficient node organization within a Mobile Ad Hoc Network (MANET). This implementation uses a dominating set-based approach for cluster formation.

## Features 

* Distributed cluster formation and maintenance 
* Adapts to the dynamic topology of MANETs

## Installation 

```bash
go get github.com/GiorgosMarga/DCA-manet
```
## Usage

**Prerequisites:**

* GraphViz (optional, for visualization)

**Steps:**

1. Create a graph input file (`sample_graph.txt`) according to the format specified below.
2. Build the executable: `go build`
3. Run the program:  `./CDA -f sample_graph.txt`

**Output:**

The program generates the following files:
* **graph_test2.dot:** Contains a GraphViz description of the original graph.
* **graph_test2.png:**  A PNG image of the original graph.
* **clustered_test2.dot:** Contains a GraphViz description of the graph with clusters highlighted.
* **clustered_test2.png:** A PNG image of the graph with clusters visually marked.

## Input File Format

The input file specifies the graph structure. Each line has one of two meanings:

* **Node Definition:**  Lines containing a comma `,` define a node. The format is:
    ```
    <node_id>,<weight>
    ```
    *  **node_id:** A unique numerical identifier for the node.
    * **weight:** A numerical weight associated with the node.

* **Edge Definition:** Lines containing a hyphen `-` define an edge (connection) between two nodes. The format is:
    ```
    <node_id1>-<node_id2>
    ```

**Example:**
1,1
2,2
3,3
1-2
1-3
2-3

This example defines three nodes (1, 2, and 3) with weights of 1, 2, and 3 respectively. It also creates edges between nodes 1-2, 1-3, and 2-3.


## Code Example

Here's a basic example demonstrating how to use the package in your `main.go` file:

```go
package main

import (
    "log"

    "github.com/GiorgosMarga/DCA_manet/graph [invalid URL removed]"
)

func main() {
    f := "graph.txt" // Assuming your input file is named 'graph.txt'
    g, err := graph.GraphFromFile(f)
    if err != nil {
        log.Fatal(err)
    }
    g.MakeGraphViz(f)        // Generates the original graph visualization
    g.DCA()                  // Performs the Distributed Clustering Algorithm
    g.MakeGraphVizClustered(f) // Generates the clustered graph visualization
}
