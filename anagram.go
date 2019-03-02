package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"
	"unicode"
)

type Rep struct {
	Char rune
	N    int8
}

type Dict struct {
	Index map[Rep]*Dict
	Words []string
}

func findAnagrams(d *Dict, reps []Rep) []string {
	if len(reps) == 0 {
		return d.Words
	}
	if d, ok := d.Index[reps[0]]; ok {
		return findAnagrams(d, reps[1:])
	}
	return nil
}

func parseWord(word string) []Rep {
	chars := []rune(word)
	sort.Slice(chars, func(i, j int) bool {
		return chars[i] < chars[j]
	})
	var reps []Rep
	for _, c := range chars {
		c = unicode.ToLower(c)
		if len(reps) > 0 && reps[len(reps)-1].Char == c {
			reps[len(reps)-1].N++
		} else {
			reps = append(reps, Rep{Char: c, N: 1})
		}
	}
	return reps
}

func insertWord(d *Dict, reps []Rep, w string) {
	if len(reps) == 0 {
		d.Words = append(d.Words, w)
	} else {
		var newd *Dict
		if newd = d.Index[reps[0]]; newd == nil {
			newd = &Dict{Index: make(map[Rep]*Dict)}
			d.Index[reps[0]] = newd
		}
		insertWord(newd, reps[1:], w)
	}
}

func indexWord(d *Dict, w string) {
	insertWord(d, parseWord(w), w)
}

func parseDict(words []string) Dict {
	d := Dict{Index: make(map[Rep]*Dict)}
	for _, w := range words {
		indexWord(&d, w)
	}
	return d
}

func readDict(path string) ([]string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(bytes.NewReader(b))
	var words []string
	for scan.Scan() {
		// hack to encode Latin1 to UTF8
		bytes := scan.Bytes()
		runes := make([]rune, len(bytes))
		for i, b := range bytes {
			runes[i] = rune(b)
		}
		words = append(words, string(runes))
	}
	return words, nil
}

func main() {
	start := time.Now()
	if len(os.Args) < 3 {
		fmt.Println("Usage: " + os.Args[0] + " <dictionary_path> <word>")
		return
	}
	dictPath := os.Args[1]
	word := strings.Join(os.Args[2:], " ")
	words, err := readDict(dictPath)
	if err != nil {
		fmt.Println("Failed to read a dictionary:", err)
		return
	}

	reps := parseWord(word)
	var anagrams []string
	for _, w := range words {
		if reflect.DeepEqual(reps, parseWord(w)) {
			anagrams = append(anagrams, w)
		}
	}
	// d := parseDict(words)
	// anagrams := findAnagrams(&d, parseWord(word))
	fmt.Println(time.Since(start), strings.Join(anagrams, ", "))
}
