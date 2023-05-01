/*
Package host provides a Host object representing a single libp2p node in a peer-to-peer network.
Host is an object participating in a p2p network, which implements protocols or provides services.
It handles requests like a Server, and issues requests like a Client.
It is called Host because it is both Server and Client (and Peer may be confusing).
This library provides an object MyHost which implements Host interface and its methods.
*/
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
	tr "p2p/transfer"
)

// The string used for the multistream select protocol negotiation
const (
	multiStreamSelect   = "/multistream/1.0.0"
	multiStreamSelectNL = "/multistream/1.0.0\n"
	filename            = "random.txt"
	inputPath           = "/Users/karan/Documents/Networks/p2p/testingSender/random.txt"
	outputPath          = "/Users/karan/Documents/Networks/p2p/testingReceiver"
	ifaceName           = "eth0" //change it to "en0" if you are on a Mac
)

// Host represents a single libp2p node in a peer-to-peer network.
//
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
	// Connection returns the TCPConnection to peer of the host
	Connection() *net.TCPConn
	// StartListening listens to incoming connections on the listener
	StartListening() (net.Conn, error)
	// StartSending creates a new outgoing connection
	StartSending() (net.Conn, error)
	// NewConn returns a new stream if session is not nil
	NewConn() (net.Conn, error)
	// NewStream returns a new
	NewStream() (*yamux.Stream, error)
	StartReceiveFile()
	StartTransferFile()
}

// MyHost is an implementation of Host interface
type MyHost struct {
	peerID     peer.ID
	addrs      ma.Multiaddr
	network    net.Interface
	listener   net.Listener
	connection *net.TCPConn
	lConn      *net.Conn
	session    *yamux.Session
}

// ID returns the peer ID associated with this host
func (h *MyHost) ID() peer.ID {
	return h.peerID
}

// Addrs returns the listen addresses of the host
func (h *MyHost) Addrs() ma.Multiaddr {
	return h.addrs
}

// Network returns the Network interface of the host
func (h *MyHost) Network() net.Interface {
	return h.network
}

// Listener returns the listener of the host
func (h *MyHost) Listener() net.Listener {
	return h.listener
}

// Connection returns the TCPConnection to peer of the host
func (h *MyHost) Connection() *net.TCPConn {
	return h.connection
}

func (h *MyHost) NewStream() (*yamux.Stream, error) {
	con := h.connection
	if con == nil {
		fmt.Printf("No connection found to stream over")
		return nil, errors.New("no connection found to stream over")
	}

	yamuxSession, err := yamux.Server(con, nil)
	defer func(yamuxSession *yamux.Session) {
		err := yamuxSession.Close()
		if err != nil {
			fmt.Printf("error %s/n", err.Error())
		}
	}(yamuxSession)
	//yamuxSession, err := setupSenderYamux(con, nil)

	if err != nil {
		fmt.Printf("Error opening a new stream because of %s\n", err)
		return nil, err
	}

	stream, err := yamuxSession.AcceptStream()
	if err != nil {
		// Ignore the error if it's related to the session being closed
		if err != yamux.ErrSessionShutdown {
			panic(err)
		}
	}
	defer func(stream *yamux.Stream) {
		err := stream.Close()
		if err != nil {
			fmt.Printf("error::%s \n", err.Error())
		}
	}(stream)
	return stream, nil

}

// StartListening listens to incoming connections on the listener
func (h *MyHost) StartListening() (net.Conn, error) {
	conn, err := h.Listener().Accept()
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		fmt.Printf("Could not resolve Conn to TCPConn\n")
	}
	h.connection = tcpConn
	fmt.Printf("conected to %s \n", conn.RemoteAddr())
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
	//temp := string(buf[:i-1])

	tempbuf := buf[:i-1]

	fmt.Printf("reading again %d bytes which are: %s and multi is %s\n", i, "|"+string(tempbuf[len(tempbuf)-18:])+"|", "|"+multiStreamSelect+"|")
	if string(tempbuf[len(tempbuf)-18:]) != multiStreamSelect {
		fmt.Printf("Sender does not support multisteam select, Abort connection\n")
		return nil, errors.New("sender does not multi-stream-select")
	}
	fmt.Printf("Sender supports multisteam select. Hooray!!!\n")

	return conn, nil

}

// StartSending returns a new connection to the destination.
func (h *MyHost) StartSending() (net.Conn, error) {
	// Get the destination address from the user.
	addr, err := GetAddrFromUser()
	if err != nil {
		return nil, err
	}
	// Connect to the destination address over TCP/IP.
	conn, err := connectTcpIp4(addr)
	if err != nil {
		return nil, err
	}

	h.connection = conn

	// Set up a new yamux session for the connection.
	session, err := setupSenderYamux(conn, nil)
	if err != nil {
		return nil, err
	}

	h.session = session

	// Open a new stream over the yamux session.
	newConn, err := newConnYamux(session)
	if err != nil {
		return nil, err
	}

	// Write the stream select message to the new stream.
	_, err = newConn.Write([]byte(multiStreamSelectNL))
	if err != nil {
		fmt.Printf("Could not write on stream because %s\n", err.Error())
		return nil, err
	}

	return newConn, nil
}

// NewConn returns a new stream if session is not nil
func (h *MyHost) NewConn() (net.Conn, error) {
	if h.connection == nil {
		return nil, errors.New("no Connection to multiplex over")
	}

	conn, err := newConnYamux(h.session)
	if err != nil {
		fmt.Printf("Could not create a new stream because of %s \n", err.Error())
		return nil, err
	}
	return conn, nil
}

