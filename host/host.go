package host

import (
	"bufio"
	"fmt"
	ma "github.com/multiformats/go-multiaddr"
	"net"
	cr "p2p/crypto"
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
	Addrs() ma.Multiaddr

	// Network  returns the Network interface of the Host
	Network() net.Interface

	// Mux returns the Mux multiplexing incoming streams to protocol handlers
	//Mux() proto.Switch
}

type MyHost struct {
	peerID  peer.ID
	addrs   ma.Multiaddr
	network net.Interface
	//mux     proto.Switch
}

func (h *MyHost) ID() peer.ID {
	return h.peerID
}

func (h *MyHost) Addrs() ma.Multiaddr {
	return h.addrs
}

func (h *MyHost) Network() net.Interface {
	return h.network
}

//func (h *MyHost) Mux() proto.Switch {
//	return h.mux
//}

func getMyMultiaddr(inface string) (*net.Interface, ma.Multiaddr, error) {
	iface, err := net.InterfaceByName(inface)
	if err != nil {
		fmt.Printf("Could not find the interface because %s \n", err.Error())
		return nil, nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		fmt.Println(err)
		return iface, nil, err
	}

	for _, add := range addrs {
		fmt.Println(add.Network(), add.String())
	}
	var ipAddress net.IP
	var port string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipAddress = ipnet.IP
				port = "5200"
				break
			}
		}
	}

	addrStr := "/ip4/" + ipAddress.String() + "/tcp/" + port

	// Parsing the entered multi-address./ip4/127.0.0.1/tcp/8000
	addr, err := ma.NewMultiaddr(addrStr)
	if err != nil {
		fmt.Printf("Invalid multiaddress: %s, press r to retry or any other key to exit\n", err)
	}
	return iface, addr, nil
}

// ReceiveData reads data from the server and prints received messages.
func ReceiveData(conn net.Conn) []byte {
	reader := bufio.NewReader(conn)
	var returnLst []byte

	for {
		// Read a message from the server, terminated by a newline character ('\n').
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		// Print the received message.
		fmt.Print("Received message: ", msg)
		returnLst = append(returnLst, []byte(msg)...)
	}
	return returnLst
}

func GetHost() Host {
	_, pubKey, _ := cr.GenerateKeyPair(1, -1)
	id, _ := peer.GenerateIDFromPubKey(pubKey)
	fmt.Println("my ID:", []byte(id))
	network, addrs, _ := getMyMultiaddr("en0")

	host := &MyHost{
		peerID:  id,
		addrs:   addrs,
		network: *network,
	}
	return host
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
