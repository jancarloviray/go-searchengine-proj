package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	f "github.com/argusdusty/Ferret"
)

func engine() {
	ferret = *buildFerret()
}

func buildFerret() *f.InvertedSuffix {
	t := time.Now()
	file, err := os.Open("asv.txt")
	if err != nil {
		panic(err)
	}

	// file scanner
	scn := bufio.NewScanner(file)
	scn.Split(bufio.ScanWords)

	for scn.Scan() {
		each(scn.Text())
	}

	// STATS
	fmt.Println("Loaded document in:", time.Now().Sub(t))
	fmt.Printf("There are %v words.\n", len(Words))

	// INDEX
	t = time.Now()
	ferret := f.New(Words, Words, Values, f.UnicodeToLowerASCII)
	fmt.Println("Created index in:", time.Now().Sub(t))

	return ferret
}

func each(scannedWord string) {
	// cleanup words to be keys
	contextWord, cleanedKey := cleanKey(scannedWord)

	// an array of words
	Words = append(Words, cleanedKey)

	// current context (currently iterated word + 11 more words)
	if len(CurrentContext) > MaximumContextSize {
		// reassign, skipping the first word
		CurrentContext = CurrentContext[1:]
	}

	// append
	CurrentContext = append(CurrentContext, cleanWord(contextWord))
	context := strings.Join(CurrentContext, " ")

	// build values
	Values = append(Values, context)
}

func cleanKey(rawWord string) (before, cleanedKey string) {
	before = rawWord
	r := strings.NewReplacer("*", "", ",", "", ".", "", "^", "")
	cleanedKey = strings.ToLower(r.Replace(rawWord))
	return
}

func cleanWord(rawWord string) (cleanedWord string) {
	r := strings.NewReplacer("*", "", "^", "", "_", "")
	cleanedWord = strings.ToLower(r.Replace(rawWord))
	return
}
