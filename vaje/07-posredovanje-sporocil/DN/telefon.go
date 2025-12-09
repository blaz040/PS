package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
	"strconv"
)

type message struct {
	data   []byte
	length int
}
var N int
var id int
var waitTime = time.Millisecond * 1000
var wg sync.WaitGroup

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
func receive(addr *net.UDPAddr) message{
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

	rMsg := message{}
	rMsg.data = buffer[:mLen]
	rMsg.length = mLen

	return rMsg
}

func send(addr *net.UDPAddr, msg message) {
	// Odpremo povezavo
	conn, err := net.DialUDP("udp", nil, addr)
	checkError(err)
	defer conn.Close()
	// Pripravimo sporočilo
	//sMsg := fmt.Sprint(id)
	// sMsg = string(msg.data[:msg.length]) + sMsg
	sMsg := string(msg.data[:msg.length])
	_, err = conn.Write([]byte(sMsg))
	checkError(err)
	fmt.Println("Telefon", id, "poslal sporočilo", sMsg, "telefonu na naslovu", addr)
}

func main() {
	// Preberi argumente
	portPtr := flag.Int("p", 9000, "# start port")
	idPtr := flag.Int("id", 0, "# process id")
	NPtr := flag.Int("n", 2, "total number of processes")
	rootPtr := flag.Int("rootId", 0, " # root id")

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

	// Izmenjava sporočil
	if id == root {
		receivedMsg := make([]chan bool, N)
		for i :=range receivedMsg {
			receivedMsg[i] = make(chan bool, 1)
		}
		go func() {
			for i := 0; i < N; i++ {
				msg := receive(rootAddr)
				//senderID := int(string(msg.data[:msg.length]))
				senderID, err := strconv.Atoi(string(msg.data[:msg.length]))            // convert string to int
    				if err != nil {
   					panic(err)
    				}
				receivedMsg[senderID] <- true
			}
		}()
		for i := 0; i < N; i++ {
			if i == root {
				continue
			}
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				sendPort := rootPort + i
				remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", sendPort))
				checkError(err)
				repeat := 10

				for repeat != 0 {
					repeat--
					send(remoteAddr, message{})
					select {
					case <-receivedMsg[i]:
						fmt.Printf("Received msg from %d \n", i)
						fmt.Printf("Stopped sending at %d repetitions\n",10-repeat)
						return
					default:
						time.Sleep(waitTime)
					}
				}
			}(i)
		}
		// send(remoteAddr, message{}) rMsg := receive(localAddr) fmt.Println(string(rMsg.data[:rMsg.length]) + "0")
		wg.Wait()
	} else {
		receive(localAddr)
		data := []byte(fmt.Sprintf("%d",id))
		send(rootAddr, message{data:data,length:len(data)})
	}
}
