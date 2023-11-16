package main

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"
	"strconv"
	"bufio"
	"strings"

)

const numOfMachines = 2
var count int32 = 0
var totalByte uint64 = 0
var sendBytes uint64 = 0
var bufferSize int = 1024

var chans = [numOfMachines]chan int{}

var bytes [numOfMachines]uint64

var ips = make(map[string]int, numOfMachines)

func main() {
	for i := 0; i < numOfMachines; i++ {
		chans[i] = make(chan int)
	}
	
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "192.168.51.112:18787")
	var conn *net.TCPConn
	var err error
	for {
		conn, err = net.DialTCP("tcp", nil, tcpAddr)
		if (err != nil) {
			fmt.Println(err)
			continue
		} else {
			break
		} 
	}

	defer conn.Close()
	fmt.Println("connected!")

	go listen()

	startTime := time.Now()

	reader := bufio.NewReader(conn)
	for {
		data, _ := reader.ReadString('\n')
		content := strings.Replace(string(data), "\n", "", -1)  
		fmt.Println(content)
		if (content == "check") {
			conn.Write([]byte(strconv.Itoa(int(count))))
			continue
		}
		if (content == "start") {
			startTime = time.Now()
			fmt.Println("Current number of connections is:", count)
			for i := range chans {
				chans[i] <- 0
			}
			fmt.Println("All are released!")
			continue
		}
		if (content == "stop") {
			break
		}
		for {
			udpAddr, _ := net.ResolveUDPAddr("udp", content + ":5050")
			socket, err := net.DialUDP("udp", nil, udpAddr)
			if (err != nil) {
				fmt.Println("connection failed", err)
				return
			}
			ips[content] = int(count) 
			go write(socket, chans[count])
			atomic.AddInt32(&count, 1)
			break
		}
	}
	elapsedTime := uint64(time.Since(startTime) / time.Millisecond / 1000)
	fmt.Println("Time consumed:", elapsedTime, "s")
	totalSpeed := totalByte / 1000 / elapsedTime * 8 / 1000
	fmt.Println("Receive speed is:", totalSpeed, "Mbps")
	sendSpeed := sendBytes / 1000 / elapsedTime * 8 / 1000
	fmt.Println("Send speed is:", sendSpeed, "Mbps")
	result := strconv.Itoa(int(sendSpeed)) + " " + strconv.Itoa(int(totalSpeed)) + " "
	for _, i := range bytes {
		speed := i / 1000 / elapsedTime * 8 / 1000
		fmt.Println("Single speed is:", speed, "Mbps")
		result = result + strconv.Itoa(int(speed)) + " "
	}
	conn.Write([]byte(result))

	// 控制台聊天功能加入
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

func write(socket *net.UDPConn, ch chan int) {
	<- ch
	speed := int(1.25 * 1000 * 1000 * 1000 / 10000 / numOfMachines / bufferSize) + 1
	ticker := time.NewTicker(time.Second / 10000)
	defer ticker.Stop()
	socket.SetWriteBuffer(bufferSize * 10)
	content := make([]byte, bufferSize)
	for {
		<- ticker.C
		for i := 0; i < speed; i++ {
			socket.Write(content)
			atomic.AddUint64(&sendBytes, uint64(bufferSize))
		}
		// socket.Write(content)
		// sendBytes += uint64(bufferSize)
	}
	
}

func addrToIndex(addr string) int {
	ip := strings.Split(addr, ":")[0]
	index := ips[ip]
	return index
}

func listen() {
	fmt.Println("Listening")
	udpAddr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:5050")
	listen, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Printf("Listen failed, err:%v\n", err)
		return
	}
	listen.SetReadBuffer(1024 * 32)
	data := make([]byte, bufferSize + 1)
	for {
		n, addr, _ := listen.ReadFromUDP(data)
		atomic.AddUint64(&totalByte, uint64(n))
		bytes[addrToIndex(addr.String())] += uint64(n)		
		// if err != nil {
		// 	fmt.Printf("read failed, err:%v\n", err)
		// 	continue
		// }
	}
}



