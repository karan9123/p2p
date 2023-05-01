package main

import (
	"fmt"
	ma "github.com/multiformats/go-multiaddr"
	"p2p/host"
)

// PrintProtocols takes in a multi address and returns the list of protocols in the address
func PrintProtocols(addr ma.Multiaddr) {
	fmt.Printf("Protocols in Multiaddr %s are:\n", addr.String())
	for _, protocol := range addr.Protocols() {
		fmt.Println("PName", protocol.Name)
	}
}

func main() {

	//fmt.Printf("List the protocols you support. Don't lie and don't be shy, go ahead list em all")
	//protos, err := proto.GetProtocols("protocols.txt")

	myHost := host.GetHost("5002")
	fmt.Println(myHost.ID(), myHost.Addrs(), myHost.Network().HardwareAddr)
	PrintProtocols(myHost.Addrs())

}

func receiverMethod(myHost host.Host) {
	conn, err := myHost.StartListening()
	if err != nil {
		fmt.Println("Couldn't listen due to ", err.Error())
	}
	buf := make([]byte, 1024)
	i, err := conn.Read(buf)
	fmt.Printf("read %d bytes which are %s\n", i, buf[:i])
}

func senderMethod(myHost host.Host) {
	sendingConn, _ := myHost.SenderConn()
	i, err := sendingConn.Write([]byte("hello"))
	if err != nil {
		return
	}
	fmt.Printf("Wrote %d bytes\n", i)
}
