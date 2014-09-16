package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/argusdusty/Ferret"
	g "github.com/zenazn/goji"
	web "github.com/zenazn/goji/web"
)

var public string = "./public"
var engine ferret.InvertedSuffix

func main() {
	SearchEngine()
	routes()
	g.Serve()
}

func routes() {
	g.Get("/", http.FileServer(http.Dir(public)))

	static := web.New()
	static.Get("/scripts/*", http.FileServer(http.Dir(public)))
	static.Get("/styles/*", http.FileServer(http.Dir(public)))
	static.Get("/img/*", http.FileServer(http.Dir(public)))

	api := web.New()
	api.Get("/api/search", SearchHandler)

	g.Handle("/scripts/*", static)
	g.Handle("/styles/*", static)
	g.Handle("/img/*", static)
	g.Handle("/api/*", api)
}

type SearchResponse struct {
	//Query    string        `json:"query"`
	Duration string `json:"duration"`
	//Results  []string      `json:"results"`
	Values []interface{} `json:"values"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query()["s"][0]
	t := time.Now()
	_, values := engine.Query(s, 15)
	duration := time.Now().Sub(t).String()
	data, _ := json.Marshal(SearchResponse{duration, values})
	fmt.Fprint(w, string(data))
}

func SearchEngine() {
	engine = *buildFerret()
}

func buildFerret() *ferret.InvertedSuffix {
	t := time.Now()
	file, err := os.Open("asv.txt")
	if err != nil {
		panic(err)
	}

	// file scanner
	scn := bufio.NewScanner(file)
	scn.Split(bufio.ScanWords)

	for scn.Scan() {
		eachWord(scn.Text())
	}

	// STATS
	fmt.Println("Loaded document in:", time.Now().Sub(t))
	fmt.Printf("There are %v words.\n", len(Words))

	// INDEX
	t = time.Now()
	SearchEngine := ferret.New(Words, Words, Values, ferret.UnicodeToLowerASCII)
	fmt.Println("Created index in:", time.Now().Sub(t))

	return SearchEngine
}

const (
	MaximumContextSize = 17
	ContextAfter       = 6
)

var (
	Words          []string      // (key) The 'true' value of the words. Used as a return value
	Values         []interface{} // (value) Some data mapped to the Words. Used for sorting, and as a return value
	CurrentContext []string
)

func eachWord(scannedWord string) {
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
