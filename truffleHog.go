package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const BASE64_CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
const HEX_CHARS = "1234567890abcdefABCDEF"
const threshold = 10

const b64_minimum = 4.5
const hex_minimum = 3.0

var verbose bool = true
var curfile string = ""

func entropyfind(folder string) {
	err := filepath.Walk(folder, func(path1 string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		if isBinaryFile(path1) {
			log.Printf("INFO: %s is a bin file.\n", path1)
			return nil
		}
		if isGitFile(path1) {
			log.Printf("INFO: %s is a git file.\n", path1)
			return nil
		}

		find_entropy(path1)
		//allres += res
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

// https://stackoverflow.com/questions/39862613/how-to-split-a-string-by-multiple-delimiters
func SplitBlank(r rune) bool {
	return r == ' ' || r == '\t'
}

func find_entropy(filename string) {
	//strings_found := []string{} // store detected secrets
	//line_counter := 0           // record what line a secret was found on
	curfile = filename

	lines, err := ReadLines(filename)
	if err != nil {
		log.Fatal(err)
		return
	}

	for i, line := range lines {
		parseTextLine(i+1, line)
	}
}

func parseTextLine(row int, line string) int {
	rescount := 0
	allwords := strings.FieldsFunc(line, SplitBlank)

	for _, wd := range allwords {
		base64_strings := get_strings_of_set(wd, BASE64_CHARS) // get b64 blobs
		hex_strings := get_strings_of_set(wd, HEX_CHARS)       // get hex blobs

		for _, bs := range base64_strings {
			b64_entropy := shannon_entropy(bs, BASE64_CHARS)

			if b64_entropy > b64_minimum {
				//strings_found = append(strings_found, bs)
				res := ""
				if verbose {
					res = fmt.Sprintf("\n--------\nFile: %s\nLine: %d\nType: BASE64"+
						"\nShannon Entropy: %v\nSecret: %s\nFull Line:\n\t%s\n",
						curfile, row, b64_entropy, bs, line)

				} else {
					res = fmt.Sprintf("\nFile: %s:%d:Line:\n\t%s\n",
						curfile, row, line)
				}
				rescount++

				fmt.Println(res)
			}
		}
		for _, bs := range hex_strings {
			hex_entropy := shannon_entropy(bs, HEX_CHARS)

			if hex_entropy > hex_minimum {
				//strings_found = append(strings_found, bs)
				res := ""
				if verbose {
					res = fmt.Sprintf("\n--------\nFile: %s\nLine: %d\nType: HEX"+
						"\nShannon Entropy: %v\nSecret: %s\nFull Line:\n\t%s\n",
						filepath.Base(curfile), row, hex_entropy, bs, line)

				} else {
					res = fmt.Sprintf("\nFile: %s:%d:Line:\n\t%s\n",
						filepath.Base(curfile), row, line)
				}

				fmt.Println(res)
				rescount++
			}
		}
	}

	return rescount
}

func shannon_entropy(data, baseline string) float64 {
	bsrs := []rune(baseline)
	entropy := 0.0

	for _, r := range bsrs {
		p_x := float64(strings.Count(data, string(r))) / float64(len(data))
		if p_x > 0.0 {
			entropy += -p_x * math.Log2(p_x)
		}
	}
	return entropy
}

// return all strings in word with length > threshold that
// contain only characters in char_set (from trufflehog)
func get_strings_of_set(oneword, char_set string) []string {
	count := 0
	letters := ""
	allstrings := []string{}

	rs := []rune(oneword)
	for _, r := range rs {
		if strings.ContainsRune(char_set, r) {
			letters += string(r)
			count += 1
		} else {
			if count > threshold {
				allstrings = append(allstrings, letters)
			}
			letters = ""
			count = 0
		}
	}

	if count > threshold {
		allstrings = append(allstrings, letters)

	}

	return allstrings
}

const HeadNBytes = 10

// check head n bytes
func isBinaryFile(filename string) bool {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("WARN: Cannot read file: %s!\n", filename)
		return false
	}

	invisiblecount := 0
	if len(bs) < HeadNBytes {
		return false
	}

	for i := 0; i < HeadNBytes; i++ {
		if isBlankChar(bs[i]) {
			// tab
			continue
		}
		if (bs[i] < 32) || (bs[i] >= 127) {
			invisiblecount++
		}
	}

	if invisiblecount > 3 {
		return true
	}
	return false
}

// .gitignore
func isGitFile(filename string) bool {
	return strings.HasPrefix(filepath.Ext(filename), ".git")
}

// .clang-format
func isClangFile(filename string) bool {
	return strings.HasPrefix(filepath.Ext(filename), ".clang")
}

func isBlankChar(b byte) bool {
	switch b {
	case 0x09, 0x0d, 0x0a: // tab, return, newline
		return true
	default:
		return false
	}
}
