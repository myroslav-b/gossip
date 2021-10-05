package contentmanager

import (
	"context"
	"strings"
	"sync"
)

/*type contentInteractor interface {
	Hear([]byte, int) error
	Say([]byte, int) error
}*/

type Content struct {
	sync.Mutex
	cont       string
	contReader strings.Reader
}

func New() *Content {
	var c Content
	c.Lock()
	c.cont = "а аб абр абра абрак абрака абракад абракада абракадаб абракадабр абракадабра"
	c.contReader = *strings.NewReader(c.cont)
	c.Unlock()
	return &c
}

func (c *Content) Manager(ctx context.Context) error {
	for {
		//fmt.Print(c.contReader.Len())
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			//time.Sleep(time.Duration(1000000/1) * time.Microsecond)
			c.Lock()
			if c.contReader.Len() == 0 {
				c.contReader.Reset(c.cont)
			}
			c.Unlock()
		}
	}
}

func (c *Content) Read(buf []byte) (n int, err error) {
	c.Lock()
	defer c.Unlock()
	n, err = c.contReader.Read(buf)

	return n, err
}
