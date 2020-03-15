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
	if now.IsAlive() && (aliveNum < 2 || aliveNum > 3) {
		return dead
	} else if aliveNum == 3 {
		return alive
	}
	return now
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

// NewEmptyCells return Cells.
// but, each cells don't have elements
func NewEmptyCells(width, height int) Cells {
	cells := make(Cells, height)
	// set Cell
	for i := 0; i < height; i++ {
		cells[i] = make([]Cell, width)

		for j := 0; j < width; j++ {
			cells[i][j] = Cell{[]<-chan State{}, []chan<- State{}, false}
		}
	}
	return cells
}
