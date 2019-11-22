package urlshort

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {

		if url, ok := pathsToUrls[request.URL.Path]; ok {
			http.Redirect(writer, request, url, http.StatusFound)
			return
		}

		fallback.ServeHTTP(writer, request)
	}
}

type pathUrl struct {
	Path string
	Url  string
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	pathUrls, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(pathUrls)

	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yml []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl

	return pathUrls, yaml.Unmarshal(yml, &pathUrls)
}

func buildMap(pathUrls []pathUrl) map[string]string {
	res := make(map[string]string, len(pathUrls))
	for _, pu := range pathUrls {
		res[pu.Path] = pu.Url
	}

	return res
}

func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {

	pathUrls, err := parseJSON(jsonData)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(pathUrls)

	return MapHandler(pathMap, fallback), nil
}

func parseJSON(jsonData []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl

	return pathUrls, json.Unmarshal(jsonData, &pathUrls)
}
