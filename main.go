package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	const (
		width  = 16
		height = 16
	)
	w, err := newWin(width, height)
	if err != nil {
		log.Fatal(err)
	}
	s := newSnake(width/2, height/2, 3)
	rand.Seed(time.Now().UnixNano())
	food := newFood(width, height, s)
	score := 0
	d := up

	for {
		select {
		case d = <-w.input:
			if d == quit {
				goto Exit
			}
		default:
			if s.move(d, 1, 1, width-2, height-2) {
				goto Exit
			}
			if s.ate(food) {
				s.grow()
				food = newFood(width, height, s)
				score++
			}
			w.draw(s, food)
			time.Sleep(200*time.Millisecond - time.Duration(3*score)*time.Millisecond)
		}
	}
Exit:
	w.close()
	fmt.Printf("Score: %d\n", score)
}

func newFood(w, h int, s snake) point {
Start:
	x := rand.Intn(w-2) + 1
	y := rand.Intn(h-2) + 1
	for i := 0; i < len(s); i++ {
		if s[i].x == x && s[i].y == y {
			goto Start
		}
	}
	return point{x, y}
}
