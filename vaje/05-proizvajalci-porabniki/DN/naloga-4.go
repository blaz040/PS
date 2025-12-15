package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
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
var wgGeneratorji, wgObdelave sync.WaitGroup
var moznaNarocila []narocilo

func (iz izdelek) obdelaj() {
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

func generatorNarocil(interval time.Duration, narocila chan narocilo, quit chan struct{}) {
	defer wgGeneratorji.Done()
	for {
		select {
		case <-quit:
			return
		default:
			idx := rand.Intn(len(moznaNarocila))
			narocila <- moznaNarocila[idx]
			time.Sleep(interval)
		}
	}
}
func obdelavaNarocil(interval time.Duration, narocila chan narocilo) {
	defer wgObdelave.Done()
	for {
		if narocilo, ok := <-narocila; ok {
			narocilo.obdelaj()
		} else {
			return
		}
		time.Sleep(interval)
	}
}
func listenForButton(quit chan struct{}) {
	fmt.Scanln()
	close(quit)
}

func main() {
	izdelek1 := izdelek{"Prenosnik", 2000, 2.5}
	eknjiga1 := eknjiga{"Bajke", 5.2}
	spletniTecaj1 := spletniTecaj{"GO", 67, 15.2}

	moznaNarocila = []narocilo{izdelek1, eknjiga1, spletniTecaj1}

	narocila := make(chan narocilo)
	quit := make(chan struct{})

	nProizvodenj := 5
	nPorabnikov := 2

	intervalProizvodnje := time.Millisecond * 50
	intervalObdelave := time.Millisecond * 50

	wgGeneratorji.Add(nProizvodenj)
	for i := 0; i < nProizvodenj; i++ {
		go generatorNarocil(intervalProizvodnje, narocila, quit)
	}
	wgObdelave.Add(nPorabnikov)
	for i := 0; i < nPorabnikov; i++ {
		go obdelavaNarocil(intervalObdelave, narocila)
	}
	go listenForButton(quit)

	wgGeneratorji.Wait()
	close(narocila)
	wgGeneratorji.Wait()
	fmt.Printf("\nPromet: %g €\nStevilo narocil: %d\n", promet, stNarocil)
}
