package main

import (
	"fmt"
	"sync"
)

type narocilo interface {
	obdelaj()
}

type izdelek struct {
	imeIzdelka string
	cena       float64
	teza       float64
}

type eknjiga struct {
	naslovKnjige string
	cena         float64
}

type spletniTecaj struct {
	imeTecaja   string
	trajanjeUre int
	cenaUre     float64
}

var lock sync.Mutex
var promet float64
var stNarocil int
var wg sync.WaitGroup

func (iz izdelek) obdelaj() {
	defer wg.Done()
	lock.Lock()

	stNarocil++
	fmt.Println("\nŠtevilka naročila:", stNarocil)

	fmt.Printf(
		"Ime izdelka: %s \n"+
			"Cena: %g €\n"+
			"Teža: %g kg\n", iz.imeIzdelka, iz.cena, iz.teza)

	promet += iz.cena

	lock.Unlock()
}
func (ek eknjiga) obdelaj() {
	defer wg.Done()

	lock.Lock()

	stNarocil++
	fmt.Println("\nŠtevilka naročila:", stNarocil)

	fmt.Printf(
		"Naslov knjige: %s\n"+
			"Cena: %g\n", ek.naslovKnjige, ek.cena)

	promet += ek.cena

	lock.Unlock()
}
func (st spletniTecaj) obdelaj() {
	defer wg.Done()

	lock.Lock()

	stNarocil++
	fmt.Println("\nŠtevilka naročila:", stNarocil)

	fmt.Printf(
		"Ime tecaja: %s\n"+
			"Trajanje : %d h\n"+
			"Cena na uro: %g €\n", st.imeTecaja, st.trajanjeUre, st.cenaUre)

	promet += float64(st.trajanjeUre) * st.cenaUre

	lock.Unlock()
}
func main() {
	izdelek1 := izdelek{"Prenosnik", 2000, 2.5}
	eknjiga1 := eknjiga{"Bajke", 5.2}
	spletniTecaj1 := spletniTecaj{"GO", 67, 15.2}

	narocila := []narocilo{izdelek1, eknjiga1, spletniTecaj1}

	for _, n := range narocila {
		wg.Add(1)
		go n.obdelaj()
	}
	wg.Wait()
	fmt.Printf("\nPromet: %g\nStevilo narocil: %d\n", promet, stNarocil)
}
