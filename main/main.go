package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"url-shortner/urlshort"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	yamlFile := flag.String("yaml", "", "Path to a YAML file")
	jsonFile := flag.String("json", "", "Path to a JSON file")
	flag.Parse()

	if *yamlFile == "" && *jsonFile == "" {
		fmt.Println("⚠️ No YAML or JSON file provided. Falling back to database routes from SQLite.")
	}

	mux := defaultMux()

	db, err := sql.Open("sqlite3", "data/urls.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	handler, err := urlshort.SqliteHandler(db, mux)
	if err != nil {
		log.Fatal(err)
	}

	if *yamlFile != "" {
		yaml, err := os.ReadFile(*yamlFile)
		if err != nil {
			fmt.Printf("Error reading YAML file: %v\n", err)
			return
		}
		handler, err = urlshort.YAMLHandler([]byte(yaml), handler)
		if err != nil {
			panic(err)
		}
	}

	if *jsonFile != "" {
		jsondata, err := os.ReadFile(*jsonFile)
		if err != nil {
			fmt.Printf("Error reading YAML file: %v\n", err)
		}
		handler, err = urlshort.JSONHandler([]byte(jsondata), handler)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
