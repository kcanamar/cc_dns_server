package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	
	/*
		net.ResolveUDPAddr is used to resolve a UDP address string (like "127.0.0.1:2053") into a net.UDPAddr struct. 
		The "udp" network type is passed as the first argument to specify this should be a UDP address.
		The address string "127.0.0.1:2053" specifies where we want to bind in this case localhost on port 2053.
		net.ResolveUDPAddr returns a net.UDPAddr struct representing that address and an error.
		If no error, we now have a valid UDP address struct that can be used to bind to the provided address.
	*/
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}
	
	/*
		net.ListenUDP is used to create a UDP listener. 
		The first argument "udp" specifies we want UDP, not TCP. 
		The second argument udpAddr is a net.UDPAddr struct representing the address we want to bind the listener to.
		The UDP listener returns a net.UDPConn which represents the listening UDP socket, and an error.
		The error is checked - if not nil, it means binding failed and we print the error and return.
		If no error, the bind was successful. We now have a valid UDP connection that can receive packets sent to that address.
		The defer udpConn.Close() will make sure the connection is closed when the main function returns an error or when the program exits.
	*/
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()
	
	/*
		Create a buffer (buf) to hold incoming UDP packet data.
		A preallocated byte buffer for working with binary data.
		Able to read from a UDP socket using udpConn.ReadFromUDP() and pass in buf to store the incoming packet data.
	*/
	buf := make([]byte, 512)
	
	/*
		Enter an infinite loop to continuously listen for UDP packets.
	*/
	for {

		/*
			Call udpConn.ReadFromUDP() to read a UDP packet into the buffer. 
			This returns the number of bytes read (size), 
			the source address (source), 
			and an error if any.
		*/
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}
	
		/*
			Extract the actual packet data from the buffer by slicing it to the size read. 
			Convert this to a string to print it out.
			Print out the size, source address and packet data received for debugging.
		*/
		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)
	
		// Create an empty "dummy" response packet.
		response := []byte{}
	
		/*
			Call udpConn.WriteToUDP() to send the response packet back to the source address that sent the original packet.
		*/
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
