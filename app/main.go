package main

import (
	"fmt"
	"net"
	"encoding/binary"
)

type DNSHeader struct {
	ID      uint16
	QR      bool
	OPCODE	uint8
	AA      bool
	TC      bool
	RD      bool
	RA      bool
	Z       uint8
	RCODE   uint8
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

// Packs the header into a uint16 values
func (h *DNSHeader) packHeader() uint16 {
	var header uint16 = 0
	if h.QR { header += 1 << 15 }
	header |= uint16(h.OPCODE) << 11
	if h.AA { header += 1 << 10 }
	if h.TC { header += 1 << 9 }
	if h.RD { header += 1 << 8 }
	if h.RA { header += 1 << 7 }
	header |= uint16(h.Z) << 6
	header |= uint16(h.RCODE)
	return header
}

// Encodes the header into a byte array
func (h *DNSHeader) encodeDNSHeader() []byte {
	buffer := make([]byte, 12)

	binary.BigEndian.PutUint16(buffer[0:], h.ID)
	binary.BigEndian.PutUint16(buffer[2:], h.packHeader())
	binary.BigEndian.PutUint16(buffer[4:], h.QDCount)
	binary.BigEndian.PutUint16(buffer[6:], h.ANCount)
	binary.BigEndian.PutUint16(buffer[8:], h.NSCount)
	binary.BigEndian.PutUint16(buffer[10:], h.ARCount)

	return buffer
}

func main() {
	// UDP ADDRESS and PORT
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}
	
		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)
	
		// DNS Header
		header := DNSHeader{
			ID: 1234,
			QR: true,
			OPCODE: 0,
			AA: false,
			TC: false,
			RD: false,
			RA: false,
			Z: 0,
			RCODE: 0,
			QDCount: 0,
			ANCount: 0,
			NSCount: 0,
			ARCount: 0,
		}

		response := header.encodeDNSHeader()
	
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
