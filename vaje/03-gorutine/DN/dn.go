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
func sensor(vrsta string, ran float32, init float32, meritve chan Meritev, stop chan struct{}) {
	c := init
	for {
		select {
		case <-stop:
			return
		default:
			coef := rand.Float32()*2 - 1
			c = c + ran*coef
			var vrednost float32 = round(c, 2)
			meritve <- Meritev{vrsta, vrednost}

			time.Sleep(100 * time.Millisecond)
		}
	}
}
func readKey(input chan bool) {
	fmt.Scanln()
	input <- true
}
func main() {

	meritve := make(chan Meritev, 3)
	input := make(chan bool)

	stop := make(chan struct{})

	go sensor("Temperatura", 3, 15, meritve, stop)
	go sensor("Vlaga", 5, 50, meritve, stop)
	go sensor("Tlak", 3, 900, meritve, stop)

	go readKey(input)

	for {
		select {
		case msg := <-meritve:
			fmt.Println(msg.vrsta, ":", msg.vrednost)
		case <-input:
			close(stop)
		case <-time.After(5 * time.Second):
			fmt.Println("Sistem neodziven za 5s. \nKončujem....")
			return
		}

	}
}
