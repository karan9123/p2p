# libp2p-tcp-yamux-Tftp

 libp2p-tcp-yamux-ftp is a library that provides a libp2p implementation using TCP over IP as its transport layer and yamux as its multiplexer. The library also includes a trivial FTP implementation at its application stack.

This library can be used to create decentralized applications that utilize the libp2p network stack for peer-to-peer communication. The library provides an easy-to-use API for creating and managing libp2p nodes, establishing connections, and transferring data.

## Libp2p Features
- Peer-to-peer communication: libp2p allows nodes to communicate directly with each other without relying on a central server.
- Transport agnostic: libp2p supports multiple transport protocols such as TCP/IP, UDP, and WebSockets, making it flexible for various network environments.
- Modular architecture: libp2p is designed to be modular, allowing developers to choose and combine different protocols and components based on their specific requirements.
- Secure communication: libp2p supports encryption and authentication mechanisms, ensuring secure communication between nodes.
- NAT traversal: libp2p includes NAT traversal techniques to enable direct communication between nodes behind NATs and firewalls.

## Implementation Features

- Uses TCP over IP as its transport layer
- Uses yamux as its multiplexer
- Includes a trivial-FTP implemented along the lines of RFC-1350 at its application stack
- Provides an easy-to-use API for creating and managing libp2p nodes


## Getting Started

### Installation
The library can be imported:

```
import "github.com/karan9123/p2p"
```
or can be cloned from Github
```
git clone https://github.com/karan9123/p2p.git
```

## Usage

### Creating a libp2p node

To create a libp2p node, you can use the following code:

To run a receiver:
```
host = host.getHost()
receiverMethod(myHost)
```
To run a sender:
```
host = host.getHost()
senderMethod(host)
```
This code dials a peer by its [multiaddrs](https://github.com/multiformats/multiaddr)  creates a TP 
connection over the established libp2p connection. The connection is then used to transfer data over multiple streams using [Yamux](https://github.com/hashicorp/yamux#readme) multiplexer.

# Host
This code defines an interface and implementation for a Host object representing a 
single libp2p node in a peer-to-peer network. The Host object participates in a p2p 
network, implementing protocols or providing services, and handles requests like a server, 
while issuing requests like a client. The code also includes functionality for starting to 
listen for incoming connections, creating a new outgoing connection, starting to receive a 
file, and starting to transfer a file.

# Crypto
This package contains functions and types related to cryptography. 
It provides functionality to generate, manage, and use various types of cryptographic keys 
such as RSA, Ed25519, and ECDSA. It also provides a mechanism to marshal and unmarshal keys 
in a standardized format.

The package defines three types of KeyType constants: 
RSA, Ed25519, and ECDSA, which are used throughout the package to denote the 
type of cryptographic key.

# Transfer

Transfer is a Go package that provides functionality for a Trivial FTP over a network using TCP/IP.
It consists of two main functions, UploadFile() and ReceiveFile(), that respectively send and receive files.

Trivial FTP (TFTP) is a simple file transfer protocol that operates on top of the User Datagram Protocol (UDP) 
using port number 69. It was initially designed for bootstrapping diskless machines and for 
transferring configuration files in a low-security setting. The TFTP protocol is described in RFC 1350. However, this code works with TCP/IP which 
ensures no data loss.



