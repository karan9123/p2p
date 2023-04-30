package host

import (
	"bufio"
	"fmt"
	"github.com/libp2p/go-libp2p/core/protocol"
	ma "github.com/multiformats/go-multiaddr"
	"net"
	"p2p/peer"
)

// Host represents a single libp2p node in a peer-to-peer network.

// Host is an object participating in a p2p network, which
// implements protocols or provides services. It handles
// requests like a Server, and issues requests like a Client.
// It is called Host because it is both Server and Client (and Peer
// may be confusing).
type Host interface {
	// ID returns the (local) peer.ID associated with this Host
	ID() peer.ID

	// Addrs Returns the listen addresses of the Host
	Addrs() []ma.Multiaddr

	// Network  returns the Network interface of the Host
	Network() net.Interface

	// Mux returns the Mux multiplexing incoming streams to protocol handlers
	Mux() protocol.Switch
}

// ReceiveData reads data from the server and prints received messages.
func ReceiveData(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		// Read a message from the server, terminated by a newline character ('\n').
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		// Print the received message.
		fmt.Print("Received message: ", msg)
	}
}

/*func main() {

	addr, err := GetAddrFromUser()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	PrintProtocols(addr)
	ipAddr, err := addr.ValueForProtocol(ma.P_IP4)
	if err != nil {
		fmt.Println("Could not get ipv4 address from the given multiadress")
	}
	tcpPort, err := addr.ValueForProtocol(ma.P_TCP)
	if err != nil {
		fmt.Println("Could not get TCP Port from the given multiadress")
	}
	peerAddr := ipAddr + ":" + tcpPort
	fmt.Println(peerAddr)

	//Resolving the TCP address
	ipAdd, err := net.ResolveTCPAddr("tcp", peerAddr)
	if err != nil {
		fmt.Printf("Error resolving TCP address: %s\n", err)
		return
	}
	fmt.Println(ipAdd)

	// Establish a TCP connection to the server.
	conn, err := net.DialTCP("tcp", nil, ipAdd)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer func(conn *net.TCPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("error: ", err.Error())
		}
	}(conn) // Ensure the connection is closed when the function returns.

	fmt.Println("Connected to peer:", addr.String())

	go receiveData(conn)

}*/
