package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var (
	numPits = 1
)

func main() {
	for {
		if !askStartGame(os.Stdin) {
			break
		}
		game := newGame(os.Stdin)
		for !game.gameOver {
			game.doRoomStatus()
			game.doTurn()
			fmt.Println("")
		}
	}
}

func readLine(in io.Reader) (string, error) {
	s := bufio.NewScanner(in)
	s.Split(bufio.ScanLines)
	if !s.Scan() {
		return "", errors.New("could not read next line")
	}
	return s.Text(), nil
}

func askStartGame(in io.Reader) bool {
	fmt.Print("Do you want to start a new game [Yn]? ")
	line, err := readLine(in)
	if err != nil{
		return false
	}
	if line == "Y" || line == "y" || line == "" {
		return true
	}
	return false
}

func newGame(in io.Reader) *game {
	g := &game{
		in: in,
		gameOver: false,
		maze: &mapDGraph{},
		arrows: 3,
		pits: map[int]struct{}{},
	}
	g.maze.addEdge(0, 1, 1)
	g.maze.addEdge(0, 4, 1)
	g.maze.addEdge(1, 2, 1)
	g.maze.addEdge(2, 3, 1)
	g.maze.addEdge(3, 4, 1)
	g.maze.addEdge(3, 0, 1)
	g.maze.addEdge(4, 3, 1)
	g.maze.addEdge(4, 5, 1)
	g.maze.addEdge(5, 1, 1)
	g.maze.addEdge(5, 2, 1)
	g.maze.addEdge(1, 5, 1)

	g.placeWumpus()
	g.placePits()
	return g
}

func (g *game)placeWumpus() {
	g.player = g.randomRoom()
	for i := 0; i < 10; i++ {
		r := g.randomRoom()
		if r != g.player {
			g.wumpus = r
			return
		}
	}
	panic("coud not place wumpus")
}

func (g *game)placePits() {
	for i := 0; i < numPits; i++ {
		placed := false
		for j := 0; j < 10; j++ {
			r := g.randomRoom()
			if r != g.player && !g.hasPitIn(r) {
				g.pits[r] = struct{}{}
				placed = true
				break
			}
		}
		if !placed {
			panic("could not place pit")
		}
	}
}

func (g *game)randomRoom() int {
	return rand.Intn(g.maze.len())
}

type game struct {
	in io.Reader
	gameOver bool
	maze dGraph
	player int
	arrows int
	wumpus int
	pits map[int]struct{}
}

func (g *game)doTurn() {
	action := g.askAction()
	switch action.typ {
	case ACT_MOVE:
		g.doMove(action)
	case ACT_SHOOT:
		g.doShoot(action)
	case ACT_QUIT:
		g.doQuit()
	}
}

func (g *game)doMove(a action) {
	if !g.roomsConnect(g.player, a.moves[0]) {
		panic("could not connect to room.")
	}
	g.player = a.moves[0]
	g.doRoomEnter()
}

func (g *game)doRoomEnter() {
	if g.player == g.wumpus {
		fmt.Println("MUNCH, MUNCH, MUNCH.  You were eaten.")
		g.gameOver = true
		return
	}
	if g.hasPitIn(g.player) {
		fmt.Println("Aaaaaaaaaahhhhhh! You fell into a pit and died.")
		g.gameOver = true
		return
	}
}

func (g *game)doRoomStatus() {
	fmt.Printf("You are in room %d.\n", g.player)

	s := []string{}
	for _, e := range g.maze.edges(g.player) {
		s = append(s, fmt.Sprintf("%d", e.end))
	}
	fmt.Printf("There are exits to room(s): %s\n", strings.Join(s, ", "))
	fmt.Printf("You have %d arrow(s) left.\n", g.arrows)


	wumpus := false
	pit := false
	for _, e := range g.maze.edges(g.player) {
		wumpus = wumpus || e.end == g.wumpus
		pit = pit || g.hasPitIn(e.end)
	}

	if pit {
		fmt.Println("I feel a draft!")
	}
	if wumpus {
		fmt.Println("I smell a wumpus!")
	}
}

func (g *game)doShoot(a action) {
	c := g.player
	for _, n := range a.moves {
		if !g.roomsConnect(c, n) {
			fmt.Println("Splat! Your arrow hit a wall")
			break
		}
		if n == g.wumpus {
			fmt.Println("You got the wumpus!  You WIN!")
			g.gameOver = true
			return
		}
		c = n
	}
	fmt.Println("You missed, and woke the wumpus, and it moves to another room.")
	tunnels := g.maze.edges(g.player)
	moveTo := tunnels[rand.Intn(len(tunnels))].end
	if moveTo == g.player {
		fmt.Println("The wumpus moves into your room, and picks eats you.")
		g.gameOver = true
	}
	g.arrows--
	if g.arrows == 0 {
		fmt.Println("You have no more arrows left.  Eventually the wumpus hunts you down and eats you.")
		g.gameOver = true
	}
}

func (g *game)hasPitIn(room int) bool {
	_, ok := g.pits[room]
	return ok
}

func (g *game)roomsConnect(start, end int) bool {
	for _, e := range g.maze.edges(start) {
		if e.end == end {
			return true
		}
	}
	return false
}

func (g *game)doQuit() {
	fmt.Println("Coward.")
	g.gameOver = true
}

type actionType int

const (
	ACT_MOVE actionType = iota
	ACT_SHOOT
	ACT_QUIT
)

type action struct {
	typ actionType
	moves []int
}

func (g *game)askAction() action {
	for {
		a := action{}
		fmt.Print("Action? ")
		line, err := readLine(g.in)
		if err != nil {
			panic("could not read input")
		}

		s := bufio.NewScanner(strings.NewReader(line))
		s.Split(bufio.ScanWords)
		// if it's an empty line then we move on to the next line
		if (!s.Scan()) {
			continue
		}

		cmd := s.Text()
		switch cmd {
		case "M", "m":
			a.typ = ACT_MOVE
			moves, err := parseRoomList(*s)
			if err != nil {
				fmt.Println("%s", err)
				continue
			}
			if len(moves) != 1 {
				fmt.Printf("You can only use one room")
				continue
			}
			if !g.roomsConnect(g.player, moves[0]) {
				fmt.Printf("The current room does not connect to room %d.\n", moves[0])
				continue
			}
			a.moves = moves
			return a

		case "S", "s":
			a.typ = ACT_SHOOT
			moves, err := parseRoomList(*s)
			if err != nil {
				fmt.Println("%s", err)
				continue
			}
			if len(moves) > 5 {
				fmt.Println("You can ony shoot an arrow up to five moves.")
				continue
			}
			if !g.roomsConnect(g.player, moves[0]) {
				fmt.Printf("The current room does not connect to room %d.\n", moves[0])
				continue
			}
			a.moves = moves
			return a
		case "Q", "q":
			a.typ = ACT_QUIT
			return a
		default:
			fmt.Println("known commads are 'm ROOM', 's ROOM{1,5}', 'q'.")
			continue
		}


	}
	panic("did not parse action correctly")
}


func parseRoomList(s bufio.Scanner) ([]int, error) {
	n := []int{}
	for s.Scan() {
		r, err := strconv.Atoi(s.Text())
		if err != nil {
			return n, fmt.Errorf("Could not parse room number: %s\n", err)
		}
		n = append(n, r)
	}
	if len(n) == 0 {
		return n, errors.New("No room numbers specified.")
	}
	return n, nil
}
