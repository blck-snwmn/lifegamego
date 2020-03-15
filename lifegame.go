package lifegame

import "context"

// State is life game cell state
type State bool

const (
	dead  State = false
	alive State = true
)

func (s State) String() string {
	if s.IsAlive() {
		return "Alive"
	}
	return "Dead"
}

// IsAlive return true if state is alive
func (s State) IsAlive() bool {
	return s == alive
}

func changeState(now State, aliveNum int) State {
	if now.IsAlive() {
		//alive
		switch aliveNum {
		case 2, 3:
			return alive
		default: //0, 1, 4, 5, 6, 7, 8
			return dead
		}
	} else if aliveNum == 3 {
		return alive
	}
	return dead
}

// Cell is life game cell
type Cell struct {
	from  []<-chan State
	to    []chan<- State
	state State
}

// SetAlive set state
func (c *Cell) SetAlive() {
	c.state = alive
}

func (c *Cell) sendState(ctx context.Context) {
	for _, ch := range c.to {
		select {
		case <-ctx.Done():
			return
		case ch <- c.state:
		}
	}
}

func (c *Cell) wake(ctx context.Context, count int) {
	defer func() {
		for _, ch := range c.to {
			close(ch)
		}
	}()
	// send initial state
	c.sendState(ctx)

	for i := 0; i < count; i++ {
		count := 0
		for _, ch := range c.from {
			select {
			case <-ctx.Done():
				return
			case s := <-ch:
				if s.IsAlive() {
					count++
				}
			}
		}
		c.state = changeState(c.state, count)
		c.sendState(ctx)
	}
}

// Cells is cell array
type Cells [][]Cell

// New return LifeGame and channels to show
func New(width, height int) (*LifeGame, [][]<-chan State) {
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
func (lg *LifeGame) Start(ctx context.Context) {
	for i, c := range lg.Cells {
		for j := range c {
			go lg.Cells[i][j].wake(ctx, lg.tickNum)
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
