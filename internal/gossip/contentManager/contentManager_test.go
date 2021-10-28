package contentmanager

import (
	"io"
	"reflect"
	"testing"
)

// =============		(queueRead)
// ^----------------	(buf)
func TestContentManager_1(t *testing.T) {
	b := true
	rb := newRingBuf(10)
	b = b && rb.add([]byte("one"))
	b = b && rb.add([]byte("two"))
	b = b && rb.add([]byte("three"))
	b = b && rb.add([]byte("four"))
	b = b && rb.add([]byte("five"))
	if !b {
		t.Errorf("Something went wrong: data initialization error.")
	}
	_, b = rb.sub()
	if !b {
		t.Errorf("Something went wrong: data initialization error.")
	}

	var cont Content
	cont.reader = nil
	cont.name = "test"
	cont.ring = *rb
	cont.markQueueRead = 0
	cont.queueRead = make([]byte, 0)

	buf := make([]byte, 8)
	buf = append(buf[:0], []byte("********")...)
	//buf = append(buf[:0], []byte("xxxxx")...)

	n, err := cont.Read(buf)

	if (n != len([]byte("two"))) || (err != io.EOF) || (!reflect.DeepEqual(buf, []byte("two*****"))) {
		t.Errorf("\nWant: %v, %v, %v\n got: %v, %v, %v\n", len([]byte("two")), io.EOF, []byte("two*****"), n, err, buf)
	}
}

// ===================		(queueRead)
// ^----------------		(buf)
func TestContentManager_2(t *testing.T) {
	b := true
	rb := newRingBuf(10)
	b = b && rb.add([]byte("one"))
	b = b && rb.add([]byte("twotwotwo"))
	b = b && rb.add([]byte("three"))
	b = b && rb.add([]byte("four"))
	b = b && rb.add([]byte("five"))
	if !b {
		t.Errorf("Something went wrong: data initialization error.")
	}
	_, b = rb.sub()
	if !b {
		t.Errorf("Something went wrong: data initialization error.")
	}

	var cont Content
	cont.reader = nil
	cont.name = "test"
	cont.ring = *rb
	cont.markQueueRead = 0
	cont.queueRead = make([]byte, 0)

	buf := make([]byte, 8)
	buf = append(buf[:0], []byte("********")...)
	//buf = append(buf[:0], []byte("xxx")...)

	n, err := cont.Read(buf)

	if (n != len([]byte("twotwotw"))) || (err != nil) || (!reflect.DeepEqual(buf, []byte("twotwotw"))) {
		t.Errorf("\nWant: %v, %v, %v\n got: %v, %v, %v\n", len([]byte("twotwotw")), nil, []byte("twotwotw"), n, err, buf)
	}
}

// ====================		(queueRead)
//    ^-----------			(buf)
func TestContentManager_3(t *testing.T) {
	b := true
	rb := newRingBuf(10)
	b = b && rb.add([]byte("one"))
	b = b && rb.add([]byte("twotwotwotwotwotwotwo"))
	b = b && rb.add([]byte("three"))
	b = b && rb.add([]byte("four"))
	b = b && rb.add([]byte("five"))
	if !b {
		t.Errorf("Something went wrong: data initialization error.")
	}
	_, b = rb.sub()
	if !b {
		t.Errorf("Something went wrong: data initialization error.")
	}

	var cont Content
	cont.reader = nil
	cont.name = "test"
	cont.ring = *rb
	cont.markQueueRead = 0
	cont.queueRead = make([]byte, 0)

	buf := make([]byte, 8)
	buf = append(buf[:0], []byte("********")...)
	//buf = append(buf[:0], []byte("xxx")...)

	n, err := cont.Read(buf)
	n, err = cont.Read(buf)

	if (n != len([]byte("otwotwot"))) || (err != nil) || (!reflect.DeepEqual(buf, []byte("otwotwot"))) {
		t.Errorf("\nWant: %v, %v, %v\n got: %v, %v, %v\n", len([]byte("otwotwot")), nil, []byte("otwotwot"), n, err, buf)
	}
}

