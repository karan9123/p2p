package transfer

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// Packet represents a custom payload packet that contains the data and its size.
type Packet struct {
	Size int32
	Data []byte
}

const (
	filename1  = "random.txt"
	inputPath  = "/Users/karan/Documents/Networks/p2p/testingSender/random.txt"
	outputPath = "/Users/karan/Documents/Networks/p2p/testingReceiver"
)

/*func main() {
	// Call uploadFile and receiveFile with your specific parameters
	//	I need to check the command line for the mode = (sender/receiver)

	mode := flag.String("mode", "sender", "To run as a sender or receiver")
	flag.Parse()

	// If sender

	if *mode == "sender" {
		println("Starting program as sender....")
		tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
		println("tcpServer", tcpServer.String())
		if err != nil {
			fmt.Errorf("Error in 'ResolveTCPAddr': %s", err.Error())
		}
		conn, err := net.Dial(TYPE, tcpServer.String())
		if err != nil {
			errstr := fmt.Errorf("Error in 'Dial': %s", err.Error())
			println(errstr.Error())
		} else {
			println("connection ", conn)
		}
		err = UploadFile(conn, filename, inputPath, 8)
		if err != nil {
			fmt.Errorf("Error in 'uploadFile': %s", err.Error())
		}
	}

	// else if receiver
	if *mode == "receiver" {
		println("Starting program as receiver....")
		listen, err := net.Listen(TYPE, HOST+":"+PORT)
		if err != nil {
			errormsg := fmt.Errorf("Error in 'Listen': %s", err.Error())
			println(errormsg.Error())
		}
		defer listen.Close()

		receiverConn, err := listen.Accept()
		if err != nil {
			errormsg := fmt.Errorf("Error in 'Accept': %s", err.Error())
			println(errormsg.Error())
		}

		err = ReceiveFile(receiverConn, outputPath)
		if err != nil {
			errormsg := fmt.Errorf("Error in 'receiveFile': %s", err.Error())
			println(errormsg.Error())

		}
	}

}*/

// UploadFile sends a file to a connected receiver over the specified net.Conn.
func UploadFile(conn net.Conn, filename, inputPath string, blockSize int) error {
	file, err := os.Open(inputPath)
	fmt.Printf("opening %s with inputPath %s \n", filename, inputPath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	println("file size : ", fileInfo.Size())
	println("file name : ", fileInfo.Name())

	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	conn.Write([]byte(filename + "\n"))
	conn.Write([]byte(strconv.FormatInt(fileInfo.Size(), 10) + "\n"))

	buffer := make([]byte, blockSize)

	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading file: %v", err)
		}

		if bytesRead > 0 {
			packet := Packet{
				Size: int32(bytesRead),
				Data: buffer[:bytesRead],
			}

			sendPacket(conn, packet)
		}
	}

	return nil
}

// sendPacket sends a Packet struct over a net.Conn connection.
func sendPacket(conn net.Conn, packet Packet) error {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.BigEndian, packet.Size)
	if err != nil {
		return fmt.Errorf("error encoding packet size: %v", err)
	}

	_, err = buf.Write(packet.Data)
	if err != nil {
		return fmt.Errorf("error writing packet data to buffer: %v", err)
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("error sending packet: %v", err)
	}

	return nil
}

// ReceiveFile receives a file from a connected sender over the specified net.Conn.
func ReceiveFile(conn net.Conn, outputPath string) error {
	reader := bufio.NewReader(conn)

	filename, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading filename: %v", err)
	}
	filename = strings.TrimSpace(filename)
	fmt.Printf("fileName: %s \n", filename)

	fileSizeStr, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading file size: %v", err)
	}
	fmt.Printf("filesizeStr: %s\n", fileSizeStr)
	fileSizeStr = strings.ReplaceAll(fileSizeStr, " ", "")
	fileSizeStr = strings.ReplaceAll(fileSizeStr, "\n", "")
	fmt.Println("f", fileSizeStr)
	k := string(fileSizeStr)
	fmt.Println("k", []byte(k))
	fileSize, err := strconv.Atoi(k)
	if err != nil {
		fmt.Printf("error converting file size because of%s\n", err.Error())

	}
	fmt.Println("filesize:", fileSize)

	outputFile, err := os.Create(outputPath + "/random.txt")
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()
	fmt.Printf("file %s created\n", outputFile)

	var receivedBytes int
	for receivedBytes < fileSize {
		packet, err := receivePacket(conn)
		if err != nil {
			return err
		}

		outputFile.Write(packet.Data)
		receivedBytes += int(packet.Size)
	}

	return nil
}

// receivePacket receives a Packet struct from a net.Conn connection.
func receivePacket(conn net.Conn) (Packet, error) {
	var packet Packet

	err := binary.Read(conn, binary.BigEndian, &packet.Size)
	if err != nil {
		return packet, fmt.Errorf("error decoding packet size: %v", err)
	}

	packet.Data = make([]byte, packet.Size)
	_, err = io.ReadFull(conn, packet.Data)
	if err != nil {
		return packet, fmt.Errorf("error reading packet data: %v", err)
	}

	return packet, nil
}
