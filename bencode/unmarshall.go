package bencode

import (
	"fmt"
	"strings"

	"github.com/trixky/gobencode/utils"
)

// unmarshallPath unmarshall the Path attribute from a bencode
func (f *File) unmarshallPath(data interface{}) error {
	if interface_list, ok := data.([]interface{}); ok {
		paths := []string{}

		for _, interface_path := range interface_list {
			if path, ok := interface_path.(string); ok {
				paths = append(paths, path)
			} else {
				return ErrorNeedToBeAStringList
			}
		}

		f.DecomposedPath = paths
		f.Path = strings.Join(paths, "/")
	} else {
		return ErrorNeedToBeAList
	}

	return nil
}

// unmarshallPieceLength unmarshall the Piece Length attribute from a bencode info section
func (i *Info) unmarshallPieceLength(info_dictionary map[string]interface{}) error {
	if piece_length, ok := info_dictionary[DictionaryKeyPieceLength].(int); ok {
		i.PieceLength = piece_length
		return nil
	}

	return fmt.Errorf("%w: %v", ErrorIntegerElementMissingInDictionary, DictionaryKeyPieceLength)
}

// unmarshallName unmarshall the Name attribute from a bencode info section
func (i *Info) unmarshallName(info_dictionary map[string]interface{}) error {
	if name, ok := info_dictionary[DictionaryKeyName].(string); ok {
		i.DirectoryName = name
		return nil
	}

	return fmt.Errorf("%: %v", ErrorStringElementMissingInDictionary, DictionaryKeyName)
}

// unmarshallPieces unmarshall the Pieces attribute from a bencode info section
func (i *Info) unmarshallPieces(info_dictionary map[string]interface{}) error {
	if pieces_bytes, ok := info_dictionary[DictionaryKeyPieces].(string); ok {
		len_pieces_bytes := len(pieces_bytes) // need to be a multiple of 20

		if len_pieces_bytes%20 != 0 {
			return fmt.Errorf("%w: %v", ErrorLengthIsNotMultipleOf20, DictionaryKeyPieces)
		}

		len_pieces := len_pieces_bytes / 20

		pieces := make([]Piece, len_pieces)

		for i := 0; i < len_pieces; i++ {
			start := i * 20
			end := start + 20
			copy(pieces[i][:], ([]byte(pieces_bytes)[start:end]))
		}

		i.Pieces = pieces

		return nil
	}

	return fmt.Errorf("%v: %v", ErrorIntegerElementMissingInDictionary, DictionaryKeyPieces)
}

// unmarshallFiles unmarshall the Files attribute from a bencode info section
func (i *Info) unmarshallFiles(info_dictionary map[string]interface{}) error {
	info_files, ok := info_dictionary[DictionaryKeyFiles]

	if ok {
		if files_list, ok := info_files.([]interface{}); ok {
			for _, file := range files_list {
				if file_dictionary, ok := file.(map[string]interface{}); ok {
					if file_length, ok := file_dictionary[DictionaryKeyLength].(int); ok {
						if file_path, ok := file_dictionary[DictionaryKeyPath]; ok {
							file := File{
								Length: file_length,
							}

							if err := file.unmarshallPath(file_path); err != nil {
								return fmt.Errorf("file corrupted: %w", err)
							}

							if len(i.DirectoryName) > 0 {
								file.CompletePath = i.DirectoryName + "/" + file.Path
							}

							i.Files = append(i.Files, file)
						} else {
							return fmt.Errorf("file corrupted: %w (%s)", ErrorElementMissingInDictionary, DictionaryKeyPath)
						}
					} else {
						return fmt.Errorf("file corrupted: %w (%s)", ErrorIntegerElementMissingInDictionary, DictionaryKeyLength)
					}
				} else {
					return fmt.Errorf("file corrupted: %w", ErrorNeedToBeADictionaryList)
				}
			}
		} else {
			return ErrorNeedToBeAList
		}
	} else {
		if file_length, ok := info_dictionary[DictionaryKeyLength].(int); ok {
			i.Files = append(i.Files, File{
				Length:       file_length,
				Path:         i.DirectoryName,
				CompletePath: i.DirectoryName,
			})
		} else {
			return fmt.Errorf("%w: %v", ErrorIntegerElementMissingInDictionary, DictionaryKeyLength)
		}
	}

	return nil
}

// unmarshallStringElement unmarshall a string from a specific dictionary key
func (b *Bencode) unmarshallStringElement(key string) (string, error) {
	dictionary, ok := b.Data.(map[string]interface{})

	if !ok {
		return "", ErrorDataIsNotADictionary

	}

	value, ok := dictionary[key].(string)

	if !ok {
		return "", fmt.Errorf("%w", ErrorStringElementMissingInDictionary)
	}

	return value, nil
}

