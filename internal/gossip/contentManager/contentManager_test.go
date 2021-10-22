package contentmanager

import (
	"bytes"
	"context"
	"reflect"
	"testing"
	"time"
)

func TestContentManager(t *testing.T) {
	cases := []struct {
		have []byte
		want []byte
	}{
		{[]byte("abraka\ndabra"), []byte("abrakadabra")},
		{[]byte("a bra kadabra"), []byte("a bra kadabra")},
	}

	for _, c := range cases {
		//var r io.Reader
		bts := make([]byte, 0, len(c.have))
		bts = append(bts, c.have...)
		r := bytes.NewReader(bts)
		content := New(r, "test")
		go content.Manager(context.Background())
		//r = bytes.NewReader(c.have)
		time.Sleep(1000 * time.Millisecond)
		bts = append(bts, c.have...)
		time.Sleep(1000 * time.Millisecond)
		got := make([]byte, len(c.have)+6)
		n, err := content.Read(got)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("Content\n have:\n %v \n  want:\n %v \n got:\n %v \n", c.have, c.want, got)
		}
		if n != len(got) {
			t.Errorf("Length content\n have: %v, want: %v, got: %v\n", len(c.have), len(c.want), n)
		}
		if err != nil {
			t.Errorf("Error: %v", err)
		}
	}
}
