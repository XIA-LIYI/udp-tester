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

var count int32 = 0
var totalByte uint64 = 0
var chans = [2]chan int{}

var bytes [2]uint64

var sendingByte int = 1250000000 / 25

func main() {
	for i := 0; i < 24; i++ {
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
			udpAddr, _ := net.ResolveUDPAddr("udp4", content + ":5050")
			socket, err := net.DialUDP("udp4", nil, udpAddr)
			if (err != nil) {
				fmt.Println("connection failed", err)
				return
			}
			go write(socket, chans[count])
			atomic.AddInt32(&count, 1)
			break
		}
	}
	elapsedTime := uint64(time.Since(startTime) / time.Millisecond / 1000)
	fmt.Println("Time consumed:", elapsedTime, "s")
	totalSpeed := totalByte / 1000 / elapsedTime * 8 / 1000
	fmt.Println("Speed is:", totalSpeed, "Mbps")
	result := ""
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

func write(socket *net.UDPConn,ch chan int) {
	<- ch
	ticker := time.NewTicker(time.Second / 1000)
	defer ticker.Stop()
	// conn.SetWriteBuffer(1000000)
	content := make([]byte, sendingByte)
	for {
		// <- ticker.C
		socket.Write(content)
	}
	
}

func listen() {
	fmt.Println("Listening")
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP: net.IPv4(0, 0, 0, 0),
		Port: 5050,
	})
	if err != nil {
		fmt.Printf("Listen failed, err:%v\n", err)
		return
	}
	for {
		var data [1024]byte
		n, addr, err := listen.ReadFromUDP(data[:])
		atomic.AddUint64(&totalByte, uint64(n))
		if err != nil {
			fmt.Printf("read failed, err:%v\n", err)
			continue
		}
		fmt.Printf("data:%s addr:%v count:%d\n", string(data[0:count]), addr, n)
	}
}


func onReceive(conn *net.TCPConn, index int32) {
	// fmt.Println("start receiving")
	// conn.SetReadBuffer(128000)
	buf := make([]byte, 156250)
	for {
		num, _ := conn.Read(buf)
		atomic.AddUint64(&totalByte, uint64(num))
		bytes[index] += uint64(num)
	}

}

func onSend(conn *net.TCPConn, ch chan int) {
	// fmt.Println("start sending")
	<- ch
	ticker := time.NewTicker(time.Second / 1000)
	defer ticker.Stop()
	// conn.SetWriteBuffer(1000000)
	content := make([]byte, sendingByte)

	// fmt.Println("start sending")
	for {
		// <- ticker.C
		conn.Write(content)
	}


}

