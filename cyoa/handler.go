package cyoa

import (
	"html/template"
	"net/http"
)

type StoryHandler struct {
	Arcs map[string]Arc
}

func (s *StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	data, ok := s.Arcs[r.URL.Path[1:]]
	if !ok {
		http.Redirect(w, r, "/intro", http.StatusFound)
		return
	}
	tmpl := template.Must(template.ParseFiles("story.html.tmpl"))
	tmpl.Execute(w, data)
}
