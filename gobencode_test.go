package gobencode

import (
	"os"
	"testing"
)

func TestUnmarshallFromReader(t *testing.T) {

	tests_file := []string{
		"./.test_files/arch.torrent",
		"./.test_files/kubuntu.torrent",
		"./.test_files/minecraft.torrent",
		"./.test_files/ubuntu.torrent",
	}

	for _, test := range tests_file {
		r, err := os.Open(test)

		if err != nil {
			t.Fatalf("failed to open file [%s]: %v", test, err)
			continue
		}

		bc, err := UnmarshallFromReader(r)

		if err != nil {
			t.Fatalf("failed to parse file [%s]: %v", test, err)
			continue
		}

		if err = bc.UnmarshallAll(); err != nil {
			t.Errorf("failed to unmarshall file [%s]: %v\n", test, err)
			continue
		}

		err = bc.GetInfoHash()

		if err != nil {
			t.Errorf("failed to get info hash of file [%s...]: %v\n", test[:50], err)
			continue
		}
	}
}
