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

**Input File Format:**
*  [Describe the format 

