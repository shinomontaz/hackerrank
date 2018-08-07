package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	//	graph := generateExample(10, 20)
	graph := readGraphFromCli()

	for _, node := range graph.nodes {
		fmt.Println(*node, node.Connections)
		for _, conn := range node.Connections {
			fmt.Println("-", conn)
		}
	}

	// А теперь решаем!
	start := graph.nodes[1 /*rand.Intn(len(graph.nodes))*/]

	visited := DFS(graph, start)
	fmt.Println(len(visited))

	for _, i := range visited {
		fmt.Printf("%d", i)
	}

}

func generateExample(N, M int) *Graph {
	nodes := generatePoints(N)

	perm := rand.Perm(len(nodes) - 1)

	for _, i := range perm {
		nodes[i].Connections = append(nodes[i].Connections, nodes[i+1])
		nodes[i+1].Connections = append(nodes[i+1].Connections, nodes[i])
	}
	for i := 0; i < M-N; i++ {
		idx1 := rand.Intn(len(nodes))

		connections := make([]int, 0)
		for _, conn := range nodes[idx1].Connections {
			connections = append(connections, conn.ID)
		}

		if len(nodes)-len(connections) == 0 {
			continue
		}
		sort.Ints(connections)
		idx2 := rand.Intn(len(nodes) - len(connections))

		found := false
		for _, conn := range nodes[idx1].Connections {
			if conn == nodes[idx2] || conn == nodes[idx1] {
				found = true
				break
			}
		}

		if !found {
			nodes[idx1].Connections = append(nodes[idx1].Connections, nodes[idx2])
		}
		found = false
		for _, conn := range nodes[idx2].Connections {
			if conn == nodes[idx1] || conn == nodes[idx2] {
				found = true
				break
			}
		}

		if !found {
			nodes[idx2].Connections = append(nodes[idx2].Connections, nodes[idx2])
		}
	}
	return &Graph{nodes: nodes}
}

type Graph struct {
	nodes []*Point
}

type Point struct {
	ID          int
	Connections []*Point
	labeled     bool
}

func generatePoints(n int) []*Point {
	res := make([]*Point, 0)

	for i := 1; i < n; i++ {
		res = append(res, &Point{
			ID:          i,
			Connections: make([]*Point, 0),
		})
	}

	return res
}

func DFS(graph *Graph, vertex *Point) (order []int) {
	stack := make([]*Point, 0)
	stack = append(stack, vertex)

	for len(stack) > 0 {
		vert := stack[0]
		stack = stack[1:]
		if !vert.labeled {
			vert.labeled = true
			order = append(order, vert.ID)
			for _, w := range vert.Connections {
				stack = append([]*Point{w}, stack...)
			}
		}

	}

	return order
}

func readGraphFromCli() *Graph {
	var n, m int
	if _, err := fmt.Scanf("%d %d", &n, &m); err != nil {
		log.Fatalf("can'r read input: %s", err)
	}

	nodes := generatePoints(n)

	var from, to int
	for i := 0; i < m; i++ {
		if _, err := fmt.Scanf("%d %d", &from, &to); err != nil {
			log.Fatalf("can't read input: %s", err)
		}
		nodes[from-1].Connections = append(nodes[from-1].Connections, nodes[to-1])
		nodes[to-1].Connections = append(nodes[to-1].Connections, nodes[from-1])
	}

	return &Graph{nodes: nodes}
}
