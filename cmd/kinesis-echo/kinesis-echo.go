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
	f, err := os.OpenFile("out.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	l := log.New(f, "", log.LstdFlags)
	l.Printf("(%d) %q", len(b), b)
}