func (h *MyHost) StartReceiveFile() {
	conn, err := h.Listener().Accept()
	fmt.Printf("conected to %s \n", conn.RemoteAddr())
	if err != nil {
		fmt.Printf("Could not accept connection on %s because %s\n", h.Addrs(), err.Error())
	}
	err = tr.ReceiveFile(conn, outputPath)
	if err != nil {
		err := fmt.Errorf("rrror in 'receiveFile': %s", err.Error())
		println(err.Error())

	}

}
func (h *MyHost) StartTransferFile() {
	// Get the destination address from the user.
	addr, err := GetAddrFromUser()
	if err != nil {
		fmt.Printf("Unable to get add because %s\n", err.Error())
	}
	// Connect to the destination address over TCP/IP.
	conn, err := connectTcpIp4(addr)
	if err != nil {
		fmt.Printf("Unable to connect TCP because %s\n", err.Error())
	}

	h.connection = conn

	// Set up a new yamux session for the connection.
	session, err := setupSenderYamux(conn, nil)
	if err != nil {
		fmt.Printf("Unable to set up Yamux because %s\n", err.Error())
	}

	// Open a new stream over the yamux session.
	newConn, err := newConnYamux(session)
	if err != nil {
		fmt.Printf("Unable to send over streaed connection because %s\n", err.Error())
	}

	// Write the stream select message to the new stream.

	err = tr.UploadFile(newConn, filename, inputPath, 8)
	if err != nil {
		fmt.Printf("Could not write on stream because %s\n", err.Error())
	}
}

// getMyMultiaddr returns the IP address and multiaddress of the specified network interface and TCP port.
func getMyMultiaddr(inface, tcpPort string) (*net.Interface, ma.Multiaddr, error) {
	// Get the network interface by name.
	iface, err := net.InterfaceByName(inface)
	if err != nil {
		fmt.Printf("Could not find the interface because %s \n", err.Error())
		return nil, nil, err
	}
	// Get the interface's addresses.
	addrs, err := iface.Addrs()
	if err != nil {
		fmt.Println(err)
		return iface, nil, err
	}

	/*	// Print the network interface's addresses.
		for _, add := range addrs {
			fmt.Println(add.Network(), add.String())
		}*/

	// Find the first IPv4 address that is not a loopback address.
	var ipAddress net.IP
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipAddress = ipnet.IP
				break
			}
		}
	}

	// Construct a multiaddress from the IP address and TCP port.
	addrStr := "/ip4/" + ipAddress.String() + "/tcp/" + tcpPort
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

// GetHost returns a new Host object after generating a key pair, obtaining the local IPv4 address and TCP port,
// and setting up a TCP listener.
// Parameters:
// - port: a string representing the port number to listen to.
// Returns:
// - a Host object representing the local host.
func GetHost(port string) Host {
	// Generate a key pair and obtain the peer ID.
	_, pubKey, _ := cr.GenerateKeyPair(1, -1)
	id, _ := peer.GenerateIDFromPubKey(pubKey)

	// Obtain the local multiaddress, IPv4 address, and TCP port.

	network, addrs, _ := getMyMultiaddr(ifaceName, port)

	ip4, tcpPort, _ := GetIp4TcpFromMultiaddr(addrs)

	// Set up a TCP listener on the local IPv4 address and TCP port.
	listener, err := net.Listen("tcp", ip4+":"+tcpPort)
	if err != nil {
		// Print an error message if the TCP listener could not be set up.
		fmt.Printf("Could not listen on Address %s \n because %s", ip4+":"+tcpPort, err.Error())
	}

	// Create a new MyHost object with the obtained parameters and return it as a Host object.
	host := &MyHost{
		peerID:     id,
		addrs:      addrs,
		network:    *network,
		listener:   listener,
		connection: nil,
		session:    nil,
		lConn:      nil,
	}
	return host
}

// GetIp4TcpFromMultiaddr extracts the IPv4 address and TCP port from a given multiaddress.
// Parameters:
// - addr: a Multiaddr object representing the multiaddress to extract from.
// Returns:
// - a string representing the IPv4 address.
// - a string representing the TCP port.
// - an error object if either the IPv4 address or TCP port could not be extracted.
func GetIp4TcpFromMultiaddr(addr ma.Multiaddr) (string, string, error) {
	// Extract the IPv4 address from the given multiaddress.
	ipAddr, err := addr.ValueForProtocol(ma.P_IP4)
	if err != nil {
		// Print an error message if the IPv4 address could not be extracted.
		fmt.Printf("Could not get ipv4 address from the given multiadress because %s \n", err.Error())
		return "", "", err
	}
	// Extract the TCP port from the given multiaddress.
	tcpPort, err := addr.ValueForProtocol(ma.P_TCP)
	if err != nil {
		// Print an error message if the TCP port could not be extracted.
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

func newStream(sess *yamux.Session) (*yamux.Stream, error) {
	k, err := sess.OpenStream()
	if err != nil {
		fmt.Printf("Error opening a new stream because of %s\n", err)
		return nil, err
	}
	return k, nil

}

// newStream takes in a yamux session and returns a multiplexed connection
func newConnYamux(session *yamux.Session) (net.Conn, error) {
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
