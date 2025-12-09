package dn

import (
	"flag"
	"fmt"
	"net"
	"time"
)

type message struct {
	data   []byte
	length int
}

var receivedMsg chan bool
var N int
var id int
var waitTime = time.Millisecond * 500

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
func receive(addr *net.UDPAddr) {
	// Poslušamo
	conn, err := net.ListenUDP("udp", addr)
	checkError(err)
	defer conn.Close()
	fmt.Println("Telefon", id, "posluša na", addr)
	buffer := make([]byte, 1024)
	// Preberemo sporočilo
	mLen, err := conn.Read(buffer)
	checkError(err)
	fmt.Println("Telefon", id, "prejel sporočilo:", string(buffer[:mLen]))

	receivedMsg <- true
}

func send(addr *net.UDPAddr, msg message) {
	// Odpremo povezavo
	conn, err := net.DialUDP("udp", nil, addr)
	checkError(err)
	defer conn.Close()
	// Pripravimo sporočilo
	sMsg := fmt.Sprint(id) + "-"
	sMsg = string(msg.data[:msg.length]) + sMsg
	_, err = conn.Write([]byte(sMsg))
	checkError(err)
	fmt.Println("Telefon", id, "poslal sporočilo", sMsg, "telefonu na naslovu", addr)
}

func main() {
	// Preberi argumente
	portPtr := flag.Int("p", 9000, "# start port")
	idPtr := flag.Int("id", 0, "# process id")
	NPtr := flag.Int("n", 2, "total number of processes")
	rootPtr := flag.Int("n", 0, " # root id")

	flag.Parse()

	rootPort := *portPtr
	id = *idPtr
	N = *NPtr
	root := *rootPtr
	basePort := rootPort + id

	// Ustvari potrebne mrežne naslove
	rootAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", rootPort))
	checkError(err)

	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", basePort))
	checkError(err)

	// Ustvari kanal, ki bo signaliziral, da so vsi procesi pripravljeni
	receivedMsg = make(chan bool)

	// Izmenjava sporočil
	if id == root {
		for i := 0; i <= N; i++ {
			if i == root {
				continue
			}
			go func() {
				sendPort := rootPort + i

				remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", sendPort))
				checkError(err)
				repeat := 5

				go receive(rootAddr)

				for repeat != 0 {
					repeat--
					send(remoteAddr, message{})
					select {
					case <-receivedMsg:
						fmt.Printf("Received msg from %d", i)
						return
					default:
						time.Sleep(waitTime)
					}
				}
			}()
		}

		// send(remoteAddr, message{})
		//rMsg := receive(localAddr)
		//fmt.Println(string(rMsg.data[:rMsg.length]) + "0")
	} else {
		receive(localAddr)
		send(rootAddr, message{})
	}
}
