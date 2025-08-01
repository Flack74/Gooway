package urlshort

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
)

// import "gopkg.in/yaml.v3" for yaml

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var p []pathYmlURL
	err := yaml.Unmarshal(yml, &p)
	if err != nil {
		return nil, err
	}
	pathsToUrls := make(map[string]string)
	for _, pu := range p {
		pathsToUrls[pu.Path] = pu.URL
	}

	return MapHandler(pathsToUrls, fallback), err
}

type pathYmlURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func JSONHandler(jsondata []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pup []pathJsonURL
	err := json.Unmarshal(jsondata, &pup)
	if err != nil {
		return nil, err
	}
	pathsToURL := make(map[string]string)
	for _, j := range pup {
		pathsToURL[j.Path] = j.URL
	}
	return MapHandler(pathsToURL, fallback), err
}

type pathJsonURL struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

func SqliteHandler(db *sql.DB, fallback http.Handler) (http.HandlerFunc, error) {
	pathToUrls := make(map[string]string)
	rows, err := db.Query("SELECT path, url FROM urls")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var path string
		var url string
		err = rows.Scan(&path, &url)
		if err != nil {
			return nil, err
		}
		pathToUrls[path] = url
	}
	return MapHandler(pathToUrls, fallback), err
}
