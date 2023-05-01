## Transfer
Transfer is a Go package that provides functionality for a Trivial FTP over a network using TCP/IP. 
It consists of two main functions, UploadFile() and ReceiveFile(), that respectively send and receive files.

Installation
To use the Transfer package, you need to have Go installed on your machine. You can install it by following the instructions on the official Go website.

Once you have installed Go, you can install the Transfer package by running the following command:

## Usage
The Transfer package is used by calling the UploadFile() and ReceiveFile() functions with the appropriate parameters.

### Sending a file
To send a file, you need to specify the following parameters:

- conn: the net.Conn connection to the receiver.
- filename: the name of the file to be sent.
- inputPath: the path to the file to be sent.
- blockSize: the size of the buffer to be used when sending the file.


### Receiving a file
To receive a file, you need to specify the following parameters:

- conn: the net.Conn connection from the sender.
- outputPath: the path to the directory where the received file should be saved.