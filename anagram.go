package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

func readDict(path string, wordc chan string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		// hack to encode Latin1 to UTF8
		bytes := scan.Bytes()
		runes := make([]rune, len(bytes))
		for i, b := range bytes {
			runes[i] = rune(b)
		}
		wordc <- string(runes)
	}
	close(wordc)
	return nil
}

func getThreadsNum() int {
	nthreads := runtime.NumCPU()/2 - 1
	if nthreads < 1 {
		return 1
	}
	return nthreads
}

func sortedRunes(w string) []rune {
	runes := []rune(strings.ToLower(w))
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return runes
}

func anagrams(path, word string) ([]string, error) {
	var (
		wordc    = make(chan string, 1024)
		srunes   = sortedRunes(word)
		nthreads = getThreadsNum()
		wg       sync.WaitGroup

		mu       sync.Mutex
		anagrams []string
	)
	for n := 0; n < nthreads; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for word := range wordc {
				if string(srunes) == string(sortedRunes(word)) {
					mu.Lock()
					anagrams = append(anagrams, word)
					mu.Unlock()
				}

			}
		}()
	}

	err := readDict(path, wordc)
	if err != nil {
		return nil, err
	}

	wg.Wait()
	return anagrams, nil
}

func main() {
	start := time.Now()
	if len(os.Args) < 3 {
		fmt.Println("Usage: " + os.Args[0] + " <dictionary_path> <word>")
		return
	}
	dictPath := os.Args[1]
	word := strings.Join(os.Args[2:], " ")

	ans, err := anagrams(dictPath, word)
	if err != nil {
		fmt.Println("Failed to read a dictionary:", err)
		return
	}
	fmt.Println(time.Since(start), strings.Join(ans, ", "))
}
