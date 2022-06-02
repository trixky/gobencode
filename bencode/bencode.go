// Package bencode provide a Bencode structure for interpret and use parsed data
package bencode

import (
	"errors"
	"math/rand"
	"time"
)

const (
	DictionaryKeyAnnounce     = "announce"
	DictionaryKeyAnnounceList = "announce-list"
	DictionaryKeyComment      = "comment"
	DictionaryKeyCreatedBy    = "created by"
	DictionaryKeyCreationDate = "creation date"
	DictionaryKeyInfo         = "info"
	DictionaryKeyName         = "name"
	DictionaryKeyLength       = "length"
	DictionaryKeyPath         = "path"
	DictionaryKeyPieceLength  = "piece length"
	DictionaryKeyPieces       = "pieces"
	DictionaryKeyUrlList      = "url-list"
	DictionaryKeyFiles        = "files"
)

var (
	ErrorNoEndpointFound                      = errors.New("no endpoint found")
	ErrorDataIsNotADictionary                 = errors.New("data is not a dictionary")
	ErrorElementMissingInDictionary           = errors.New("element missing in dictionary")
	ErrorStringElementMissingInDictionary     = errors.New("string element missing in dictionary")
	ErrorIntegerElementMissingInDictionary    = errors.New("integer element missing in dictionary")
	ErrorDictionaryElementMissingInDictionary = errors.New("dictionary element missing in dictionary")
	ErrorNeedToBeADictionaryList              = errors.New("need to be a dictionary list")
	ErrorNeedToBeAStringList                  = errors.New("need to be a string list")
	ErrorNeedToBeAList                        = errors.New("need to be a list")
	ErrorLengthIsNotMultipleOf20              = errors.New("length is not a multiple of 20")
	ErrorFilePathIsMissing                    = errors.New("file path is missing")
	ErrorDecomposedFilePathIsMissing          = errors.New("decomposedfile path is missing")
	ErrorFileNameIsMissing                    = errors.New("file name is missing")
	ErrorDirectoryNameIsMissing               = errors.New("directory name is missing")
)

type Piece [20]byte

type File struct {
	Length         int
	Path           string
	DecomposedPath []string
	CompletePath   string
}

type Info struct {
	Files         []File
	PieceLength   int
	Pieces        []Piece
	DirectoryName string
}

type Bencode struct {
	Data                   interface{}
	Announce               string
	AnnounceList           [][]string
	RandomizedAnnounceList []string
	Comment                string
	CreatedBy              string
	CreationDate           int
	Info                   Info
	InfoHash               [20]byte
	UrlList                []string
}

// RandomizeAnnounceList generates a Randomized Announce List from the initial announce list
//
// http://www.bittorrent.org/beps/bep_0012.html
// "URLs within each tier will be processed in a randomly chosen order;
// in other words, the list will be shuffled when first read, and then parsed in order"
func (b *Bencode) RandomizeAnnounceList() error {
	b.RandomizedAnnounceList = []string{}
	rand.Seed(time.Now().UnixNano())

	for _, sub_announce_list := range b.AnnounceList {
		sub_announce_list_copy := sub_announce_list

		for l := len(sub_announce_list_copy); l > 0; l-- {
			random_index := rand.Intn(l)
			b.RandomizedAnnounceList = append(b.RandomizedAnnounceList, sub_announce_list_copy[random_index])
			sub_announce_list_copy = append(sub_announce_list_copy[:random_index], sub_announce_list_copy[random_index+1:]...)
		}
	}

	return nil
}
