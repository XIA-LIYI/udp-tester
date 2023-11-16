package main

import (
	// "bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
	// "sync/atomic"
	// "time"
)
const numOfMachines = 8
var connections [numOfMachines]*net.TCPConn
var ips [numOfMachines]string
var count int = 0
var allReady bool = false
var numOfNodesReady int32 = 0
var canClose chan int = make(chan int)
var startPort = 5600


func main() {
	go monitorInput()
	var tcpAddr *net.TCPAddr

	tcpAddr, _ = net.ResolveTCPAddr("tcp", "192.168.51.112:18787")

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if (err != nil) {
		fmt.Println(err)
	} 

	for {
		tcpConn, err := tcpListener.AcceptTCP()

		if err != nil {
			continue
		}
		connections[count] = tcpConn
		ips[count] = strings.Split(tcpConn.RemoteAddr().String(), ":")[0]
		// check()
		count += 1
		fmt.Println("A client connected:" + tcpConn.RemoteAddr().String())
		fmt.Println("Total number of connections:", count)
		if (count == numOfMachines) {
			tcpListener.Close()
			break
		}
	}

	broadcast()
	time.Sleep(time.Second)
	
	fmt.Println("check for check, start for start, stop for stop")
	start()
	fmt.Println("starting")
	time.Sleep(time.Second * 20)
	getResult()
	<- canClose
	// for {
	// 	var msg string
	// 	fmt.Scanln(&msg)
	// 	if msg == "quit" {
	// 		break
	// 	}
	// 	b := []byte(msg + "\n")
	// 	conn.Write(b)
	// }

}

func broadcast() {
	
	for _, key := range ips {
		curr := 0 
		for _, conn := range connections {
			conn.Write([]byte(key + ":" + strconv.Itoa(curr + startPort) + "\n") )
			curr += 1
		} 
	}
}

func monitorInput() {
	for {
		var msg string
		fmt.Scanln(&msg)
		if msg == "check" {
			check()
		}
		if msg == "start" {
			start()
		}
		if msg == "stop" {
			getResult()
			canClose <- 1
		}
	}

}

func getResult() {
	for i, conn := range connections {
		fmt.Printf(ips[i] + ": ")
		for {
			conn.Write([]byte("stop\n"))
			buf := make([]byte, 100)
			num, err := conn.Read(buf)
			if (err != nil) {
				continue
			}
			content := string(buf)[:num]
			fmt.Printf(content)
			fmt.Printf("\n")
			break
		}
	}
}

func listen(conn *net.TCPConn) {
	for {
		buf := make([]byte, 100)
		num, _ := conn.Read(buf)
		content := string(buf)[:num]
		fmt.Println(content)
	}
}

func check() {
	for i, conn := range connections {
		fmt.Printf("Checking ip:" + ips[i] + " ")
		conn.Write([]byte("check\n"))
		buf := make([]byte, 100)
		num, _ := conn.Read(buf)
		content := string(buf)[:num]
		fmt.Println(content)
		time.Sleep(time.Second / 2)
	}

}

func start() {
	for _, conn := range connections {
		conn.Write([]byte("start\n"))
	}

}