package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var input io.Reader

	switch len(os.Args) {
	case 1:
		input = os.Stdin
	case 2:
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		input = f
	default:
		log.Fatal("This program expects either 0 or 1 arguments.")
	}

	r := csv.NewReader(input)
	r.FieldsPerRecord = -1

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Remove the very first "record" (i.e 'Category: All categories') if exists
	if len(records[0]) == 1 {
		records = append(records[:0], records[1:]...)
	}

	// Save names to a slice
	names := records[0][1:] // Skip 'weeks' column
	commonSuffix := longestCommonSuffix(names)
	if commonSuffix != "" {
		for i, name := range names {
			names[i] = strings.TrimSuffix(name, commonSuffix)
		}
	}
	records = append(records[:0], records[1:]...)

	avg := make([]int, len(names))

	for _, record := range records {
		record = record[1:] // Skip 'weeks' column
		for i, s := range record {
			n, err := strconv.Atoi(s)
			if err != nil {
				log.Fatal(err)
			}
			avg[i] += n
		}
	}

	for i := 0; i < len(avg); i++ {
		avg[i] = avg[i] / len(records)
	}

	for i, n := range avg {
		n = n/10 + 1
		fmt.Printf("%s %s (%d)\n", strings.Repeat("▓", n), names[i], avg[i])
	}
}

func longestCommonSuffix(a []string) string {
	if len(a) == 0 {
		return ""
	}

	suffix := a[0]
	if len(a) == 1 {
		return suffix
	}

	for _, s := range a[1:] {
		suffixLength := len(suffix)
		sLength := len(s)

		if suffixLength == 0 || sLength == 0 {
			return ""
		}

		maxLength := suffixLength
		if sLength < maxLength {
			maxLength = sLength
		}

		for i := 0; i < maxLength; i++ {
			j := suffixLength - i - 1
			k := sLength - i - 1
			if suffix[j] != s[k] {
				suffix = suffix[j+1:]
				break
			}
		}
	}
	return suffix
}
