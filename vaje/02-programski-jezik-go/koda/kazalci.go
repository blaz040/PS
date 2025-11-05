/*
Program kazalci prikazuje uporabo kazalcev v jeziku go
*/
package main

import "fmt"

type point struct {
	x int
	y int
}

func print(p *int) {
	fmt.Println(*p)
}
func main() {
	i, j := 42, 1337
	// Ustvarimo kazalec p, ki kaže na spremenljivko i
	p := &i
	// Dostop do vrednosti, na katero kaže kazalec
	defer print(p)
	// Spremenimo vrednost preko kazalca
	*p = j
	fmt.Println(i)
	i = 55
	// Pri kazalcih na strukture lahko uporabimo poenostavljeno sintakso brez *, ko spreminjamo vrednosti preko kazalca
	t := point{}
	p2 := &t
	p2.x = 15

	fmt.Println(t)

}
