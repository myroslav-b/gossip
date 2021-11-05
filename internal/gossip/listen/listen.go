package listen

import (
	"context"
	"fmt"
	"io"
	"net"
)

const (
	maxDatagramSize = 8192
)

func Listen(ctx context.Context, wrt io.Writer, addrStr string) error {
	addr, err := net.ResolveUDPAddr("udp4", addrStr)
	if err != nil {
		return err
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return err
	}

	conn.SetReadBuffer(maxDatagramSize)

	buffer := make([]byte, maxDatagramSize)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			//buffer := make([]byte, maxDatagramSize)
			numBytes, _, err := conn.ReadFromUDP(buffer)
			fmt.Fprintln(wrt, string(buffer[:numBytes]))
			if err != nil {
				return err
			}

			//log.Print(string(buffer[:numBytes]))
		}
	}
}
