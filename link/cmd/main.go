package main

import (
	"fmt"
	"github.com/mocak/gophercises/link/pkg/link"
	"os"
)

type Link struct {
	Href string
	Text string
}

func main() {

	file, err := os.Open("ex3.html")
	if err != nil {
		panic(err)
	}

	links, err := link.Parse(file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", links)
}
