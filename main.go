package main

import (
	"context"
	"fmt"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
	"p2p/host"
	tr "p2p/transfer"
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

	myHost := host.GetHost("5007")
	fmt.Println(myHost.ID(), myHost.Addrs(), myHost.Network().HardwareAddr)
	PrintProtocols(myHost.Addrs())

	//testingTransferSender(myHost)
	//testingTransferReceiver(myHost)

	//receiverMethod(myHost)
	senderMethod(myHost)

}

func testNoise() {
	// Let there be Alice, Bob, and Charlie.

	alice, err := noise.NewNode()
	if err != nil {
		panic(err)
	}

	bob, err := noise.NewNode()
	if err != nil {
		panic(err)
	}

	charlie, err := noise.NewNode()
	if err != nil {
		panic(err)
	}

	// Alice, Bob, and Charlie are following an overlay network protocol called Kademlia to discover, interact, and
	// manage each others peer connections.

	ka, kb, kc := kademlia.New(), kademlia.New(), kademlia.New()

	alice.Bind(ka.Protocol())
	bob.Bind(kb.Protocol())
	charlie.Bind(kc.Protocol())

	if err := alice.Listen(); err != nil {
		panic(err)
	}

	if err := bob.Listen(); err != nil {
		panic(err)
	}

	if err := charlie.Listen(); err != nil {
		panic(err)
	}

	// Have Bob and Charlie learn about Alice. Bob and Charlie do not yet know of each other.

	if _, err := bob.Ping(context.TODO(), alice.Addr()); err != nil {
		panic(err)
	}

	if _, err := charlie.Ping(context.TODO(), bob.Addr()); err != nil {
		panic(err)
	}

	// Using Kademlia, Bob and Charlie will learn of each other. Alice, Bob, and Charlie should
	// learn about each other once they run (*kademlia.Protocol).Discover().

	fmt.Printf("Alice discovered %d peer(s).\n", len(ka.Discover()))
	fmt.Printf("Bob discovered %d peer(s).\n", len(kb.Discover()))
	fmt.Printf("Charlie discovered %d peer(s).\n", len(kc.Discover()))

}

func testingTransferReceiver(myHost host.Host) {
	conn, err := myHost.StartListening()
	if err != nil {
		fmt.Println("Couldn't listen due to ", err.Error())
	}
	err = tr.ReceiveFile(conn, outputPath)
}

func testingTransferSender(myHost host.Host) {
	sendingConn, _ := myHost.SenderConn()
	err := tr.UploadFile(sendingConn, filename, inputPath, 8)
	if err != nil {
		fmt.Printf("error in upload file due to %s \n", err.Error())
	}
}

func receiverMethod(myHost host.Host) {
	_, _ = myHost.StartListening()
}

func senderMethod(myHost host.Host) {
	_, _ = myHost.SenderConn()
	//i, err := sendingConn.Write([]byte("hello\n"))
	//if err != nil {
	//	return
	//}
	//fmt.Printf("Wrote %d bytes\n", i)

}
