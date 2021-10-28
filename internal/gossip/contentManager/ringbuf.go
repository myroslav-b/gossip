package contentmanager

type ringBuf struct {
	power int
	body  [][]byte
	head  int
	tail  int
}

func newRingBuf(power int) *ringBuf {
	var ring ringBuf
	ring.body = make([][]byte, power)
	for i := 0; i < power; i++ {
		ring.body[i] = make([]byte, 0)
	}
	ring.power = power
	ring.head = 0
	ring.tail = 0
	return &ring
}

func (ring *ringBuf) add(b []byte) bool {
	if (ring.head+1)%ring.power != ring.tail {
		//ring.body[ring.head] = ring.body[ring.head][:0]
		ring.body[ring.head] = append(ring.body[ring.head][:0], b...)
		ring.head = (ring.head + 1) % ring.power
		return true
	} else {
		return false
	}
}

func (ring *ringBuf) sub() ([]byte, bool) {
	if ring.head != ring.tail {
		b := make([]byte, 0, len(ring.body[ring.tail]))
		b = append(b[:0], ring.body[ring.tail]...)
		ring.tail = (ring.tail + 1) % ring.power
		return b, true
	} else {
		return nil, false
	}
}
