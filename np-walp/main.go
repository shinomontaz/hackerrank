package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	graph := generateExample(10, 20)
	//	graph := readGraphFromCli()
	//	graph := readMock()

	for _, node := range graph.nodes {
		fmt.Println(*node, node.Connections)
		for _, conn := range node.Connections {
			fmt.Println("-", conn)
		}
	}

	// А теперь решаем!
	start := graph.nodes[rand.Intn(len(graph.nodes))]
	graph.nodes[start.ID-1].labeled = true
	branches := make([][]int, 0)
	for _, conn := range start.Connections {
		clone := graph.Clone()
		branches = append(branches, clone.DFS(conn))
	}

	fmt.Println(start.ID, len(branches))

	for _, visited := range branches {
		fmt.Println("-")
		for _, i := range visited {
			fmt.Printf("%d ", i)
		}
	}

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

	for i := 1; i <= n; i++ {
		res = append(res, &Point{
			ID:          i,
			Connections: make([]*Point, 0),
		})
	}

	return res
}

func (graph *Graph) DFS(vertex *Point) (order []int) {
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

func (graph *Graph) Clone() *Graph {
	clone := &Graph{}
	clone.nodes = make([]*Point, 0, len(graph.nodes))
	for _, point := range graph.nodes {
		clone.nodes = append(clone.nodes, &Point{ID: point.ID, labeled: point.labeled, Connections: make([]*Point, 0, len(point.Connections))})
	}

	for i, point := range graph.nodes {
		for _, conn := range point.Connections {
			clone.nodes[i].Connections = append(clone.nodes[i].Connections, clone.nodes[conn.ID-1])
		}
	}

	return clone
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

func readMock() *Graph {
	str := `4 5
3 1
3 4
2 4
2 3
4 1`
	scanner := bufio.NewScanner(strings.NewReader(str))
	var n, m int

	scanner.Scan()
	items := strings.Split(scanner.Text(), " ")
	n, _ = strconv.Atoi(items[0])
	m, _ = strconv.Atoi(items[1])

	nodes := generatePoints(n)

	var from, to int
	i := 0
	for scanner.Scan() {
		if i > m {
			panic("!")
		}
		items := strings.Split(scanner.Text(), " ")

		from, _ = strconv.Atoi(items[0])
		to, _ = strconv.Atoi(items[1])

		nodes[from-1].Connections = append(nodes[from-1].Connections, nodes[to-1])
		nodes[to-1].Connections = append(nodes[to-1].Connections, nodes[from-1])
		i++
	}

	return &Graph{nodes: nodes}

}
