package talk

import (
	"context"
	"io"
	"net"
	"time"

	contentmanager "github.com/myroslav-b/gossip/internal/gossip/contentManager"
)

type TraficControler interface {
	Speak() bool
}

func Talk(ctx context.Context, rdr io.Reader, tc TraficControler, addrStr string, freq uint32) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if tc.Speak() {
				err := sendMessage(rdr, addrStr)
				if err != nil {
					//log.Printf("[ERROR] %+v", err)
					return err
				}
			}
			time.Sleep(time.Duration(1000000/freq) * time.Microsecond)
		}
	}
}

func sendMessage(rdr io.Reader, addrStr string) error {
	addr, err := net.ResolveUDPAddr("udp4", addrStr)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return err
	}

	buf, err := rdrToBuf(rdr)
	if (err != nil) && (err != contentmanager.ErrContentBufferEmpty) {
		return err
	}

	if err != contentmanager.ErrContentBufferEmpty {
		conn.Write(buf)
	}

	return nil
}

func rdrToBuf(rdr io.Reader) ([]byte, error) {
	nBuf := 16
	buffer := make([]byte, nBuf)
	bigBuffer := make([]byte, 0)
	for {
		numBytes, err := rdr.Read(buffer)
		bigBuffer = append(bigBuffer, buffer[:numBytes]...)
		switch {
		case err == io.EOF:
			return bigBuffer, nil
		case err == contentmanager.ErrContentBufferEmpty:
			return bigBuffer, contentmanager.ErrContentBufferEmpty
		case (err != io.EOF) && (err != nil) && (err != contentmanager.ErrContentBufferEmpty):
			return nil, err
		}
	}
}
