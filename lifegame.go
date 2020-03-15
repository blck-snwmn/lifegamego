package lifegame

import "context"

// New return LifeGame and channels to show
func New(width, height, tickNum int) (*LifeGame, [][]<-chan State) {
	lg := &LifeGame{}
	lg.tickNum = tickNum
	lg.Cells = NewEmptyCells(width, height)
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

			existsTop := top >= 0
			existsBtm := btm < height
			existsRight := right < width
			existsLeft := left >= 0

			// top
			lg.genCell(i, j, top, j, existsTop)
			// bottom
			lg.genCell(i, j, btm, j, existsBtm)
			// left
			lg.genCell(i, j, i, left, existsLeft)
			// right
			lg.genCell(i, j, i, right, existsRight)
			// top-left
			lg.genCell(i, j, top, left, existsTop && existsLeft)
			// top-right
			lg.genCell(i, j, top, right, existsTop && existsRight)
			// bottom-left
			lg.genCell(i, j, btm, left, existsBtm && existsLeft)
			// bottom-right
			lg.genCell(i, j, btm, right, existsBtm && existsRight)
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
