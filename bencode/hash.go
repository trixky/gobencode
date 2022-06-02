package bencode

import (
	"crypto/sha1"
)

// GetInfoHash encodes the info section an generate his hash
func (b *Bencode) GetInfoHash() error {
	encoded_info, err := encodeElement(b.Info)

	if err != nil {
		return err
	}

	b.InfoHash = sha1.Sum([]byte(encoded_info))

	return nil
}
