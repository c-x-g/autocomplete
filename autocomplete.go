package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

var fptr *string
var cache map[string]string // cache to avoid repeatedly searching for same values
var file *os.File

func main() {

	var err error

	file, err = os.Create("autocomplete.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fptr = flag.String("fpath", "shakespeare-complete.txt", "")
	flag.Parse()

	cache = make(map[string]string)

	autoCompleteHandler := http.HandlerFunc(autoComplete)
	http.Handle("/autocomplete", autoCompleteHandler)
	http.ListenAndServe(":9000", nil)
}

// https://golangbyexample.com/net-http-package-get-query-params-golang/#:~:text=Often%20in%20the%20context%20of,have%20one%20or%20multiple%20values.
func autoComplete(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	abbrev, _ := query["term"]

	badForm, _ := regexp.MatchString(`\W|\d`, abbrev[0])

	if badForm || abbrev[0] == "" {
		w.Write([]byte("The query should have length of at least one and consist of characters only\n"))
	} else {
		w.WriteHeader(200)
		if val, ok := cache[abbrev[0]]; ok {
			w.Write([]byte(val))
			writeToFile(val)
		} else {
			res := suggest(abbrev[0])
			w.Write([]byte(res))
			writeToFile(res)
		}
	}
}

func suggest(pattern string) string {

	// for reading files line by line:
	// https://www.educative.io/edpresso/file-reading-in-golang
	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	s := bufio.NewScanner(f)

	/*
		strings as keys and ints as values to keep track
		of the occurrences of each word that matches the pattern
	 */
	var m map[string]int = make(map[string]int)

	for s.Scan() {

		// split and regex in golang:
		// https://yourbasic.org/golang/split-string-into-slice/#split-on-comma-or-other-substring
		// https://yourbasic.org/golang/regexp-cheat-sheet/#character-classes

		re := regexp.MustCompile(`\W|\d`)
		t := re.Split(s.Text(), -1)

		for i := 0; i < len(t); i++ {
			lowerString := strings.ToLower(string(t[i]))
			isMatch, _ := regexp.MatchString(`^`+pattern+`.*`, lowerString)

			if isMatch {
				m[lowerString]++
			}
		}
	}
	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}

	result := sortByValues(m, pattern)

	// this accounts for situations where there are no results
	if len(result) == (len(pattern) + 2) {
		cache[pattern] = "No results for " + pattern + "\n"
		return cache[pattern]
	}

	cache[pattern] = result // store this value in the cache so in the same session it is not
							// calculated repeatedly
	return result
}
/*
    find the 25 topmost frequent words matching the pattern
	by sorting the map by its values
    https://ispycode.com/GO/Sorting/Sort-map-by-value
*/
func sortByValues(m map[string]int, pattern string) string {

	var reverseMap map[int][]string = make(map[int][]string)
	var ints []int

	result := pattern + ":"
	for k, _ := range m {
		reverseMap[m[k]] = append(reverseMap[m[k]], k)
		ints = append(ints, m[k])
	}

	// Sort in descending order
	// https://stackoverflow.com/questions/37695209/golang-sort-slice-ascending-or-descending/40932847
	sort.Slice(ints, func(i, j int) bool {
		return ints[i] > ints[j]
	})

	var count int
	count = 25

	var curr int // keep track of current value in ints and skip if it is duplicated
	curr = -1

	top25:
	for i := 0; i < len(ints); i++ {

		if ints[i] == curr { //ignore duplicated values in ints
			continue
		}
		curr = ints[i]

		for j := 0; j < len(reverseMap[ints[i]]); j++ {
			result = result + " " + reverseMap[ints[i]][j] // concatenate the top 25 most frequent results
			count = count - 1
			if count == 0 {
				break top25
			}
		}
	}

	return result + "\n"
}

// write result to autocomplete.txt
func writeToFile(value string) {
	_, err := io.WriteString(file, value)
	if err != nil {
		log.Fatal(err)
	}
	file.Sync()
}
