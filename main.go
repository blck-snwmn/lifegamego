package main

import (
	"fmt"
	"time"
)

const csi = "\033["

// State is life game cell state
type State bool

const (
	dead  State = false
	alive State = true
)

// Cell is life game cell
type Cell struct {
	from  []<-chan State
	to    []chan<- State
	state State
}

func (c *Cell) changeState(count int) {
	defer func() {
		for _, ch := range c.to {
			close(ch)
		}
	}()
	// send initial state
	for _, ch := range c.to {
		ch <- c.state
	}

	for i := 0; i < count; i++ {
		count := 0
		for _, ch := range c.from {
			if <-ch {
				count++
			}
		}
		if c.state {
			//alive
			switch count {
			case 2, 3:
				c.state = alive
			case 0, 1:
				c.state = dead
			default:
				c.state = dead
			}

		} else if count == 3 {
			//dead
			c.state = alive
		}

		for _, ch := range c.to {
			ch <- c.state
		}
	}
}

// Cells is cell array
type Cells [][]Cell

// NewLifeGame return LifeGame and channels to show
func NewLifeGame(width, height int) (*LifeGame, [][]<-chan State) {
	lg := &LifeGame{}
	lg.tickNum = 10
	// init
	cells := make(Cells, height)

	seq := 0
	// set Cell
	for i := 0; i < height; i++ {
		cells[i] = make([]Cell, width)

		for j := 0; j < width; j++ {
			cells[i][j] = Cell{[]<-chan State{}, []chan<- State{}, false}
			seq++
		}
	}
	lg.Cells = cells
	dwr := lg.genCells(width, height)
	return lg, dwr
}

// LifeGame manage cells for life game
type LifeGame struct {
	Cells   [][]Cell //Temporarily public
	tickNum int
}

// Start start life game
func (lg *LifeGame) Start() {
	for i, c := range lg.Cells {
		for j := range c {
			go lg.Cells[i][j].changeState(lg.tickNum)
		}
	}
}

func (lg *LifeGame) genCells(width, height int) [][]<-chan State {
	drawer := make([][]<-chan State, height)
	// set chan
	for i := 0; i < height; i++ {
		drawer[i] = make([]<-chan State, width)

		for j := 0; j < width; j++ {

			top := i - 1
			btm := i + 1
			left := j - 1
			right := j + 1

			isInTop := top >= 0
			isInBtm := btm < height
			isInRight := right < width
			isInLeft := left >= 0

			// top
			lg.genCell(i, j, top, j, isInTop)
			// bottom
			lg.genCell(i, j, btm, j, isInBtm)
			// left
			lg.genCell(i, j, i, left, isInLeft)
			// right
			lg.genCell(i, j, i, right, isInRight)
			// top-left
			lg.genCell(i, j, top, left, isInTop && isInLeft)
			// top-right
			lg.genCell(i, j, top, right, isInTop && isInRight)
			// bottom-left
			lg.genCell(i, j, btm, left, isInBtm && isInLeft)
			// bottom-right
			lg.genCell(i, j, btm, right, isInBtm && isInRight)
			{
				// drawer
				c := make(chan State, 1)
				drawer[i][j] = c
				lg.Cells[i][j].to = append(lg.Cells[i][j].to, c)
			}
		}
	}

	return drawer
}

func (lg *LifeGame) genCell(sl, sc, dl, dc int, cond bool) {
	if cond {
		c := make(chan State, 1)
		lg.Cells[dl][dc].from = append(lg.Cells[dl][dc].from, c)
		lg.Cells[sl][sc].to = append(lg.Cells[sl][sc].to, c)
	}
}

func main() {
	// tickNum := 1
	height := 3
	width := 3
	cs, d := NewLifeGame(height, width)
	cs.Cells[1][0].state = alive
	cs.Cells[1][1].state = alive
	cs.Cells[1][2].state = alive

	cs.Start()

	for i := 0; i < 10; i++ {
		for _, line := range d {
			for _, column := range line {
				s := <-column
				var v string
				if s {
					v = csi + "47m " + csi + "0m"
				} else {
					v = csi + "8m " + csi + "0m"
				}
				fmt.Print(v)
			}
			fmt.Print("\n")
		}
		time.Sleep(time.Second)
		fmt.Printf(csi+"%dF", height)
	}
	fmt.Printf(csi+"%dE", height)
}
