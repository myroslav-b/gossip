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

/*type contentInteractor interface {
	Hear([]byte, int) error
	Say([]byte, int) error
}*/

const powerRingBuf = 64

var ErrContentBufferEmpty = errors.New("Content Buffer Empty")

type Content struct {
	sync.Mutex
	reader        io.Reader
	name          string
	ring          ringBuf
	queueRead     []byte
	markQueueRead int
	//cont   []byte
	//rCont  bytes.Reader
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
	//c.cont = []byte{}
	//c.rCont = *bytes.NewReader(c.cont)
	c.Unlock()
	return &c
}

/*func New(r io.Reader, name string) *Content {
	var c Content
	var err error
	c.reader = r
	c.Lock()
	switch {
	case name == "":
		c.name, err = os.Hostname()
		if err != nil {
			c.name = "noname"
		}
	default:
		c.name = name
	}

	//var str string
	//fmt.Fscan(c.reader, &str)
	//c.cont = c.name + str
	c.cont = []byte{}
	//c.cont = ""
	c.rCont = *bytes.NewReader(c.cont)
	//c.rCont = *strings.NewReader(c.cont)
	c.Unlock()
	return &c
}*/

/*func (c *Content) input(errInput *error) {
	haltInput := false
	b := true
	bts := []byte{}
	scanner := bufio.NewScanner(c.reader)
	c.Lock()
	cont := []byte(strings.Join([]string{c.name, ": "}, ""))
	c.Unlock()
	for scanner.Scan() {
		if haltInput = atomicHaltInput.Load(); haltInput == true {
			return
		}
		bts = scanner.Bytes()
		c.Lock()
		cont = append(cont, bts...)
		b = c.ring.add(cont)
		if !b {
			log.Printf("[WARNING] skipped data: content buffer full")
		}
		c.Unlock()
	}
	err := scanner.Err()
	if err != nil {
		errInput = &err
		return
	}
}*/

/*func (c *Content) input() {
	bts := []byte{}
	scanner := bufio.NewScanner(c.reader)
	c.Lock()
	c.cont = []byte(strings.Join([]string{c.name, ": "}, ""))
	c.Unlock()
	for {
		if scanner.Scan() {
			bts = scanner.Bytes()
			c.Lock()
			//c.cont = []byte(strings.Join([]string{c.name, ": "}, ""))
			c.cont = append(c.cont, bts...)
			c.rCont.Reset(c.cont)
			c.Unlock()
		}
	}
}*/

func (c *Content) Manager(ctx context.Context) error {
	//var atomicHaltInput atomic.Value
	var errInput error
	//go c.input(&errInput)
	go func() {
		for {
			scanner := bufio.NewScanner(c.reader)
			c.Lock()
			prefix := []byte(strings.Join([]string{c.name, ": "}, ""))
			c.Unlock()
			for scanner.Scan() {
				/*if atomicHaltInput.Load().(bool) {
					return
				}*/
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
				errInput = err
				return
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			//haltInput = true
			//atomicHaltInput.Store(true)
			return ctx.Err()
		default:
			if errInput != nil {
				return errInput
			}
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
			//buf = buf[:0]
			c.markQueueRead = 0
			return 0, ErrContentBufferEmpty
		}
	}

	switch {
	case c.markQueueRead+len(buf) < len(c.queueRead):
		//buf = append(buf[:0], c.queueRead[c.markQueueRead:len(buf)]...)
		l := len(buf)
		for i := 0; i < l; i++ {
			buf[i] = c.queueRead[c.markQueueRead+i]
		}
		c.markQueueRead = c.markQueueRead + len(buf)
		return l, nil
	case c.markQueueRead+len(buf) >= len(c.queueRead):
		//buf = append(buf[:0], c.queueRead[c.markQueueRead:len(c.queueRead)]...)
		l := len(c.queueRead) - c.markQueueRead
		for i := 0; i < l; i++ {
			buf[i] = c.queueRead[c.markQueueRead+i]
		}
		//buf = append(buf,)
		c.markQueueRead = 0
		return l, io.EOF
	default:
		return 0, errors.New("ErrImpossible")
	}
}

/*func (c *Content) Read(buf []byte) (n int, err error) {
	c.Lock()
	defer c.Unlock()
	bts, b := c.ring.sub()
	fmt.Println(bts)
	if !b {
		//buf = bts
		buf = []byte{}
		return 0, io.EOF //ErrContentBufferEmpty
	}
	buf = bts
	return len(buf), io.EOF
}*/
