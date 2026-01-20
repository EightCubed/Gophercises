package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	urlshort "github.com/EightCubed/Gophercises/urlshort/pkg"
)

func main() {
	mux := defaultMux()

	yamlFile := flag.String("yaml", "", "A string flag with ")
	flag.Parse()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	if *yamlFile == "" {
		fmt.Println("No YAML file provided, skipping YAML handler")
		http.ListenAndServe(":8080", mapHandler)
		return
	}

	yaml, err := os.ReadFile(*yamlFile)
	if err != nil {
		panic(err)
	}

	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
