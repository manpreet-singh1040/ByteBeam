package app

import (
	"fmt"
	"net"
	"time"
)

type App struct{
	SenderIp string
	SenderPort int
	RecieverIp string
	RecieverPort int
	StartReceiver bool
}

// func (a *App) DiscoverSender() {
// 	fmt.Println("Discovering sender...")
// 	//to implement discovery via heartbeat
// }



func (a *App) StartDiscoveryOfSender() {
	stop := make(chan struct{})
	go a.Heartbeat(stop)

	addr := net.UDPAddr{
		Port: a.RecieverPort,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error starting UDP server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("UDP server listening on", addr.String())
	buffer := make([]byte, 1024)
	for {
		n, senderAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}
		message := string(buffer[:n])
		if message == "ack" {
			fmt.Println("Received ack from sender")
			a.SenderIp = senderAddr.IP.String()
			a.StartReceiver = true
			close(stop)
			return
		}
	}
}

func (a *App) Heartbeat(stop <- chan struct{}){
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:{
				fmt.Println("Sending heartbeat...")
				addr := net.UDPAddr{
					Port: a.SenderPort,
					IP:   net.ParseIP(a.SenderIp),
				}
				conn, err := net.DialUDP("udp", nil, &addr)
				if err != nil {
					fmt.Println("Error dialing UDP:", err)
					return
				}
				conn.Write([]byte("heartbeat"))
				conn.Close()
			}
		case <- stop:
			fmt.Println("Stopping heartbeat...")
			return
		}
	}
}
