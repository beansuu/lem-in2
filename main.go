package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Room struct { // aa
	name    string
	links   []*Room
	isStart bool
	isEnd   bool
}

type Ant struct {
	id      int
	current *Room
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: Please provide a file name")
		return
	}
	filePath := os.Args[1]

	// open file and close
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("ERROR: Cannot open file %s\n", filePath)
		return
	}
	defer file.Close()

	// parse the file to get rooms and ants
	rooms, ants := parseFile(file)
	if rooms == nil || ants == nil {
		fmt.Println("ERROR: Invalid data format")
		return
	}

	if !simulate(rooms, ants) {
		fmt.Println("ERROR: No path found from start to end")
	}
}

func parseFile(file *os.File) (map[string]*Room, []*Ant) {
	scanner := bufio.NewScanner(file)
	rooms := make(map[string]*Room)
	var startRoom, endRoom *Room
	var numberOfAnts int

	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "##start"):
			scanner.Scan()
			startRoom = createRoom(scanner.Text(), rooms)
		case strings.HasPrefix(line, "##end"):
			scanner.Scan()
			endRoom = createRoom(scanner.Text(), rooms)
		case !strings.HasPrefix(line, "#"):
			if strings.Contains(line, " ") {
				createRoom(line, rooms)
			} else if strings.Contains(line, "-") {
				splitted := strings.Split(line, "-")
				if len(splitted) == 2 {
					room1, exists1 := rooms[splitted[0]]
					room2, exists2 := rooms[splitted[1]]
					if exists1 && exists2 {
						room1.links[room2.name] = room2
						room2.links[room1.name] = room1
					}
				}
			} else if numberOfAnts == 0 {
				fmt.Sscan(line, &numberOfAnts)
			}
		}
	}
	// create ants and place them in the room
	var ants []*Ant
	for i := 1; i <= numberOfAnts; i++ {
		ants = append(ants, &Ant{id: i, room: startRoom})
	}
	return rooms, ants
}

func createRoom(line string, rooms map[string]*Room) *Room {
	splitted := strings.Fields(line)
	if len(splitted) < 3 {
		return nil
	}
	name := splitted[0]
	room := &Room{name: name, links: make(map[string]*Room)}
	rooms[name] = room
	return room
}