// unmarshallIntegerElement unmarshall a integer from a specific dictionary key
func (b *Bencode) unmarshallIntegerElement(key string) (int, error) {
	dictionary, ok := b.Data.(map[string]interface{})

	if !ok {
		return 0, ErrorDataIsNotADictionary
	}

	value, ok := dictionary[key].(int)

	if !ok {
		return 0, fmt.Errorf("%w", ErrorIntegerElementMissingInDictionary)
	}

	return value, nil
}

// UnmarshallAnnounce unmarshall the Announce attribute
func (b *Bencode) UnmarshallAnnounce() error {
	value, err := b.unmarshallStringElement(DictionaryKeyAnnounce)

	if err != nil {
		return err
	}

	b.Announce = value

	return nil
}

// UnmarshallAnnounceList unmarshall the Announce List attribute
func (b *Bencode) UnmarshallAnnounceList() error {
	dictionary, ok := b.Data.(map[string]interface{})

	if !ok {
		return ErrorDataIsNotADictionary
	}

	value, ok := dictionary[DictionaryKeyAnnounceList]

	if !ok {
		return ErrorElementMissingInDictionary
	}

	announce_list, err := utils.ToListOfStringList(value)

	if err != nil {
		return err
	}

	b.AnnounceList = announce_list

	return nil
}

// UnmarshallComment unmarshall the Comment attribute
func (b *Bencode) UnmarshallComment() error {
	value, err := b.unmarshallStringElement(DictionaryKeyComment)

	if err != nil {
		return err
	}

	b.Comment = value

	return nil
}

// UnmarshallCreatedBy unmarshall the Created By attribute
func (b *Bencode) UnmarshallCreatedBy() error {
	value, err := b.unmarshallStringElement(DictionaryKeyCreatedBy)

	if err != nil {
		return err
	}

	b.CreatedBy = value

	return nil
}

// UnmarshallCreationDate unmarshall the Creation Date attribute
func (b *Bencode) UnmarshallCreationDate() error {
	value, err := b.unmarshallIntegerElement(DictionaryKeyCreationDate)

	if err != nil {
		return err
	}

	b.CreationDate = value

	return nil
}

// UnmarshallInfo unmarshall the Info section attribute
func (b *Bencode) UnmarshallInfo() error {

	info := Info{}

	dictionary, ok := b.Data.(map[string]interface{})

	if !ok {
		return ErrorDataIsNotADictionary
	}

	info_dictionary, ok := dictionary[DictionaryKeyInfo].(map[string]interface{})

	if !ok {
		return ErrorDictionaryElementMissingInDictionary
	}

	// ---------- piece length
	if err := info.unmarshallPieceLength(info_dictionary); err != nil {
		return err
	}

	// ---------- pieces
	if err := info.unmarshallPieces(info_dictionary); err != nil {
		return err
	}

	// ---------- name
	if err := info.unmarshallName(info_dictionary); err != nil {
		return err
	}

	// ---------- files
	if err := info.unmarshallFiles(info_dictionary); err != nil {
		return err
	}

	b.Info = info

	return nil
}

// UnmarshallUrlList unmarshall the Url List attribute
func (b *Bencode) UnmarshallUrlList() error {
	dictionary, ok := b.Data.(map[string]interface{})

	if !ok {
		return ErrorDataIsNotADictionary
	}

	value, ok := dictionary[DictionaryKeyUrlList]

	if !ok {
		return ErrorElementMissingInDictionary
	}

	url_list, err := utils.ToStringList(value)

	if err != nil {
		return err
	}

	b.UrlList = url_list

	return nil
}

// UnmarshallAll unmarshall all attribute
func (b *Bencode) UnmarshallAll() (err error) {
	endpoint_errors := []error{}

	// ---------- endpoints
	if err := b.UnmarshallAnnounce(); err != nil {
		endpoint_errors = append(endpoint_errors, err)
	}
	if err := b.UnmarshallAnnounceList(); err != nil {
		endpoint_errors = append(endpoint_errors, err)
	}
	if err := b.UnmarshallUrlList(); err != nil {
		endpoint_errors = append(endpoint_errors, err)
	}

	if len(endpoint_errors) == 3 {
		return fmt.Errorf("%w: %v", ErrorNoEndpointFound, endpoint_errors)
	}

	if err := b.RandomizeAnnounceList(); err != nil {
		return err
	}

	// ---------- meta
	b.UnmarshallComment()
	b.UnmarshallCreatedBy()
	b.UnmarshallCreationDate()

	// ---------- info/files
	if err := b.UnmarshallInfo(); err != nil {
		return fmt.Errorf("failed to unmarshall info: %w", err)
	}

	if err := b.GetInfoHash(); err != nil {
		return fmt.Errorf("failed to get the info hash: %w", err)
	}

	return nil
}
