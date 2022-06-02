package bencode

import (
	"bufio"
	"os"
	"testing"

	"github.com/trixky/gobencode/parser"
)

func TestGetInfoHash(t *testing.T) {
	tests := []struct {
		file     string
		expected [20]byte
	}{
		{
			file:     "../.test_files/arch.torrent",
			expected: [20]byte{155, 76, 20, 137, 191, 204, 216, 32, 93, 21, 35, 69, 247, 168, 170, 213, 45, 154, 31, 87},
		},
		{
			file:     "../.test_files/kubuntu.torrent",
			expected: [20]byte{62, 206, 208, 101, 4, 33, 21, 130, 87, 37, 224, 129, 200, 16, 132, 218, 41, 250, 120, 2},
		},
		{
			file:     "../.test_files/ubuntu.torrent",
			expected: [20]byte{44, 107, 104, 88, 214, 29, 169, 84, 61, 66, 49, 167, 29, 180, 177, 201, 38, 75, 6, 133},
		},
		{
			file:     "../.test_files/ubuntu.torrent",
			expected: [20]byte{44, 107, 104, 88, 214, 29, 169, 84, 61, 66, 49, 167, 29, 180, 177, 201, 38, 75, 6, 133},
		},
	}

	for index, test := range tests {
		f, err := os.Open(test.file)

		if err != nil {
			t.Errorf("failed to read file %d: %v", index, err)
			continue
		}

		data, err := parser.ParseElement(bufio.NewReader(f))

		if err != nil {
			t.Errorf("failed to get data of file %d: %v", index, err)
			continue
		}

		bc := Bencode{
			Data: data,
		}

		err = bc.UnmarshallAll()

		if err != nil {
			t.Errorf("failed to unmarshall file %d: %v", index, err)
			continue
		}

		err = bc.GetInfoHash()

		if err != nil {
			t.Errorf("failed to get info hash of file %d: %v", index, err)
			continue
		}

		if test.expected != bc.InfoHash {
			t.Errorf("expected %v | %v output", test.expected, bc.InfoHash)
			continue
		}
	}
}
