package main

import (
	f "github.com/argusdusty/Ferret"
)

const (
	MaximumContextSize = 17
	ContextAfter       = 6
	public             = "./public"
)

var (
	ferret         f.InvertedSuffix
	Words          []string
	Values         []interface{}
	CurrentContext []string
)

func main() {
	engine()
	routes()
	server()
}
