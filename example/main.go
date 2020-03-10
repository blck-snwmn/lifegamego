package main

import (
	"fmt"
	"time"

	"github.com/blck-snwmn/lifegame"
)

const csi = "\033["

func main() {
	// tickNum := 1
	height := 3
	width := 3
	cs, d := lifegame.NewLifeGame(height, width)
	cs.Cells[1][0].SetAlive()
	cs.Cells[1][1].SetAlive()
	cs.Cells[1][2].SetAlive()

	cs.Start()

	for i := 0; i < 10; i++ {
		for _, line := range d {
			for _, column := range line {
				s := <-column
				var v string
				if s.IsAlive() {
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
