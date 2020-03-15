package lifegame

import (
	"context"
	"testing"
)

func TestCell_tick(t *testing.T) {
	size := 4
	fromChs := make([]<-chan State, size)
	sendChs := make([]chan State, size)

	toCh := make(chan State)

	for i := 0; i < size; i++ {
		ch := make(chan State)
		fromChs[i] = ch
		sendChs[i] = ch
	}
	defer func() {
		close(toCh)
		for _, ch := range sendChs {
			close(ch)
		}
	}()

	c := Cell{from: fromChs, to: []chan<- State{toCh}, state: alive}
	ctx := context.Background()

	go func() {
		for _, ch := range sendChs {
			ch <- alive
		}
	}()
	go c.tick(ctx)
	if s := <-toCh; dead != s {
		t.Errorf("received state mismatch. want=%v, got=%v", dead, s)
	}

	go func() {
		aliveNum := 3
		for _, ch := range sendChs[0:aliveNum] {
			ch <- alive
		}
		for _, ch := range sendChs[aliveNum:] {
			ch <- dead
		}
	}()
	go c.tick(ctx)
	if s := <-toCh; alive != s {
		t.Errorf("received state mismatch. want=%v, got=%v", alive, s)
	}

	go func() {
		aliveNum := 3
		for _, ch := range sendChs[0:aliveNum] {
			ch <- alive
		}
		for _, ch := range sendChs[aliveNum:] {
			ch <- dead
		}
	}()
	go c.tick(ctx)
	if s := <-toCh; alive != s {
		t.Errorf("received state mismatch. want=%v, got=%v", alive, s)
	}
}