// =================		(queueRead)
//         ^-----------		(buf)
func TestContentManager_4(t *testing.T) {
	b := true
	rb := newRingBuf(10)
	b = b && rb.add([]byte("one"))
	b = b && rb.add([]byte("01234567890123456789"))
	b = b && rb.add([]byte("three"))
	b = b && rb.add([]byte("four"))
	b = b && rb.add([]byte("five"))
	if !b {
		t.Errorf("Something went wrong: data initialization error.")
	}
	_, b = rb.sub()
	if !b {
		t.Errorf("Something went wrong: data initialization error.")
	}

	var cont Content
	cont.reader = nil
	cont.name = "test"
	cont.ring = *rb
	cont.markQueueRead = 0
	cont.queueRead = make([]byte, 0)

	buf := make([]byte, 8)
	buf = append(buf[:0], []byte("********")...)
	//buf = append(buf[:0], []byte("xxx")...)

	n, err := cont.Read(buf)
	n, err = cont.Read(buf)
	n, err = cont.Read(buf)

	if (n != len([]byte("6789"))) || (err != io.EOF) || (!reflect.DeepEqual(buf, []byte("67892345"))) {
		t.Errorf("\nWant: %v, %v, %v\n got: %v, %v, %v\n", len([]byte("6789")), io.EOF, []byte("67892345"), n, err, buf)
	}
}

// =================		(queueRead)
//         ^-----------		(buf)
//
// =========
// ^-----------
func TestContentManager_5(t *testing.T) {
	b := true
	rb := newRingBuf(10)
	b = b && rb.add([]byte("one"))
	b = b && rb.add([]byte("01234567890123456789"))
	b = b && rb.add([]byte("three"))
	b = b && rb.add([]byte("four"))
	b = b && rb.add([]byte("five"))
	if !b {
		t.Errorf("Something went wrong: data initialization error.")
	}
	_, b = rb.sub()
	if !b {
		t.Errorf("Something went wrong: data initialization error.")
	}

	var cont Content
	cont.reader = nil
	cont.name = "test"
	cont.ring = *rb
	cont.markQueueRead = 0
	cont.queueRead = make([]byte, 0)

	buf := make([]byte, 8)
	buf = append(buf[:0], []byte("********")...)
	//buf = append(buf[:0], []byte("xxx")...)

	n, err := cont.Read(buf)
	n, err = cont.Read(buf)
	n, err = cont.Read(buf)
	n, err = cont.Read(buf)

	if (n != len([]byte("three"))) || (err != io.EOF) || (!reflect.DeepEqual(buf, []byte("three345"))) {
		t.Errorf("\nWant: %v, %v, %v\n got: %v, %v, %v\n", len([]byte("three")), io.EOF, []byte("three345"), n, err, buf)
	}
}

// =================		(queueRead)
//         ^-----------		(buf)
//
// =========
// ^-----------
//
// .
// ^-----------
func TestContentManager_6(t *testing.T) {
	b := true
	rb := newRingBuf(10)
	b = b && rb.add([]byte("one"))
	b = b && rb.add([]byte("01234567890123456789"))
	b = b && rb.add([]byte("three"))
	b = b && rb.add([]byte(""))
	b = b && rb.add([]byte("five"))
	if !b {
		t.Errorf("Something went wrong: data initialization error.")
	}
	_, b = rb.sub()
	if !b {
		t.Errorf("Something went wrong: data initialization error.")
	}

	var cont Content
	cont.reader = nil
	cont.name = "test"
	cont.ring = *rb
	cont.markQueueRead = 0
	cont.queueRead = make([]byte, 0)

	buf := make([]byte, 8)
	buf = append(buf[:0], []byte("********")...)
	//buf = append(buf[:0], []byte("xxx")...)

	n, err := cont.Read(buf)
	n, err = cont.Read(buf)
	n, err = cont.Read(buf)
	n, err = cont.Read(buf)
	n, err = cont.Read(buf)

	if (n != len([]byte(""))) || (err != io.EOF) || (!reflect.DeepEqual(buf, []byte("three345"))) {
		t.Errorf("\nWant: %v, %v, %v\n got: %v, %v, %v\n", len([]byte("")), io.EOF, []byte("three345"), n, err, buf)
	}
}
