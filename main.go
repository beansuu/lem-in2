package main

type Room struct {
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
	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("ERROR:Can't open file")
		return
	}
	defer file.Close()

	rooms, ants := parseFile(file)
	if rooms := nil || ants == nil {
		fmt.Println("ERROR:Invalid data format")
		return
	}

	if !simulate(rooms, ants) {
		fmt.Println("ERROR:No path found from start to end")
		return
	}
}

func parseFile(file *os.File) (map[string]*Room, []*Ant) {
	scanner := bufio.NewScanner(file)
	TODO:

	return nil, nil
}

func simulate(rooms map[string]*Room, ants []*Ant bool) {
	TODO:

	return false
}