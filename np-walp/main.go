package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	graph := generateExample(10, 20)

	for _, node := range graph.nodes {
		fmt.Println(*node, node.Connections)
	}

	// А теперь решаем!
	start := rand.Intn(10)

}

func generateExample(N, M int) *Graph {
	nodes := generatePoints(N)
	//	edges := make([]*Edge, 0, M)
	for i := 0; i < len(nodes)-1; i++ {
		//		edges = append(edges, &Edge{Start: nodes[i], End: nodes[i+1]})
		nodes[i].Connections = append(nodes[i].Connections, nodes[i+1])
		nodes[i+1].Connections = append(nodes[i+1].Connections, nodes[i])
	}
	for i := 0; i < M-N; i++ {
		idx1 := rand.Intn(len(nodes))

		connections := make([]int, 0)
		for _, conn := range nodes[idx1].Connections {
			for idx, node := range nodes {
				if node == conn {
					connections = append(connections, idx)
				}
			}
		}
		if len(nodes)-len(connections) == 0 {
			continue
		}
		sort.Ints(connections)
		idx2 := rand.Intn(len(nodes) - len(connections))

		if idx2 > connections[0] {
			shift := 0
			for _, connIdx := range connections {
				if idx2 < connIdx {
					idx2 += shift
					break
				}
				if idx2 == connIdx {
					shift++
				}
			}
		}

		nodes[idx1].Connections = append(nodes[idx1].Connections, nodes[idx2])
		nodes[idx2].Connections = append(nodes[idx2].Connections, nodes[idx1])
	}
	return &Graph{nodes: nodes}
}

type Graph struct {
	nodes []*Point
	//	edges []*Edge
}

/*
type Edge struct {
	Start *Point
	End   *Point
}
*/
type Point struct {
	Lat         float64
	Lng         float64
	Connections []*Point
}

func generatePoints(n int) []*Point {
	res := make([]*Point, 0)

	for i := 1; i < n; i++ {
		res = append(res, &Point{
			Lat:         rand.Float64() * 100,
			Lng:         rand.Float64() * 100,
			Connections: make([]*Point, 0),
		})
	}

	return res
}
