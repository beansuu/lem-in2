package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type AntNest struct {
	ants    int
	tunnels map[string][]string
	start   string
	end     string
}

type Ant struct {
	id       int
	position int
	path     []string
}

func ParseFile(filename string) *AntNest {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read number of ants
	scanner.Scan()
	ants, _ := strconv.Atoi(scanner.Text())
	nest := &AntNest{
		ants:    ants,
		tunnels: make(map[string][]string),
	}

	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case line == "##start":
			scanner.Scan()
			parts := strings.Split(scanner.Text(), " ")
			nest.start = parts[0]

		case line == "##end":
			scanner.Scan()
			parts := strings.Split(scanner.Text(), " ")
			nest.end = parts[0]

		case strings.Contains(line, "-"):
			parts := strings.Split(line, "-")
			nest.AddTunnel(parts[0], parts[1])
		}
	}

	return nest
}

func (a *AntNest) AddTunnel(from, to string) {
	if _, exists := a.tunnels[from]; !exists {
		a.tunnels[from] = []string{}
	}
	a.tunnels[from] = append(a.tunnels[from], to)

	if _, exists := a.tunnels[to]; !exists {
		a.tunnels[to] = []string{}
	}
	a.tunnels[to] = append(a.tunnels[to], from)
}

func (a *AntNest) BFS() []string {
	queue := [][]string{{a.start}}
	visited := make(map[string]bool)

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		lastRoom := path[len(path)-1]

		if lastRoom == a.end {
			return path
		}

		for _, nextRoom := range a.tunnels[lastRoom] {
			if !visited[nextRoom] || nextRoom == a.end {
				newPath := make([]string, len(path))
				copy(newPath, path)
				newPath = append(newPath, nextRoom)
				queue = append(queue, newPath)

				if nextRoom != a.end {
					visited[nextRoom] = true
				}
			}
		}
	}

	return nil
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

	paths := nest.MultipleBFS()
	nest.MoveAnts(paths)
}

func (a *AntNest) MoveAnts(paths [][]string) {
	ants := make([]Ant, a.ants)

	sortPathsByLength(paths)

	for i := range ants {
		ants[i] = Ant{
			id:       i + 1,
			position: -1,
			path:     paths[i%len(paths)],
		}
	}

	allReachedEnd := false
	for !allReachedEnd {
		moves := []string{}
		occupiedRooms := map[string]bool{}
		allReachedEnd = true

		for i := range ants {
			ant := &ants[i] // Use pointer to update ants slice directly
			if ant.position < len(ant.path)-1 {
				allReachedEnd = false
				nextRoom := ant.path[ant.position+1]
				if !occupiedRooms[nextRoom] {
					if ant.position != -1 { // If the ant has already started
						delete(occupiedRooms, ant.path[ant.position])
					}
					ant.position++
					moves = append(moves, fmt.Sprintf("L%d-%s", ant.id, nextRoom))
					occupiedRooms[nextRoom] = true
				}
			}
		}

		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		}
	}
}

func (a *AntNest) MultipleBFS() [][]string {
	allPaths := [][]string{}
	tunnelsBackup := a.copyTunnels() // Backup the tunnels

	for {
		path := a.BFS()
		if path == nil || len(path) == 0 {
			break
		}

		allPaths = append(allPaths, path)

		// Temporarily remove the found path's tunnels from the nest
		for i := 0; i < len(path)-1; i++ {
			room := path[i]
			nextRoom := path[i+1]
			a.tunnels[room] = removeFromStringSlice(a.tunnels[room], nextRoom)
			a.tunnels[nextRoom] = removeFromStringSlice(a.tunnels[nextRoom], room)
		}
	}

	a.tunnels = tunnelsBackup // Restore the tunnels

	return allPaths
}

func removeFromStringSlice(slice []string, s string) []string {
	index := -1
	for i, item := range slice {
		if item == s {
			index = i
			break
		}
	}
	if index == -1 {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

func sortPathsByLength(paths [][]string) {
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})
}

func (a *AntNest) copyTunnels() map[string][]string {
	tunnelsCopy := make(map[string][]string)
	for key, val := range a.tunnels {
		tunnelsCopy[key] = append([]string(nil), val...)
	}
	return tunnelsCopy
}
