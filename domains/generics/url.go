package generics

import (
	"fmt"
	"sort"
	"strings"
)

// URL is a custom implementation for representing an url as an object, because url.URL doesn't handle the job quite
// well.
type URL struct {
	// Host is the full host of the url, including the protocol.
	// For example, "https://example.com" (and not "example.com").
	Host string `json:"host" yaml:"host"`
	// Path is the path of the url, including the leading slash.
	Path string `json:"path" yaml:"path"`
	// Query is the object representation of query parameters.
	Query map[string][]string `json:"query" yaml:"query"`
}

// WithQuery adds data to the query object, without deleting old data. It creates a copy of the source URL.
func (url URL) WithQuery(data map[string]interface{}) URL {
	query := map[string][]string{}

	cp := URL{
		Host:  url.Host,
		Path:  url.Path,
		Query: query,
	}

	for k, v := range url.Query {
		cp.Query[k] = append(cp.Query[k], v...)
	}

	for k, v := range data {
		cp.Query[k] = append(cp.Query[k], fmt.Sprintf("%v", v))
	}

	return cp
}

// String returns the string representation of the url.
// Query parameters are sorted alphabetically in the final results, so it can be tested statically.
func (url URL) String() string {
	if url.Query == nil || len(url.Query) == 0 {
		return fmt.Sprintf("%s%s", url.Host, url.Path)
	}

	// Use an intermediary slice layer, rather than directly converting to string. This allows to sort and join easily.
	var queryStringList []string
	for k, v := range url.Query {
		queryStringList = append(queryStringList, fmt.Sprintf("%s=%s", k, strings.Join(v, ",")))
	}
	// Make order consistent, especially for unit testing.
	sort.Strings(queryStringList)

	return fmt.Sprintf("%s%s?%s", url.Host, url.Path, strings.Join(queryStringList, "&"))
}
