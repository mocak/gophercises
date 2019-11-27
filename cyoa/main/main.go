package main

import (
	"encoding/json"
	"github.com/mocak/gophercises/cyoa"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	arcs := make(map[string]cyoa.Arc)
	jsonBytes, err := ioutil.ReadFile("gopher.json")
	if err != nil {
		log.Fatal(err)
	}
	storyHandler := &cyoa.StoryHandler{Arcs: arcs}

	json.Unmarshal(jsonBytes, &arcs)
	http.Handle("/", storyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
