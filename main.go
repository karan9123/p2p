package main

import (
	"bufio"
	"fmt"
	ma "github.com/multiformats/go-multiaddr"
	"os"
	"p2p/host"
)

// PrintProtocols takes in a multi address and returns the list of protocols in the address
func PrintProtocols(addr ma.Multiaddr) {
	fmt.Printf("Protocols in Multiaddr %s are:\n", addr.String())
	for _, protocol := range addr.Protocols() {
		fmt.Println("PName", protocol.Name)
	}
}

const (
	filename   = "random.txt"
	inputPath  = "/Users/karan/Documents/Networks/p2p/testingSender/random.txt"
	outputPath = "/Users/karan/Documents/Networks/p2p/testingReceiver"
)

func main() {

	myHost := host.GetHost("5031")
	fmt.Printf("ID: %s\nAddress: %s\n", myHost.ID(), myHost.Addrs())

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	if input == "rec" {
		receiverMethod(myHost)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()
		if input == "go" {
			fmt.Printf("in receive go\n")
			receiveAgain(myHost)
		}
	}

	if input == "send" {
		senderMethod(myHost)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()
		if input == "go" {
			sendAgain(myHost)
		}
	}

	/*scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	if input == "rec" {
		myHost.StartReceiveFile()
	} else if input == "send" {
		myHost.StartTransferFile()
	}*/

}

func receiverMethod(myHost host.Host) {
	_, _ = myHost.StartListening()
}

func senderMethod(myHost host.Host) {
	_, _ = myHost.StartSending()
}

func sendAgain(myHost host.Host) {
	sendingConn, err := myHost.NewConn()
	i, err := sendingConn.Write([]byte("hello, This is sent over a new " +
		"Stream using the same connection. Isn't it efficient???\n"))
	if err != nil {
		return
	}
	fmt.Printf("Wrote %d bytes\n", i)
}

func receiveAgain(myHost host.Host) {
	fmt.Printf("In receive again %s \n", myHost.Listener().Addr())
	stream, err := myHost.NewStream()
	if err != nil {
		fmt.Printf("Can't listen on %s because %s \n", myHost.Listener().Addr().String(), err.Error())
	}
	fmt.Printf("Listening on %s \n", stream.LocalAddr())
	buf := make([]byte, 1024)
	i, err := stream.Read(buf)
	if err != nil {
		fmt.Printf("Error %s encountered while readin\n", err.Error())
	}
	fmt.Printf("reading %d bytes which are: %s\n", i, buf[:i])

}
