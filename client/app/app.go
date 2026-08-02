package app

import(
	"fmt"
	"net"
)

type App struct{
	senderIp string
	senderPort int
	receiverIp string
	receiverPort int
	filePathToSend string
}

func (a *App) DiscoverReceiver() {
	fmt.Println("Discovering receiver...")
	//to implement discovery via heartbeat
	addr := net.UDPAddr{
		Port: a.senderPort,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error starting UDP client:", err)
		return
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, senderAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}
		message := string(buffer[:n])
		if message == "heartbeat" {
			fmt.Println("Received heartbeat from receiver at", senderAddr)
			//send a response back to the receiver to acknowledge the heartbeat
			receiverAddr := net.UDPAddr{
				Port: a.receiverPort,
				IP:   net.ParseIP(a.receiverIp),
			}
			a.receiverIp = senderAddr.IP.String()
			_, err = conn.WriteToUDP([]byte("ack"), &receiverAddr)
			if err != nil {
				fmt.Println("Error sending acknowledgment:", err)
				return
			}
			return
		}
	}
}