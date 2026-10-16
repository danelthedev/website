package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	data := load()
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	// also parse partials
	tmpl = template.Must(tmpl.ParseGlob("templates/partials/*.html"))

	render := func(w http.ResponseWriter, r *http.Request, page string, apps []App) {
		filter := r.URL.Query().Get("filter")
		if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("partial") == "games" {
			var f []App
			for _, g := range apps {
				if filter == "web" && !g.WebGame {
					continue
				}
				if filter == "dl" && g.WebGame {
					continue
				}
				f = append(f, g)
			}
			_ = tmpl.ExecuteTemplate(w, "games_grid", f)
			return
		}
		_ = tmpl.ExecuteTemplate(w, "layout", map[string]any{"Data": data, "Page": page, "Apps": apps})
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			var others []App
			for _, c := range data.Projects {
				if c.Category == "misc" || c.Category == "backgammon" {
					others = append(others, c.Apps...)
				}
			}
			for i, j := 0, len(others)-1; i < j; i, j = i+1, j-1 {
				others[i], others[j] = others[j], others[i]
			}
			_ = tmpl.ExecuteTemplate(w, "layout", map[string]any{"Data": data, "Page": "home", "Apps": others})
		case "/games":
			render(w, r, "games", data.Projects[1].Apps)
		case "/dailies":
			render(w, r, "dailies", data.Projects[0].Apps)
		case "/misc":
			render(w, r, "misc", data.Projects[2].Apps)
		case "/about":
			_ = tmpl.ExecuteTemplate(w, "layout", map[string]any{"Data": data, "Page": "about"})
		case "/backgammon":
			for _, c := range data.Projects {
				if c.Category == "backgammon" {
					render(w, r, "backgammon", c.Apps)
					return
				}
			}
			http.NotFound(w, r)
		default:
			if strings.HasPrefix(r.URL.Path, "/game/") {
				slug := strings.TrimPrefix(r.URL.Path, "/game/")
				for _, c := range data.Projects {
					for _, a := range c.Apps {
						if a.Slug == slug {
							_ = tmpl.ExecuteTemplate(w, "modal", a)
							return
						}
					}
				}
			}
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/game/", func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/game/")
		for _, c := range data.Projects {
			for _, a := range c.Apps {
				if a.Slug == slug {
					_ = tmpl.ExecuteTemplate(w, "modal", a)
					return
				}
			}
		}
		http.NotFound(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Println("listening :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
