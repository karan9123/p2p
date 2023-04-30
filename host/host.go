package host

import (
	"bufio"
	"errors"
	"fmt"
	ma "github.com/multiformats/go-multiaddr"
	"net"
	"os"
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

	// Peerstore returns the Host's repository of Peer Addresses and Keys.
	//Peerstore() peerstore.Peerstore

	// Returns the listen addresses of the Host
	Addrs() []ma.Multiaddr

	// Networks returns the Network interface of the Host
	Network() net.Interface

	// Mux returns the Mux multiplexing incoming streams to protocol handlers
	//Mux() protocol.Switch

	// Connect ensures there is a connection between this host and the peer with
	// given peer.ID. Connect will absorb the addresses in pi into its internal
	// peerstore. If there is not an active connection, Connect will issue a
	// h.Network.Dial, and block until a connection is open, or an error is
	// returned. // TODO: Relay + NAT.
	//Connect(ctx context.Context, pi peer.AddrInfo) error

	// SetStreamHandler sets the protocol handler on the Host's Mux.
	// This is equivalent to:
	//   host.Mux().SetHandler(proto, handler)
	// (Threadsafe)
	//SetStreamHandler(pid protocol.ID, handler network.StreamHandler)

	// SetStreamHandlerMatch sets the protocol handler on the Host's Mux
	// using a matching function for protocol selection.
	//SetStreamHandlerMatch(protocol.ID, func(protocol.ID) bool, network.StreamHandler)

	// RemoveStreamHandler removes a handler on the mux that was set by
	// SetStreamHandler
	//RemoveStreamHandler(pid protocol.ID)

	// NewStream opens a new stream to given peer p, and writes a p2p/protocol
	// header with given ProtocolID. If there is no connection to p, attempts
	// to create one. If ProtocolID is "", writes no header.
	// (Threadsafe)
	//NewStream(ctx context.Context, p peer.ID, pids ...protocol.ID) (network.Stream, error)

	// Close shuts down the host, its Network, and services.
	//Close() error

	// ConnManager returns this hosts connection manager
	//ConnManager() connmgr.ConnManager

	// EventBus returns the hosts eventbus
	//EventBus() event.Bus
}

// getAddrFromUser takes in user input to return a multi address or an error
func getAddrFromUser() (ma.Multiaddr, error) {
	for {
		fmt.Print("Enter the server multi-address (e.g., /ip4/127.0.0.1/tcp/8000): ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		addrStr := scanner.Text()

		// Parsing the entered multi-address.
		addr, err := ma.NewMultiaddr(addrStr)
		if err != nil {
			fmt.Printf("Invalid multiaddress: %s, press r to retry or any other key to exit\n", err)
			scanner.Scan()
			input := scanner.Text()
			if input == "r" {
				continue
			} else {

				return nil, errors.New("exit from get Address")
			}
		} else {
			return addr, nil
		}
	}
}

// printProtocols takes in a multi address and returns the list of protocols in the address
func printProtocols(addr ma.Multiaddr) {
	fmt.Printf("Protocols in Multiaddr %s are:\n", addr.String())
	for _, protocol := range addr.Protocols() {
		fmt.Println("PName", protocol.Name)
	}
}

// receiveData reads data from the server and prints received messages.
func receiveData(conn net.Conn) {
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

func main() {

	addr, err := getAddrFromUser()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	printProtocols(addr)
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

}
