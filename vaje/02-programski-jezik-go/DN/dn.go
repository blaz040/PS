package main

import "fmt"

type Student struct {
	ime     string
	priimek string
	ocene   []int
}

func dodajOceno(studenti map[string]Student, vpisnaStevilka string, ocena int) {
	student, ok := studenti[vpisnaStevilka]
	if !ok {
		fmt.Println("Študenta ni na seznamu :(")
		return
	}
	if 0 <= ocena && ocena <= 10 {
		student.ocene = append(student.ocene, ocena)
		studenti[vpisnaStevilka] = student

	} else {
		fmt.Println("Ocena ni v pravem rangu")
	}

}

func povprecje(studenti map[string]Student, vpisnaStevilka string) float64 {
	student, ok := studenti[vpisnaStevilka]
	if !ok {
		fmt.Println("Študent ne obstaja")
		return -1.0
	}
	var sum float64
	for _, v := range student.ocene {
		if v < 6 {
			return 0.0
		}
		sum += float64(v)
	}
	sum /= float64(len(student.ocene))

	return sum
}

func izpisRedovalnice(studenti map[string]Student) {
	fmt.Println("REDOVALNICA:")
	for v := range studenti {
		var student = studenti[v]
		fmt.Printf("%s - %s %s: %d\n", v, student.ime, student.priimek, student.ocene)
	}

}

func izpisiKoncniUspeh(studenti map[string]Student) {
	for v := range studenti {
		var student = studenti[v]
		var povp_ocena = povprecje(studenti, v)
		var uspeh = "Povprečen študent"

		switch {
		case povp_ocena >= 9:
			uspeh = "Odločen študent!"
		case povp_ocena < 6:
			uspeh = "Neuspešen študent"
		}
		fmt.Printf("%s %s: povprečna ocena %0.1f -> %s\n", student.ime, student.priimek, povp_ocena, uspeh)
	}
}
func main() {

	var studenti map[string]Student
	studenti = make(map[string]Student)

	studenti["1"] = Student{ime: "Jože", priimek: "Pločnik", ocene: []int{9, 9, 9}}
	studenti["2"] = Student{"Blaž", "Pirc", []int{10, 10, 10}}
	studenti["3"] = Student{"Maj", "Zorko", []int{8, 7, 6}}

	dodajOceno(studenti, "2", 6)
	dodajOceno(studenti, "1", 5)
	dodajOceno(studenti, "1", 11)

	fmt.Printf("%s povprečje: %0.2f\n", studenti["2"].ime, povprecje(studenti, "2"))
	fmt.Printf("%s povprečje: %0.2f\n", studenti["1"].ime, povprecje(studenti, "1"))
	fmt.Printf("Neko povprečje: %0.2f\n", povprecje(studenti, "4"))

	izpisRedovalnice(studenti)
	izpisiKoncniUspeh(studenti)
}
