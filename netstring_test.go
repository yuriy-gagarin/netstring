package netstring_test

import (
	"bufio"
	"bytes"
	"io"
	"testing"

	"github.com/yuriy-gagarin/netstring"
)

type bb [][]byte
type b []byte

func TestNetstring(t *testing.T) {
	tcs := []struct {
		input    b
		expected bb
	}{
		{
			b(`10:helloworld,7:1234567,`),
			bb{b("helloworld"), b("1234567")},
		},
		{
			b(`10:helloworld,7:1`),
			bb{b("helloworld")},
		},
		{
			b(`10helloworld,7:1`),
			bb{},
		},
		{
			b(`oworld,7:1234567,`),
			bb{b("1234567")},
		},
		{
			b(`8:zawarudo,,,fffffff,,7:1234567,`),
			bb{b("zawarudo"), b("1234567")},
		},
	}

	for _, tc := range tcs {

		scn := bufio.NewScanner(bytes.NewBuffer(tc.input))
		scn.Split(netstring.SplitNetstring)

		i := 0

		for i < len(tc.expected) {
			scn.Scan()
			s := scn.Bytes()
			diff := bytes.Compare(tc.expected[i], s)
			if diff != 0 {
				t.Errorf("expected %v got %v, possible error %v", tc.expected[i], s, scn.Err())
			}
			i++
		}
	}
}

func TestNetstringChunks(t *testing.T) {

	input := bb{b(`paddingpaddin`), b(`10:hello`), b(`world,,,,,,,7:12`), b(`34567,`)}
	expected := bb{b("helloworld"), b("1234567")}

	pr, pw := io.Pipe()

	scn := bufio.NewScanner(pr)
	scn.Split(netstring.SplitNetstring)

	i := 0

	go func() {
		for _, v := range input {
			pw.Write(v)
		}
	}()

	for i < len(expected) {
		scn.Scan()
		s := scn.Bytes()
		diff := bytes.Compare(expected[i], s)
		if diff != 0 {
			t.Errorf("expected %v got %v, possible error %v", expected[i], s, scn.Err())
		}
		i++
	}

}
