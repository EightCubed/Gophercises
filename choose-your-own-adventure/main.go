package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"
)

type Story map[string]Chapter

type Chapter struct {
	Title   string    `json:"title"`
	Story   []string  `json:"story"`
	Options []Options `json:"options"`
}

type Options struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func main() {
	file_bytes, err := os.ReadFile("gopher.json")
	if err != nil {
		panic(err)
	}

	var Story Story
	err = json.Unmarshal(file_bytes, &Story)
	if err != nil {
		panic(err)
	}

	// data, err := printChapter(Story[arc])
	// fmt.Fprintf(w, "%v", data)

	var tmplFile = "story.tmpl"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, file_bytes)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("GET /{arc}", getStory)
	http.ListenAndServe(":8080", nil)
}

func getStory(w http.ResponseWriter, r *http.Request) {
	// Extract the "id" from the request URL path
	// arc := r.PathValue("arc")

	file_bytes, err := os.ReadFile("gopher.json")
	if err != nil {
		panic(err)
	}

	var Story Story
	err = json.Unmarshal(file_bytes, &Story)
	if err != nil {
		panic(err)
	}

	// data, err := printChapter(Story[arc])
	// fmt.Fprintf(w, "%v", data)

	var tmplFile = "story.tmpl"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, file_bytes)
	if err != nil {
		panic(err)
	}
}

func printChapter(chapter Chapter) (string, error) {
	fmt.Println("\t", chapter.Title)

	for _, paragraph := range chapter.Story {
		fmt.Println("\t", paragraph)
	}

	if len(chapter.Options) != 0 {
		fmt.Println("\nOptions : ")
		for idx, option := range chapter.Options {
			fmt.Println(idx+1, ")", option.Text)
		}
	} else {
		return "", nil
	}

	fmt.Print("\nEnter Choice : ")
	var selected_choice int
	_, err := fmt.Scan(&selected_choice)
	if err != nil {
		return "", err
	}

	return chapter.Options[selected_choice-1].Arc, nil
}
