package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/hashicorp/yamux"
	ma "github.com/multiformats/go-multiaddr"
	"io"
	"net"
	"os"
	cr "p2p/crypto"
	"p2p/host"
	_ "p2p/host"
	_ "p2p/network"
	_ "p2p/routing"
	_ "p2p/transport"
)

// GetAddrFromUser takes in user input to return a multi address or an error
func GetAddrFromUser() (ma.Multiaddr, error) {
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

// PrintProtocols takes in a multi address and returns the list of protocols in the address
func PrintProtocols(addr ma.Multiaddr) {
	fmt.Printf("Protocols in Multiaddr %s are:\n", addr.String())
	for _, protocol := range addr.Protocols() {
		fmt.Println("PName", protocol.Name)
	}
}

// newStream takes in a yamux session and returns a multiplexed connection
func newStreamYamux(session *yamux.Session) (net.Conn, error) {
	stream, err := session.Open()
	if err != nil {
		fmt.Printf("Error opening a new Stream:%s \n", err.Error())
	}
	return stream, err
}

// connectTcpIp4 returns tcp connection from multiaddr
func connectTcpIp4(addr ma.Multiaddr) (*net.TCPConn, error) {

	ipAddr, err := addr.ValueForProtocol(ma.P_IP4)
	if err != nil {
		fmt.Println("Could not get ipv4 address from the given multiadress")
		return nil, err
	}

	tcpPort, err := addr.ValueForProtocol(ma.P_TCP)
	if err != nil {
		fmt.Println("Could not get TCP Port from the given multiadress")
		return nil, err
	}

	peerAddr := ipAddr + ":" + tcpPort
	fmt.Printf("Resolved tcp addr from multiaddr: %s \n", peerAddr)

	//Resolving the TCP address
	ipAdd, err := net.ResolveTCPAddr("tcp", peerAddr)
	if err != nil {
		fmt.Printf("Error resolving TCP address(%s): %s \n", peerAddr, err)
		return nil, err
	}

	// Get a TCP connection
	conn, err := net.DialTCP("tcp", nil, ipAdd)
	if err != nil {
		fmt.Printf("Could not establish with peer %s because %s", peerAddr, err.Error())
		return nil, err
	}

	return conn, nil
}

// setupSenderYamux returns yamux session after taking in connection and config
// set config to nil for default configuration
func setupSenderYamux(conn io.ReadWriteCloser, config *yamux.Config) (*yamux.Session, error) {
	session, err := yamux.Client(conn, config)
	if err != nil {
		fmt.Printf("Could not wrap yamux on connection because %s \n", err.Error())
		return nil, err
	}
	return session, nil
}

func main() {

	addr, err := GetAddrFromUser()

	PrintProtocols(addr)

	conn, err := connectTcpIp4(addr)
	//Ensure the connection is closed when the function returns
	defer func(conn *net.TCPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("error: ", err.Error())
		}
	}(conn)

	fmt.Println("Connected to peer:", addr.String())

	session, err := setupSenderYamux(conn, nil)

	newConn, err := newStreamYamux(session)

	lent, err := newConn.Write([]byte("test"))

	if err != nil {
		fmt.Printf("Could not write on connection %s because of %s\n", newConn, err.Error())
	}
	fmt.Printf("successfully wrote %d bytes \n", lent)

	go host.ReceiveData(conn)

	privKey, pubKey, _ := cr.GenerateKeyPair(1, -1)
	fmt.Println("priv>", privKey)
	fmt.Println("pub>", pubKey)

	/*
		privKey, pubKey, _ := cr.GenerateKeyPair(1, -1)
		fmt.Println("priv>", privKey)
		fmt.Println("pub>", pubKey)

		id, _ := peer.GenerateIDFromPubKey(pubKey)

		fmt.Println(len(id))
		DefaultMuxers := Muxer(YamuxID, yamux.DefaultConfig())*/

}
