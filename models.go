package main

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
)

type App struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	WebGame     bool   `json:"web_game"`
	ReleaseDate string `json:"release_date"`
	Cover       string   `json:"cover"`
	Screenshots []string `json:"screenshots"`
	Slug        string
}

type Category struct {
	Category string `json:"category"`
	Apps     []App  `json:"apps"`
}

type Data struct {
	Projects []Category `json:"projects"`
}

func load() Data {
	b, _ := os.ReadFile("projects.json")
	var d Data
	_ = json.Unmarshal(b, &d)
	for i := range d.Projects {
		for j := range d.Projects[i].Apps {
			a := &d.Projects[i].Apps[j]
			if strings.Contains(a.URL, "itch.io") {
				p := strings.Split(strings.Trim(a.URL, "/"), "/")
				a.Slug = p[len(p)-1]
			} else {
				a.Slug = strings.ToLower(strings.ReplaceAll(a.Name, " ", "-"))
			}
		}
		if d.Projects[i].Category == "games" {
			sort.Slice(d.Projects[i].Apps, func(x, y int) bool {
				dx, dy := d.Projects[i].Apps[x].ReleaseDate, d.Projects[i].Apps[y].ReleaseDate
				if dx == "" {
					return false
				}
				if dy == "" {
					return true
				}
				return dx > dy
			})
		}
	}
	return d
}
