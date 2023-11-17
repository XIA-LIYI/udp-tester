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

const numOfMachines = 12
const numOfThreads = 1
var count int32 = 0
var totalByte uint64 = 0
var sendBytes uint64 = 0
var bufferSize int = 1500

var chans = [numOfMachines * numOfThreads]chan int{}

var bytes [numOfMachines]uint64

var ips = make(map[string]int, numOfMachines * numOfThreads)

var sendPacks [numOfMachines]uint64
var receivePacks [numOfMachines]int

func main() {
	for i := 0; i < numOfMachines * numOfThreads; i++ {
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

	for i := 0; i < numOfMachines; i++ {
		go listen(5600 + i, i)
	}
	

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
			udpAddr, _ := net.ResolveUDPAddr("udp", content)
			for i := 0; i < numOfThreads; i++ {
				socket, _ := net.DialUDP("udp", nil, udpAddr)
				go write(socket, chans[count + numOfMachines * int32(i)], int(count))
			}
			// socket, _ := net.DialUDP("udp", nil, udpAddr)
			// go write(socket, chans[count])
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
	for _, i := range sendPacks {
		result = result + strconv.Itoa(int(i)) + " "
	}
	for _, i := range receivePacks {
		result = result + strconv.Itoa(int(i)) + " "
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

func write(socket *net.UDPConn, ch chan int, index int) {
	<- ch
	speed := int(1.2 * 1000 * 100 * 1000 / 1000 / numOfMachines / bufferSize / numOfThreads)
	ticker := time.NewTicker(time.Second / 1000)
	defer ticker.Stop()
	socket.SetWriteBuffer(bufferSize * 10)
	content := make([]byte, bufferSize)
	for {
		<- ticker.C
		for i := 0; i < speed; i++ {
			socket.Write(content)
			atomic.AddUint64(&sendBytes, uint64(bufferSize))
			atomic.AddUint64(&sendPacks[index], 1)
		}
		// socket.Write(content)
		// atomic.AddUint64(&sendBytes, uint64(bufferSize))
	}
	
}

func addrToIndex(addr string) int {
	ip := strings.Split(addr, ":")[0]
	index := ips[ip]
	return index
}

func listen(port int, index int) {
	fmt.Println("Listening")
	udpAddr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:" + strconv.Itoa(port))
	listen, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Printf("Listen failed, err:%v\n", err)
		return
	}
	listen.SetReadBuffer(bufferSize * 64)
	data := make([]byte, bufferSize + 10)
	for {
		n, _, _ := listen.ReadFromUDP(data)
		atomic.AddUint64(&totalByte, uint64(n))
		receivePacks[index] += 1
		// bytes[addrToIndex(addr.String())] += uint64(n)		
		// if err != nil {
		// 	fmt.Printf("read failed, err:%v\n", err)
		// 	continue
		// }
	}
}



