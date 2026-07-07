package main

import(
	"fmt"
	"net"
)

func main(){
	addr := ":8080"
	udpAddr,err := net.ResolveUDPAddr("udp",addr)
	if(err !=nil){
		fmt.Println("error in resolvinf the udp adderss");
		panic(err)
	}
	conn,err := net.ListenUDP("udp",udpAddr)
	if(err!=nil){
		panic(err)
	}
	defer conn.Close()

	fmt.Println("hey bud nice work")
	buf := make( []byte,1024 )
	fmt.Println("UDP server listening on", addr)

	for{
		n,clientAddr,err := conn.ReadFromUDP(buf)
		if(err !=nil){
			fmt.Println("Error reading:", err)
			continue
		}
		fmt.Printf("Received from %s: %s\n", clientAddr, string(buf[:n]))
	}

}