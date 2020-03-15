package lifegame

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestState_String(t *testing.T) {
	t.Run("state.alive return `Alive`", func(t *testing.T) {
		if fmt.Sprintf("%s", alive) != "Alive" {
			t.Errorf("State.String() = %s, want %s", alive, "Alive")
		}
	})
	t.Run("state.dead return `Dead`", func(t *testing.T) {
		if fmt.Sprintf("%s", dead) != "Dead" {
			t.Errorf("State.String() = %s, want %s", dead, "Dead")
		}
	})
}

func Test_changeState(t *testing.T) {
	type args struct {
		now      State
		aliveNum int
	}
	tests := []struct {
		name string
		args args
		want State
	}{
		{"now dead, near 1 alive cell -> alive", args{dead, 1}, dead},
		{"now dead, near 2 alive cell -> dead", args{dead, 2}, dead},
		{"now dead, near 3 alive cell -> dead", args{dead, 3}, alive},
		{"now dead, near 4 alive cell -> alive", args{dead, 4}, dead},
		{"now dead, near 5 alive cell -> dead", args{dead, 5}, dead},
		{"now dead, near 6 alive cell -> dead", args{dead, 6}, dead},
		{"now dead, near 7 alive cell -> alive", args{dead, 7}, dead},
		{"now dead, near 8 alive cell -> dead", args{dead, 8}, dead},
		{"now alive, near 1 alive cell -> dead", args{alive, 1}, dead},
		{"now alive, near 2 alive cell -> alive", args{alive, 2}, alive},
		{"now alive, near 3 alive cell -> alive", args{alive, 3}, alive},
		{"now alive, near 4 alive cell -> alive", args{alive, 4}, dead},
		{"now alive, near 5 alive cell -> dead", args{alive, 5}, dead},
		{"now alive, near 6 alive cell -> alive", args{alive, 6}, dead},
		{"now alive, near 7 alive cell -> alive", args{alive, 7}, dead},
		{"now alive, near 8 alive cell -> alive", args{alive, 8}, dead},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := changeState(tt.args.now, tt.args.aliveNum); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("changeState() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestCell_sendState(t *testing.T) {
	ch := make(chan State)
	defer close(ch)
	c := Cell{to: []chan<- State{ch}}

	ctx := context.Background()

	c.state = alive
	go func() {
		c.sendState(ctx)
	}()
	if s := <-ch; alive != s {
		t.Errorf("received state mismatch. want=%v, got=%v", alive, s)
	}

	c.state = dead
	go func() {
		c.sendState(ctx)
	}()
	if s := <-ch; dead != s {
		t.Errorf("received state mismatch. want=%v, got=%v", dead, s)
	}
}

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
