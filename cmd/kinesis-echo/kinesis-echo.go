package main

import (
	"io/ioutil"
	"log"
	"os"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	l := log.New(os.Stderr, "", log.LstdFlags)
	l.Printf("(%d) %q", len(b), b)
}
