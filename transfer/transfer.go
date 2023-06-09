package transfer

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

// Packet represents a custom payload packet that contains the data and its size.
type Packet struct {
	Size int32
	Data []byte
}

const (
	HOST       = "localhost"
	PORT       = "8080"
	TYPE       = "tcp"
	filename   = "random.txt"
	inputPath  = "testingSender/random.txt"
	outputPath = "testingReceiver"
)

func main() {
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
		err = UploadFile(conn, filename, inputPath, 64)
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

}

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
	fmt.Println("SENDER_BLK_PACK_SIZE", packet.Size)
	_, err = buf.Write(packet.Data)
	if err != nil {
		return fmt.Errorf("error writing packet data to buffer: %v", err)
	}
	fmt.Println("SENDER_BLK_PACK_data", packet.Data)
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("error sending packet: %v", err)
	}

	return nil
}

//func ReceiveFile(conn net.Conn, outputPath string) error {
//	reader := bufio.NewReader(conn)
//
//	filename, err := reader.ReadString('\n')
//	if err != nil {
//		return fmt.Errorf("error reading filename: %v", err)
//	}
//	filename = strings.TrimSpace(filename)
//
//	fileSizeStr, err := reader.ReadString('\n')
//	if err != nil {
//		return fmt.Errorf("error reading file size: %v", err)
//	}
//	fileSizeStr = strings.TrimSpace(fileSizeStr)
//	fileSize, err := strconv.ParseInt(fileSizeStr, 10, 64)
//	if err != nil {
//		return fmt.Errorf("error parsing file size: %v", err)
//	}
//
//	err = os.MkdirAll(outputPath, 0755)
//
//	if err != nil {
//		fmt.Println(4)
//		return fmt.Errorf("error creating output file directory: %v", err)
//	}
//
//	cleanFPath := strings.Replace(path.Join(outputPath, filename), "\r", "", -1)
//	cleanFPath = strings.Replace(cleanFPath, "\x00", "", -1)
//	cleanFPath = strings.Replace(cleanFPath, "\x01", "", -1)
//	cleanFPath = strings.Replace(cleanFPath, "\v", "", -1)
//	fmt.Printf("%q", cleanFPath)
//	outputFile, err := os.Create(cleanFPath)
//	if err != nil {
//		fmt.Println(5)
//		return fmt.Errorf("error creating output file: %v", err)
//	}
//
//	defer outputFile.Close()
//
//	var receivedBytes int64
//	for receivedBytes < fileSize {
//		packet, err := receivePacket(conn)
//		if err != nil {
//			return err
//		}
//
//		outputFile.Write(packet.Data)
//		receivedBytes += int64(packet.Size)
//	}
//
//	return nil
//}

// ReceiveFile receives a file from a connected sender over the specified net.Conn.
func ReceiveFile(conn net.Conn, outputPath string) error {
	reader := bufio.NewReader(conn)
	println("reader size: ", reader.Size())

	filename, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(1)
		return fmt.Errorf("error reading filename: %v", err)
	}
	filename = strings.TrimSpace(filename)

	fileSizeStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(2)
		return fmt.Errorf("error reading file size: %v", err)
	}
	fileSizeInfo := strings.TrimSpace(fileSizeStr)
	fileSizeInfo = fileSizeInfo[len(fileSizeInfo)-3:]
	println("File size found ", fileSizeInfo)
	fileSize, err := strconv.ParseInt(fileSizeInfo, 10, 16)
	if err != nil {
		fmt.Println(3)
		return fmt.Errorf("error parsing file size: %v", err)
	}

	err = os.MkdirAll(outputPath, 0755)

	if err != nil {
		fmt.Println(4)
		return fmt.Errorf("error creating output file directory: %v", err)
	}

	cleanFPath := strings.Replace(path.Join(outputPath, filename), "\r", "", -1)
	cleanFPath = strings.Replace(cleanFPath, "\x00", "", -1)
	cleanFPath = strings.Replace(cleanFPath, "\x01", "", -1)
	cleanFPath = strings.Replace(cleanFPath, "\v", "", -1)
	fmt.Printf("%q", cleanFPath)
	outputFile, err := os.Create(cleanFPath)
	if err != nil {
		fmt.Println(5)
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	var receivedBytes int64
	fmt.Println("Full file size", fileSize)
	for receivedBytes < int64(fileSize) {
		packet, err := receivePacket(conn)

		if err != nil {
			fmt.Println(6)
			return err
		}
		fmt.Println("PSIZE", packet.Size)

		if len(packet.Data) == 0 {
			break
		} else {
			outputFile.Write(packet.Data)
			receivedBytes += int64(packet.Size)
		}
	}

	return nil
}

// receivePacket receives a Packet struct from a net.Conn connection.
func receivePacket(conn net.Conn) (Packet, error) {
	var packet Packet

	//d, err := io.ReadAll(conn)
	//
	//if err != nil {
	//	return packet, fmt.Errorf("error decoding packet size: %v", err)
	//}
	//
	//fmt.Println(d)

	err := binary.Read(conn, binary.BigEndian, &packet.Size)
	if err != nil {
		return packet, fmt.Errorf("error decoding packet size: %v", err)
	}

	fmt.Println("RECPACK SIZE", packet.Size)

	packet.Size = 8
	packet.Data = make([]byte, 8)

	_, err = io.ReadAtLeast(conn, packet.Data, int(packet.Size))
	if err != nil {
		return packet, fmt.Errorf("error reading packet data: %v", err)
	}

	packet.Data = []byte(strings.Replace(string(packet.Data), "\x00", "", -1))
	packet.Data = []byte(strings.Replace(string(packet.Data), "\x01", "", -1))
	//packet.Data = []byte(strings.Replace(string(packet.Data), "\x0c", "", -1))

	return packet, nil
}
