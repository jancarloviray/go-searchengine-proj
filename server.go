package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	g "github.com/zenazn/goji"
	web "github.com/zenazn/goji/web"
)

func server() {
	g.Serve()
}

func routes() {
	g.Get("/", http.FileServer(http.Dir(public)))

	static := web.New()
	static.Get("/scripts/*", http.FileServer(http.Dir(public)))
	static.Get("/styles/*", http.FileServer(http.Dir(public)))
	static.Get("/img/*", http.FileServer(http.Dir(public)))

	api := web.New()
	api.Get("/api/search", searchHandler)

	g.Handle("/scripts/*", static)
	g.Handle("/styles/*", static)
	g.Handle("/img/*", static)
	g.Handle("/api/*", api)
}

type SearchResponse struct {
	//Query    string        `json:"query"`
	//Results  []string      `json:"results"`
	Duration string        `json:"duration"`
	Values   []interface{} `json:"values"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query()["s"][0]
	t := time.Now()
	_, values := ferret.Query(s, 10)
	duration := time.Now().Sub(t).String()
	data, _ := json.Marshal(SearchResponse{duration, values})
	fmt.Fprint(w, string(data))
}
