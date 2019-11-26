package main

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/mocak/gophercises/urlshort"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	yamlFileName := flag.String("yaml", "data.yaml", "yaml file name")
	jsonFileName := flag.String("json", "data.json", "json file name")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlBytes, err := ioutil.ReadFile(*yamlFileName)
	if err != nil {
		log.Fatal(err)
	}

	yamlHandler, err := urlshort.YAMLHandler(yamlBytes, mapHandler)
	if err != nil {
		log.Fatal(err)
	}

	// Build the JSONHandler using the YAMLHandler as the
	// fallback
	jsonBytes, err := ioutil.ReadFile(*jsonFileName)
	if err != nil {
		log.Fatal(err)
	}

	jsonHandler, err := urlshort.JSONHandler(jsonBytes, yamlHandler)
	if err != nil {
		log.Fatal(err)
	}

	// Build the DBHandler using the mapHandler as the
	// fallback
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create Data
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("urls"))
		if err != nil {
			return err
		}
		if err := b.Put([]byte("/bolt"), []byte("https://github.com/boltdb/bolt")); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	dbHandler, err := urlshort.DBHandler(db, jsonHandler)

	fmt.Println("Starting the server on :8080")
	log.Fatal(http.ListenAndServe(":8080", dbHandler))
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
