package bencode

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
)

var (
	ErrorTypeNotEncodable = errors.New("type is not encodable")
)

// encodeString encodes a string in the bencode format
func encodeString(str string) string {
	return strconv.Itoa(len(str)) + ":" + str
}

// encodeInteger encodes an integer in the bencode format
func encodeInteger(i int) string {
	return "i" + strconv.Itoa(i) + "e"
}

// encodeList encodes a list in the bencode format
func encodeList(list []interface{}) (string, error) {
	encoded_list := "l"

	for _, element := range list {
		if encoded_element, err := encodeElement(element); err != nil {
			return "", err
		} else {
			encoded_list += encoded_element
		}
	}

	return encoded_list + "e", nil
}

// encodeDictionary encodes a dictionary in the bencode format
func encodeDictionary(dictionary map[string]interface{}) (string, error) {
	encoded_dictionary := "d"

	keys := []string{}

	for key := range dictionary {
		keys = append(keys, key)
	}

	// http://www.bittorrent.org/beps/bep_0003.html
	// "Keys must be strings and appear in sorted order (sorted as raw strings, not alphanumerics)""
	sort.Strings(keys)

	for _, key := range keys {
		encoded_dictionary += encodeString(key)

		if encoded_element, err := encodeElement(dictionary[key]); err != nil {
			return "", err
		} else {
			encoded_dictionary += encoded_element
		}
	}

	return encoded_dictionary + "e", nil
}

// encodePieces encodes Pieces in the bencode format
func encodePieces(pieces []Piece) string {
	encoded_pieces := strconv.Itoa(len(pieces)*20) + ":"

	for _, piece := range pieces {
		encoded_pieces += string(piece[:])
	}

	return encoded_pieces
}

// encodeFiles encodes Files in the bencode format
func encodeFiles(files []File) (string, error) {
	encodedFiles := encodeString(DictionaryKeyFiles) + "l"
	for _, file := range files {
		if len(file.DecomposedPath) == 0 {
			return "", ErrorFilePathIsMissing
		}

		encodedFiles += "d" + encodeString(DictionaryKeyLength) + encodeInteger(file.Length)
		encodedFiles += encodeString(DictionaryKeyPath) + "l"

		for _, path := range file.DecomposedPath {
			if len(path) == 0 {
				return "", ErrorFilePathIsMissing

			}
			encodedFiles += encodeString(path)
		}
		encodedFiles += "ee"
	}

	return encodedFiles + "e", nil
}

// encodeInfo encodes an Info section in the bencode format
func encodeInfo(info Info) (string, error) {
	encoded_info := "d"

	if len(info.Files) == 1 {
		if len(info.DirectoryName) == 0 {
			return "", ErrorFileNameIsMissing
		}

		encoded_info += encodeString(DictionaryKeyLength) + encodeInteger(info.Files[0].Length)
		encoded_info += encodeString(DictionaryKeyName) + encodeString(info.Files[0].Path)
	} else if len(info.Files) > 2 {
		if len(info.DirectoryName) == 0 {
			return "", ErrorDirectoryNameIsMissing
		}

		encodedFiles, err := encodeFiles(info.Files)

		if err != nil {
			return "", err
		}

		encoded_info += encodedFiles
		encoded_info += encodeString(DictionaryKeyName) + encodeString(info.DirectoryName)
	}

	encoded_info += encodeString(DictionaryKeyPieceLength) + encodeInteger(info.PieceLength)
	encoded_info += encodeString(DictionaryKeyPieces) + encodePieces(info.Pieces)

	return encoded_info + "e", nil
}

// encodeElement encodes any type of element in the bencode format
func encodeElement(element interface{}) (string, error) {
	switch element.(type) {
	case string:
		return encodeString(element.(string)), nil
	case int:
		return encodeInteger(element.(int)), nil
	case []interface{}:
		return encodeList(element.([]interface{}))
	case map[string]interface{}:
		return encodeDictionary(element.(map[string]interface{}))
	case []Piece:
		return encodePieces(element.([]Piece)), nil
	case Info:
		return encodeInfo(element.(Info))
	default:
		return "", fmt.Errorf("%w: %T", ErrorTypeNotEncodable, element)
	}
}
