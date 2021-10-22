package contentmanager

import (
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	ring := newRingBuf(10)
	cases := []struct {
		have struct {
			h int
			t int
			r [][]byte
			a []byte
		}
		want struct {
			h int
			t int
			r [][]byte
			b bool
		}
	}{
		{struct {
			h int
			t int
			r [][]byte
			a []byte
		}{5, 2, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, []byte{9}},
			struct {
				h int
				t int
				r [][]byte
				b bool
			}{6, 2, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {9}, {7}, {8}, {}}, true}},
		{struct {
			h int
			t int
			r [][]byte
			a []byte
		}{0, 0, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, []byte{}},
			struct {
				h int
				t int
				r [][]byte
				b bool
			}{1, 0, [][]byte{{0}, {}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, true}},
		{struct {
			h int
			t int
			r [][]byte
			a []byte
		}{9, 1, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, []byte{9}},
			struct {
				h int
				t int
				r [][]byte
				b bool
			}{0, 1, [][]byte{{9}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, true}},
		{struct {
			h int
			t int
			r [][]byte
			a []byte
		}{5, 5, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, []byte{0}},
			struct {
				h int
				t int
				r [][]byte
				b bool
			}{6, 5, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {0}, {7}, {8}, {}}, true}},
		{struct {
			h int
			t int
			r [][]byte
			a []byte
		}{9, 9, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, []byte{9}},
			struct {
				h int
				t int
				r [][]byte
				b bool
			}{0, 9, [][]byte{{9}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, true}},
		{struct {
			h int
			t int
			r [][]byte
			a []byte
		}{2, 5, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, []byte{99}},
			struct {
				h int
				t int
				r [][]byte
				b bool
			}{3, 5, [][]byte{{0}, {1}, {2}, {99}, {4}, {5}, {6}, {7}, {8}, {}}, true}},
		{struct {
			h int
			t int
			r [][]byte
			a []byte
		}{5, 6, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, []byte{9}},
			struct {
				h int
				t int
				r [][]byte
				b bool
			}{5, 6, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, false}},
		{struct {
			h int
			t int
			r [][]byte
			a []byte
		}{9, 0, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, []byte{9}},
			struct {
				h int
				t int
				r [][]byte
				b bool
			}{9, 0, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, false}},
	}

	for _, c := range cases {
		ring.head = c.have.h
		ring.tail = c.have.t
		ring.body = c.have.r
		b := ring.add(c.have.a)
		if (b != c.want.b) || (ring.head != c.want.h) || (ring.tail != c.want.t) || !reflect.DeepEqual(ring.body, c.want.r) {
			t.Errorf("\nHave: head = %v, tail = %v, body = %v, slice = %v; \ngot: head = %v, tail = %v, body = %v, b = %v; \nwant: head = %v, tail = %v, body = %v, b = %v", c.have.h, c.have.t, c.have.r, c.have.a, ring.head, ring.tail, ring.body, b, c.want.h, c.want.t, c.want.r, c.want.b)
		}
	}

}

func TestSub(t *testing.T) {
	ring := newRingBuf(10)
	cases := []struct {
		have struct {
			h int
			t int
			r [][]byte
		}
		want struct {
			h int
			t int
			r [][]byte
			a []byte
			b bool
		}
	}{
		{struct {
			h int
			t int
			r [][]byte
		}{5, 2, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}},
			struct {
				h int
				t int
				r [][]byte
				a []byte
				b bool
			}{5, 3, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, []byte{2}, true}},
		{struct {
			h int
			t int
			r [][]byte
		}{5, 4, [][]byte{{0}, {1}, {2}, {3}, {}, {5}, {6}, {7}, {8}, {}}},
			struct {
				h int
				t int
				r [][]byte
				a []byte
				b bool
			}{5, 5, [][]byte{{0}, {1}, {2}, {3}, {}, {5}, {6}, {7}, {8}, {}}, []byte{}, true}},
		{struct {
			h int
			t int
			r [][]byte
		}{0, 8, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {0}, {}}},
			struct {
				h int
				t int
				r [][]byte
				a []byte
				b bool
			}{0, 9, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {0}, {}}, []byte{0}, true}},
		{struct {
			h int
			t int
			r [][]byte
		}{0, 9, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}},
			struct {
				h int
				t int
				r [][]byte
				a []byte
				b bool
			}{0, 0, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, []byte{}, true}},
		{struct {
			h int
			t int
			r [][]byte
		}{0, 0, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}},
			struct {
				h int
				t int
				r [][]byte
				a []byte
				b bool
			}{0, 0, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, nil, false}},
		{struct {
			h int
			t int
			r [][]byte
		}{9, 9, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}},
			struct {
				h int
				t int
				r [][]byte
				a []byte
				b bool
			}{9, 9, [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {}}, nil, false}},
	}

	for _, c := range cases {
		ring.head = c.have.h
		ring.tail = c.have.t
		ring.body = c.have.r
		a, b := ring.sub()
		if (b != c.want.b) || (ring.head != c.want.h) || (ring.tail != c.want.t) || !reflect.DeepEqual(ring.body, c.want.r) || !reflect.DeepEqual(a, c.want.a) {
			t.Errorf("\nHave: head = %v, tail = %v, body = %v; \ngot: head = %v, tail = %v, body = %v, slice = %v, b = %v; \nwant: head = %v, tail = %v, body = %v, slice = %v, b = %v", c.have.h, c.have.t, c.have.r, ring.head, ring.tail, ring.body, a, b, c.want.h, c.want.t, c.want.r, c.want.a, c.want.b)
		}
	}

}
