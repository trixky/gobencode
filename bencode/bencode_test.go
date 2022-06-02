package bencode

import (
	"fmt"
	"testing"
)

func TestRandomizeAnnounceList(t *testing.T) {
	const test_retries = 10 // used for test different random combinaisons

	tests := []struct {
		first_chars  []string
		second_chars []string
	}{
		{
			first_chars:  []string{},
			second_chars: []string{},
		},
		{
			first_chars:  []string{},
			second_chars: []string{"1", "2", "3"},
		},
		{
			first_chars:  []string{"A", "B", "C"},
			second_chars: []string{},
		},
		{
			first_chars:  []string{"A", "B", "C"},
			second_chars: []string{"1"},
		},
		{
			first_chars:  []string{"A"},
			second_chars: []string{"1", "2", "3"},
		},
		{
			first_chars:  []string{"A", "B", "C"},
			second_chars: []string{"1", "2", "3"},
		},
		{
			first_chars:  []string{"A", "B", "C", "D", "E", "F"},
			second_chars: []string{"1", "2", "3", "4", "5", "6"},
		},
	}

	testAnnounceListIntegrity := func(announce_list []string, first_char string, second_chars []string) error {
		for i := 0; i < len(announce_list); i++ {
			element := announce_list[i]
			first_character := element[:1]
			second_character := element[1:2]

			if len(element) != 2 || !func() bool {
				for _, second_char := range second_chars {
					if second_char == second_character {
						return true
					}
				}

				return false
			}() {
				return fmt.Errorf("sub list element is corrupted: [%s]", announce_list[i])
			}
			if first_character != first_char {
				return fmt.Errorf("sub list are mixed: expected [%s] | [%s] output", first_character, first_char)
			}
			for j := i + 1; j < len(announce_list); j++ {
				if element == announce_list[j] {
					return fmt.Errorf("duplication in sub list: [%s]", element)
				}
			}
		}

		return nil
	}

	generate_announce_list := func(first_chars []string, second_chars []string) (announce_list [][]string) {
		for _, first_char := range first_chars {
			sub_announce_list := []string{}
			for _, second_char := range second_chars {
				sub_announce_list = append(sub_announce_list, first_char+second_char)
			}

			announce_list = append(announce_list, sub_announce_list)
		}

		return
	}

	check_sub_randomized_announce_list := func(randomized_announce_list []string, first_chars []string, second_chars []string) {
		for index, first_char := range first_chars {
			len_second_char := len(second_chars)
			start := index * len_second_char
			end := start + len_second_char
			if err := testAnnounceListIntegrity(randomized_announce_list[start:end], first_char, second_chars); err != nil {
				t.Fatal(err)
			}
		}
	}

	testAnnounceList := func(first_chars []string, second_chars []string) {
		bencode := Bencode{
			AnnounceList: generate_announce_list(first_chars, second_chars),
		}

		bencode.RandomizeAnnounceList()

		check_sub_randomized_announce_list(bencode.RandomizedAnnounceList, first_chars, second_chars)
	}

	for _, test := range tests {
		for i := 0; i < test_retries; i++ {
			testAnnounceList(test.first_chars, test.second_chars)
		}
	}
}
