package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
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

type DNSQuestion struct {
	Name  string
	Type  int
	Class int
}

type DNSAnswer struct {
	Name   string
	Type   int
	Class  int
	TTL    int
	Length int
	Data   string
}

// Encodes the DNSAnswer into a byte array
func (a *DNSAnswer) Encode() []byte {
	// Split the domain into labels
	labels := strings.Split(a.Name, ".")

	var domainSequence []byte

	for _, label := range labels {
		// byte length of the label
		domainSequence = append(domainSequence, byte(len(label)))

		// append the value or the label after its byte length
		domainSequence = append(domainSequence, label...)
	}

	// append the null byte to terminate the domain
	domainSequence = append(domainSequence, '\x00')

	buffer := make([]byte, 10)
	ip := net.ParseIP(a.Data).To4()
	a.Length = len(ip)

	binary.BigEndian.PutUint16(buffer[0:], uint16(a.Type))
	binary.BigEndian.PutUint16(buffer[2:], uint16(a.Class))
	binary.BigEndian.PutUint32(buffer[4:], uint32(a.TTL))
	binary.BigEndian.PutUint16(buffer[8:], uint16(a.Length))

	result := append(domainSequence, buffer...)
	result = append(result, ip...)

	fmt.Println("Answer:", result)

	return result
}

// Encodes the DNSQuestion into a byte array
func (q *DNSQuestion) Encode() []byte {
	// Split the domain into labels
	labels := strings.Split(q.Name, ".")

	var domainSequence []byte

	for _, label := range labels {
		// byte length of the label
		domainSequence = append(domainSequence, byte(len(label)))

		// append the value or the label after its byte length
		domainSequence = append(domainSequence, label...)
	}

	// append the null byte to terminate the domain
	domainSequence = append(domainSequence, '\x00')

	buffer := make([]byte, 4)

	// Packs the question type and class into buffer
	binary.BigEndian.PutUint16(buffer, uint16(q.Type))
	binary.BigEndian.PutUint16(buffer, uint16(q.Class))

	// append the domain sequence and the buffer
	result := append(domainSequence, buffer...)

	fmt.Println("Question:", result)
	return result
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
func (h *DNSHeader) Encode() []byte {
	buffer := make([]byte, 12)

	binary.BigEndian.PutUint16(buffer[0:], h.ID)
	binary.BigEndian.PutUint16(buffer[2:], h.packHeader())
	binary.BigEndian.PutUint16(buffer[4:], h.QDCount)
	binary.BigEndian.PutUint16(buffer[6:], h.ANCount)
	binary.BigEndian.PutUint16(buffer[8:], h.NSCount)
	binary.BigEndian.PutUint16(buffer[10:], h.ARCount)

	fmt.Println("Header:", buffer)
	return buffer
}

func decodeHeader(data []byte) DNSHeader {
	response := DNSHeader{
		ID: binary.BigEndian.Uint16(data[0:2]),
	}

	value := binary.BigEndian.Uint16(data[2:4])
	response.QR = (value & (1 << 15))!= 0
	response.OPCODE = uint8(value >> 11)
	response.AA = (value & (1 << 10))!= 0
	response.TC = (value & (1 << 9))!= 0
	response.RD = (value & (1 << 8))!= 0
	response.RA = (value & (1 << 7))!= 0
	response.Z = uint8(value >> 4)

	return response
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

		// Parse Received DNS Header
		parsed := decodeHeader([]byte(receivedData))
	
		// DNS Header
		header := DNSHeader{
			ID: parsed.ID,
			QR: true,
			OPCODE: parsed.OPCODE,
			AA: false,
			TC: false,
			RD: parsed.RD,
			RA: false,
			Z: 0,
			RCODE: 4,
			QDCount: 1,
			ANCount: 1,
			NSCount: 0,
			ARCount: 0,
		}

		// DNS Question
		question := DNSQuestion{
			Name: "codecrafters.io",
			Type: 1,
			Class: 1,
		}

		// DNS Answer
		answer := DNSAnswer{
			Name: "codecrafters.io",
			Type: 1,
			Class: 1,
			TTL: 60,
			Data: "8.8.8.8",
		}

		response := append(header.Encode(), question.Encode()...)
		response = append(response, answer.Encode()...)

		fmt.Println("Response:", response)
	
		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
