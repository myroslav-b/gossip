package contentmanager

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

const powerRingBuf = 64

var ErrContentBufferEmpty = errors.New("Content Buffer Empty")

type Content struct {
	sync.Mutex
	reader        io.Reader
	name          string
	ring          ringBuf
	queueRead     []byte
	markQueueRead int
}

func New(r io.Reader, name string) *Content {
	var c Content
	var err error
	c.reader = r
	c.Lock()
	c.ring = *newRingBuf(powerRingBuf)
	switch {
	case name == "":
		c.name, err = os.Hostname()
		if err != nil {
			c.name = "noname"
		}
	default:
		c.name = name
	}
	c.Unlock()
	return &c
}

func (c *Content) Manager(ctx context.Context) error {
	chErrInput := make(chan error)
	go func() {
		scanner := bufio.NewScanner(c.reader)
		for {
			c.Lock()
			prefix := []byte(strings.Join([]string{c.name, ": "}, ""))
			c.Unlock()
			for scanner.Scan() {
				bts := scanner.Bytes()
				cont := prefix
				cont = append(cont, bts...)
				c.Lock()
				b := c.ring.add(cont)
				c.Unlock()
				if !b {
					log.Printf("[WARNING] skipped data: content buffer full")
				}
			}
			err := scanner.Err()
			if err != nil {
				chErrInput <- err
				return
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case errInput := <-chErrInput:
			return errInput
		}
	}
}

func (c *Content) Read(buf []byte) (n int, err error) {
	b := false
	c.Lock()
	defer c.Unlock()

	if c.markQueueRead == 0 {
		c.queueRead, b = c.ring.sub()
		if !b {
			c.markQueueRead = 0
			return 0, ErrContentBufferEmpty
		}
	}

	switch {
	case c.markQueueRead+len(buf) < len(c.queueRead):
		l := len(buf)
		for i := 0; i < l; i++ {
			buf[i] = c.queueRead[c.markQueueRead+i]
		}
		c.markQueueRead = c.markQueueRead + len(buf)
		return l, nil
	case c.markQueueRead+len(buf) >= len(c.queueRead):
		l := len(c.queueRead) - c.markQueueRead
		for i := 0; i < l; i++ {
			buf[i] = c.queueRead[c.markQueueRead+i]
		}
		c.markQueueRead = 0
		return l, io.EOF
	default:
		return 0, errors.New("ErrImpossible")
	}
}
