package listen

import (
	"context"
	"log"
	"net"
)

const (
	maxDatagramSize = 8192
)

func Listen(ctx context.Context /*wtr io.Writer,*/, addrStr string) error {
	addr, err := net.ResolveUDPAddr("udp4", addrStr)
	if err != nil {
		return err
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return err
	}

	conn.SetReadBuffer(maxDatagramSize)

	for {
		buffer := make([]byte, maxDatagramSize)
		numBytes, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			return err
		}

		log.Print(string(buffer[:numBytes]))
	}
}
