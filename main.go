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
	"p2p/host"
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

// connectTcpIp4 returns tcp connection from multiaddr
func connectTcpIp4(addr ma.Multiaddr) (*net.TCPConn, error) {

	ipAddr, tcpPort, err := host.GetIp4TcpFromMultiaddr(addr)
	if err != nil {
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

// newStream takes in a yamux session and returns a multiplexed connection
func newStreamYamux(session *yamux.Session) (net.Conn, error) {
	stream, err := session.Open()
	if err != nil {
		fmt.Printf("Error opening a new Stream:%s \n", err.Error())
	}
	return stream, err
}

func main() {

	//fmt.Printf("List the protocols you support. Don't lie and don't be shy, go ahead list em all")
	//protos, err := proto.GetProtocols("protocols.txt")

	myHost := host.GetHost("5002")
	fmt.Println(myHost.ID(), myHost.Addrs(), myHost.Network().HardwareAddr)
	PrintProtocols(myHost.Addrs())

	go myHost.StartListening()

	addr, _ := GetAddrFromUser()
	conn, _ := connectTcpIp4(addr)
	session, _ := setupSenderYamux(conn, nil)
	newConn, _ := newStreamYamux(session)

	_, err := newConn.Write([]byte("hello"))
	if err != nil {
		fmt.Println("error in writing:", err)
		return
	}

	/*
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
	*/

	/*myHost := host.GetHost()
	fmt.Println(myHost.ID(), myHost.Addrs(), myHost.Network().HardwareAddr)
	PrintProtocols(myHost.Addrs())

	conn, err := myHost.Listener().Accept()
	if err != nil {
		fmt.Printf("Could not accept connection on %s because %s\n", myHost.Addrs(), err.Error())
	}
	fmt.Println(conn.LocalAddr())*/

}
