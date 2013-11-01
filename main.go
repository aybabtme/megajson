package main

import (
	"flag"
	"log"

	"github.com/benbjohnson/megajson/generator"
)

func init() {
	log.SetFlags(0)
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}

	path := flag.Arg(0)
	if err := generator.Walk(path); err != nil {
		log.Fatalln(err)
	}
}

func usage() {
	log.Fatal("usage: megajson OPTIONS FILE")
}
