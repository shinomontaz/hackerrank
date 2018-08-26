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

	//	graph := generateExample(10, 20)
	graph := readGraphFromCli()
	//	graph := readMock()

	solution := make([]int, 0, len(graph.nodes))

	if len(graph.nodes) < 10000 {
		for i := 0; i < len(graph.nodes); i++ {
			start := graph.nodes[i]
			clone := graph.Clone()
			solution1 := clone.DFS(start, 20)
			if len(solution1) > len(solution) {
				solution = solution1
			}
		}
	} else {
		solution = rand.Perm(len(graph.nodes))
		for i := 0; i < len(solution); i++ {
			cost1 := graph.cost(solution)
			for j := 0; j < len(solution); j++ {
				if i == j {
					continue
				}
				solution[i], solution[j] = solution[j], solution[i]
				if cost1 < graph.cost(solution) {
					solution[i], solution[j] = solution[j], solution[i]
				}
			}
		}
	}
	/*
		for i := 0; i < len(graph.nodes); i++ {
			start := graph.nodes[i]
			clone := graph.Clone()
			solution1 := clone.DFS(start, 20)
			if len(solution1) > len(solution) {
				solution = solution1
			}
		}
	*/
	fmt.Println(len(solution))
	for _, i := range solution {
		fmt.Printf("%d ", i)
	}

}

func (g *Graph) cost(way []int) (cost float64) {
	for i := 1; i < len(way); i++ {
		exists := false
		for _, conn := range g.nodes[i-1].Connections {
			if conn.ID == way[i] {
				exists = true
				break
			}
		}

		if !exists {
			cost += 1000000.0
		}
	}
	return cost
}

/*
Concurrent experiment:

	count := runtime.GOMAXPROCS(0)
	sem := make(chan struct{}, count)
	var wg sync.WaitGroup
	solutions := make(chan []int, count)
	solution := make([]int, 0)
	for i := 0; i < len(graph.nodes); i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			start := graph.nodes[i]
			clone := graph.Clone()
			solutions <- clone.DFS(start, 20)
			//			if len(solution1) > len(solution) {
			//				solution = solution1
			//			}
			<-sem
		}(i)
	}

	go func() {
		for sol := range solutions {
			if len(sol) > len(solution) {
				solution = sol
			}
		}
	}()

	wg.Wait()
*/

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

func (graph *Graph) DFS(vertex *Point, deep int) (order []int) {
	stack := make([]*Point, 0)
	stack = append(stack, vertex)
	i := 0
	for len(stack) > 0 && i < deep {
		vert := stack[0]
		stack = stack[1:]
		i++
		if !vert.labeled {
			vert.labeled = true
			order = append(order, vert.ID)
			for _, w := range vert.Connections {
				stack = append([]*Point{w}, stack...)
				i++
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
