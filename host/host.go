package host

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
	"p2p/peer"
)

const (
	multiStreamSelect   = "/multistream/1.0.0"
	multiStreamSelectNL = "/multistream/1.0.0\n"
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

	// Listener returns the Listener of the Host
	Listener() net.Listener

	StartListening() (net.Conn, error)

	SenderConn() (net.Conn, error)

	// Mux returns the Mux multiplexing incoming streams to protocol handlers
	//Mux() proto.Switch
}

type MyHost struct {
	peerID   peer.ID
	addrs    ma.Multiaddr
	network  net.Interface
	listener net.Listener
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

func (h *MyHost) Listener() net.Listener {
	return h.listener
}

func (h *MyHost) StartListening() (net.Conn, error) {
	conn, err := h.Listener().Accept()
	if err != nil {
		fmt.Printf("Could not accept connection on %s because %s\n", h.Addrs(), err.Error())
		return nil, err
	}

	buf := make([]byte, 1024)
	i, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Error %s encountered while readed\n", err.Error())
		return nil, err
	}
	//fmt.Printf("reading %d bytes which are: %s\n", i, buf[:i])
	i, err = conn.Read(buf)
	if err != nil {
		fmt.Printf("Error %s encountered while readed\n", err.Error())
		return nil, err
	}
	temp := string(buf[:i-1])

	tempbuf := buf[:i-1]

	fmt.Printf("reading again %d bytes which are: %s and multi is %s\n", i, "|"+temp+"|", "|"+multiStreamSelect+"|")
	if string(tempbuf[len(tempbuf)-18:]) != multiStreamSelect {
		fmt.Printf("Sender does not support multisteam select, Abort connection\n")
		return nil, errors.New("sender does not multi-stream-select")
	}
	fmt.Printf("Sender supports multisteam select. Hooray!!!\n")
	return conn, nil

}

func (h *MyHost) SenderConn() (net.Conn, error) {
	addr, err := GetAddrFromUser()
	if err != nil {
		return nil, err
	}
	conn, err := connectTcpIp4(addr)
	if err != nil {
		return nil, err
	}
	session, err := setupSenderYamux(conn, nil)
	if err != nil {
		return nil, err
	}
	newConn, err := newStreamYamux(session)
	if err != nil {
		return nil, err
	}
	_, err = newConn.Write([]byte(multiStreamSelectNL))
	if err != nil {
		fmt.Printf("Could not write on stream because %s\n", err.Error())
		return nil, err
	}
	return newConn, nil
}

func getMyMultiaddr(inface, tcpPort string) (*net.Interface, ma.Multiaddr, error) {
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
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipAddress = ipnet.IP
				break
			}
		}
	}

	addrStr := "/ip4/" + ipAddress.String() + "/tcp/" + tcpPort

	// Parsing the entered multi-address./ip4/127.0.0.1/tcp/8000
	addr, err := ma.NewMultiaddr(addrStr)
	if err != nil {
		fmt.Printf("Invalid multiaddress: %s, press r to retry or any other key to exit\n", err)
	}
	return iface, addr, nil
}

// ReceiveData reads data from the server and prints received messages.
func _(conn net.Conn) []byte {
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

func GetHost(port string) Host {
	_, pubKey, _ := cr.GenerateKeyPair(1, -1)
	id, _ := peer.GenerateIDFromPubKey(pubKey)

	fmt.Printf("pubkey is %s, id is %s\n", pubKey, id)
	network, addrs, _ := getMyMultiaddr("en0", port)
	ip4, tcpPort, _ := GetIp4TcpFromMultiaddr(addrs)

	listener, err := net.Listen("tcp", ip4+":"+tcpPort)
	if err != nil {
		fmt.Printf("Could not listen on Address %s \n because %s", ip4+":"+tcpPort, err.Error())
	}
	fmt.Println("my ID:", []byte(id))
	fmt.Println("listening on port: ", tcpPort)

	host := &MyHost{
		peerID:   id,
		addrs:    addrs,
		network:  *network,
		listener: listener,
	}
	return host
}

func GetIp4TcpFromMultiaddr(addr ma.Multiaddr) (string, string, error) {
	ipAddr, err := addr.ValueForProtocol(ma.P_IP4)
	if err != nil {
		fmt.Printf("Could not get ipv4 address from the given multiadress because %s \n", err.Error())
		return "", "", err
	}

	tcpPort, err := addr.ValueForProtocol(ma.P_TCP)
	if err != nil {
		fmt.Printf("Could not get TCP Port from the given multiadress %s \n", err.Error())
		return ipAddr, "", err
	}
	return ipAddr, tcpPort, nil
}

// connectTcpIp4 returns tcp connection from multiaddr
func connectTcpIp4(addr ma.Multiaddr) (*net.TCPConn, error) {

	ipAddr, tcpPort, err := GetIp4TcpFromMultiaddr(addr)
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
