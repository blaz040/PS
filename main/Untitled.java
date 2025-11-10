package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

type Meritev struct {
	vrsta    string
	vrednost float32
}

func round(vred float32, precision int) float32 {
	v := float64(vred)
	c := math.Pow10(precision)
	return float32(math.Round(v*c) / c)
}
func sensor(vrsta string, ran float32, init float32, meritve chan Meritev) {
	c := init
	for {
		coef := rand.Float32()*2 - 1
		c = c + ran*coef
		var vrednost float32 = round(c, 2)
		meritve <- Meritev{vrsta, vrednost}

		time.Sleep(1000 * time.Millisecond)
	}
}

func readKey(input chan bool) {
	fmt.Scanln()
	input <- true
}

func main() {

	meritve := make(chan Meritev, 3)
	input := make(chan bool)

	go sensor("Temperatura", 3, 15, meritve)
	go sensor("Vlaga", 5, 50, meritve)
	go sensor("Tlak", 3, 900, meritve)

	go readKey(input)

	for {
		select {
		case msg := <-meritve:
			fmt.Println(msg.vrsta, ":", msg.vrednost)
		case <-input:
			close(meritve)
		}
	}
}
