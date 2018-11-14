package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {

	// listen on all interfaces
	listener, err := net.Listen("tcp", "127.0.0.1:3307")
	if err != nil {
		panic(err)
	}

	// fmt.Println("Starting network proxy from port 3307->3306")
	// Wait for a new connection
	inbound, err := listener.Accept()
	if err != nil {
		panic(err)
	}

	// Define a (not yet open) downstream connection
	var outbound net.Conn
	buffer := make([]byte, 1024*1024)
	for {
		// Read the available data from the inbound connection
		bytesRead, err := bufio.NewReader(inbound).Read(buffer)
		// fmt.Printf("bytes read: %d\n", bytesRead)
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}

		// Print this information to stdout for the parent process to make use of
		// it
		fmt.Print(string(buffer[0:bytesRead]))

		// If this is the first message received, open the downstream connection
		if outbound == nil {
			outbound, err = net.Dial("tcp", "127.0.0.1:3306")
			if err != nil {
				panic(err)
			}
		}

		// Pass the exact bytes on to the downstream connection
		outbound.Write(buffer[0:bytesRead])
	}
}
