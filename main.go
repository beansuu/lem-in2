package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const INF int = 1 << 30

type Room struct {
	name string
	id   int
}

type AntNest struct {
	rooms    map[string]Room
	adj      [][]int
	capacity [][]int
	flow     [][]int
	parent   []int
	start    int
	end      int
	ants     int
}

type Ant struct {
	id       int
	position int
	path     []int
}

func CreateAntNest() *AntNest {
	return &AntNest{
		rooms: make(map[string]Room),
		adj:   make([][]int, 0),
	}
}

func (a *AntNest) AddRoom(name string) {
	id := len(a.rooms)
	a.rooms[name] = Room{name, id}

	// Expand capacity, flow, and adjacency matrix
	for i := range a.adj {
		a.capacity[i] = append(a.capacity[i], 0)
		a.flow[i] = append(a.flow[i], 0)
	}
	a.capacity = append(a.capacity, make([]int, len(a.rooms)))
	a.flow = append(a.flow, make([]int, len(a.rooms)))
	a.adj = append(a.adj, []int{})
}

func (nest *AntNest) AddTunnel(fromName, toName string) {
	from, ok1 := nest.rooms[fromName]
	to, ok2 := nest.rooms[toName]

	if !ok1 || !ok2 {
		fmt.Printf("Error adding tunnel between %s and %s: Rooms not found\n", fromName, toName)
		return
	}

	nest.capacity[from.id][to.id] = 1
	nest.capacity[to.id][from.id] = 1
	nest.adj[from.id] = append(nest.adj[from.id], to.id)
	nest.adj[to.id] = append(nest.adj[to.id], from.id)

	fmt.Printf("Added tunnel between %s(%d) and %s(%d)\n", fromName, from.id, toName, to.id) // Debugging line
}

func (a *AntNest) BFS(source, sink int) bool {
	visited := make([]bool, len(a.rooms))
	a.parent = make([]int, len(a.rooms))
	for i := range a.parent {
		a.parent[i] = -1
	}

	queue := []int{source}
	visited[source] = true

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		for _, v := range a.adj[u] {
			if !visited[v] && a.capacity[u][v]-a.flow[u][v] > 0 {
				visited[v] = true
				a.parent[v] = u
				queue = append(queue, v)

				if v == sink {
					return true
				}
			}
		}
	}

	return false
}

func (a *AntNest) FordFulkerson(source, sink int) int {
	maxFlow := 0

	for a.BFS(source, sink) {
		pathFlow := INF
		s := sink

		for s != source {
			pathFlow = min(pathFlow, a.capacity[a.parent[s]][s]-a.flow[a.parent[s]][s])
			s = a.parent[s]
		}

		maxFlow += pathFlow
		v := sink

		for v != source {
			u := a.parent[v]
			a.flow[u][v] += pathFlow
			a.flow[v][u] -= pathFlow
			v = u
		}
	}

	return maxFlow
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ParseFile(filename string) *AntNest {
	nest := CreateAntNest()

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	nest.ants = n

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Parsing line:", line) // Debugging line
		switch {
		case line == "##start":
			scanner.Scan()
			parts := strings.Split(scanner.Text(), " ")
			nest.AddRoom(parts[0])
			nest.start = nest.rooms[parts[0]].id
			fmt.Println("Added start room:", parts[0]) // Debugging line
		case line == "##end":
			scanner.Scan()
			parts := strings.Split(scanner.Text(), " ")
			nest.AddRoom(parts[0])
			nest.end = nest.rooms[parts[0]].id
			fmt.Println("Added end room:", parts[0]) // Debugging line
		case strings.Contains(line, "-"):
			parts := strings.Split(line, "-")
			nest.AddTunnel(parts[0], parts[1])
		default:
			// Regular room definition
			parts := strings.Split(line, " ")
			nest.AddRoom(parts[0])
			fmt.Println("Added room:", parts[0]) // Debugging line
		}
	}
	return nest
}

func (a *AntNest) ExtractPathsFromResidual() [][]int {
	var paths [][]int
	visited := make([]bool, len(a.rooms))
	var path []int

	var dfs func(u int) bool
	dfs = func(u int) bool {
		if u == a.end {
			paths = append(paths, append([]int(nil), path...))
			return true
		}

		visited[u] = true
		for _, v := range a.adj[u] {
			residualCapacity := a.capacity[u][v] - a.flow[u][v]
			fmt.Printf("From %d to %d: Capacity=%d, Flow=%d, Residual=%d\n", u, v, a.capacity[u][v], a.flow[u][v], residualCapacity) // Debug line
			if residualCapacity > 0 && !visited[v] {
				path = append(path, v)
				if dfs(v) {
					a.flow[u][v]--
					a.flow[v][u]++
					return true
				}
				path = path[:len(path)-1]
			}
		}
		return false
	}

	for dfs(a.start) {
		path = []int{a.start}
		visited = make([]bool, len(a.rooms))
	}

	return paths
}

func (a *AntNest) MoveAnts(paths [][]int) {
	if len(paths) == 0 {
		fmt.Println("No paths available for ants.")
		return
	}

	ants := make([]Ant, a.ants)
	sortPathsByLength(paths)

	// Create an inverse mapping for rooms for efficient look-up
	roomNames := make(map[int]string)
	for name, room := range a.rooms {
		roomNames[room.id] = name
	}

	for i := range ants {
		ants[i] = Ant{
			id:       i + 1,
			position: -1,
			path:     paths[i],
		}
	}

	allReachedEnd := false
	for !allReachedEnd {
		moves := []string{}
		allReachedEnd = true

		for i := range ants {
			ant := &ants[i]

			if ant.position == -1 && (i == 0 || ants[i-1].position > 0) {
				ant.position++
			} else if ant.position > -1 && ant.position < len(ant.path)-1 {
				ant.position++
			}

			if ant.position >= 0 && ant.position < len(ant.path) {
				roomName := roomNames[ant.path[ant.position]]
				moves = append(moves, fmt.Sprintf("L%d-%s", ant.id, roomName))
			}

			if ant.position != len(ant.path)-1 {
				allReachedEnd = false
			}
		}

		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		}
	}
}

func sortPathsByLength(paths [][]int) {
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the path to the input file.")
		return
	}

	filename := os.Args[1]
	nest := ParseFile(filename)
	if nest == nil {
		fmt.Println("Error reading the input file.")
		return
	}

	maxAnts := nest.FordFulkerson(nest.start, nest.end)
	fmt.Println("Maximum number of ants that can travel without colliding:", maxAnts)

	paths := nest.ExtractPathsFromResidual()

	nest.MoveAnts(paths)
}
