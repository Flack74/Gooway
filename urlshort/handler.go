package urlshort

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
)

// MapHandler returns an http.HandlerFunc that map paths to URLs.
// If path not found, fallback handler is called
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusMovedPermanently)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// BuildMapHandler populates the map from a slice of PathURL structs.
func BuildMapHandler(pathsToURL map[string]string, p []PathURL) {
	for _, pu := range p {
		pathsToURL[pu.Path] = pu.URL
	}
}

func YAMLHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return GenericHandler(data, yaml.Unmarshal, fallback)
}

func JSONHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return GenericHandler(data, json.Unmarshal, fallback)
}

// GenericHandler unmarshals data into PathURL slice and builds MapHandler
func GenericHandler(data []byte, unmarshalFunc func([]byte, any) error, fallback http.Handler) (http.HandlerFunc, error) {
	var p []PathURL
	err := unmarshalFunc(data, &p)
	if err != nil {
		return nil, err
	}

	pathsToURL := make(map[string]string)
	BuildMapHandler(pathsToURL, p)

	return MapHandler(pathsToURL, fallback), nil
}

// SqliteHandler reads paths and URLs from Sqlite DB and returns a handler.
// The SQLite table schema is assumed as:
// CREATE TABLE urls (path TEXT PRIMARY KEY, url TEXT NOT NULL);
func SqliteHandler(db *sql.DB, fallback http.Handler) (http.HandlerFunc, error) {
	var p []PathURL
	rows, err := db.Query("SELECT path, url FROM urls")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var pu PathURL
		err = rows.Scan(&pu.Path, &pu.URL)
		if err != nil {
			return nil, err
		}
		p = append(p, pu)
	}
	pathsToURL := make(map[string]string)
	BuildMapHandler(pathsToURL, p)

	return MapHandler(pathsToURL, fallback), err
}

// PathURL holds a path and corresponding URL.
type PathURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}
