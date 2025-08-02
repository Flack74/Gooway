package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"url-shortner/urlshort"

	_ "github.com/mattn/go-sqlite3" // SQLite driver import
)

func main() {
	// Define command line flags for YAML and JSON config file paths
	yamlFile := flag.String("yaml", "", "Path to a YAML file")
	jsonFile := flag.String("json", "", "Path to a JSON file")
	flag.Parse()

	// Warn if no config files are provided, will fallback to DB routes
	if *yamlFile == "" && *jsonFile == "" {
		fmt.Println("⚠️ No YAML or JSON file provided. Falling back to database routes from SQLite.")
	}

	// Create default HTTP request multiplexer with a basic hello handler
	mux := defaultMux()

	// Open SQLite database connection
	db, err := sql.Open("sqlite3", "data/urls.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Start with the default mux as the handler
	handler := http.Handler(mux)

	// Attempt to load SQLite routes and override handler if successful
	if sqliteHandler, err := urlshort.SqliteHandler(db, mux); err != nil {
		fmt.Println("⚠️ Could not load SQLite routes:", err)
	} else {
		handler = sqliteHandler
	}

	// If YAML file provided, read and wrap the handler with YAML routes
	if *yamlFile != "" {
		fmt.Println("ℹ️ YAML routes will override SQLite routes if there are conflicts.")
		yamldata, err := os.ReadFile(*yamlFile)
		if err != nil {
			fmt.Printf("Error reading YAML file: %v\n", err)
			return
		}
		handler, err = urlshort.YAMLHandler(yamldata, handler)
		if err != nil {
			panic(err)
		}
	}

	// If JSON file provided, read and wrap the handler with JSON routes
	if *jsonFile != "" {
		fmt.Println("ℹ️ JSON routes will override YAML and SQLite routes if there are conflicts.")

		jsondata, err := os.ReadFile(*jsonFile)
		if err != nil {
			fmt.Printf("Error reading JSON file: %v\n", err)
			return
		}
		handler, err = urlshort.JSONHandler(jsondata, handler)
		if err != nil {
			panic(err)
		}
	}

	// Start the HTTP server on port 8080 with the composed handler chain
	fmt.Println("Starting the server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// defaultMux sets up a simple HTTP server mux with a hello world handler
func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

// hello is the default fallback HTTP handler
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
