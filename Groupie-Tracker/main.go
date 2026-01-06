package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type ArtistAPI struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type PageData struct {
	Artists []ArtistPage
}

type ArtistPage struct {
	Id        int
	Name      string
	Relations Relation
}

// 🔧 Fonction générique pour appeler l'API
func fetchAPI(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", home)

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func home(w http.ResponseWriter, r *http.Request) {

	// 🔹 Récupération artistes
	var artistsAPI []ArtistAPI
	err := fetchAPI("https://groupietrackers.herokuapp.com/api/artists", &artistsAPI)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 🔹 Récupération relations
	var relationsAPI struct {
		Index []Relation `json:"index"`
	}
	err = fetchAPI("https://groupietrackers.herokuapp.com/api/relations", &relationsAPI)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 🔹 Fusion artistes + relations
	var artistsPage []ArtistPage
	for _, artist := range artistsAPI {
		for _, rel := range relationsAPI.Index {
			if rel.Id == artist.Id {
				artistsPage = append(artistsPage, ArtistPage{
					Id:        artist.Id,
					Name:      artist.Name,
					Relations: rel,
				})
			}
		}
	}

	// 🔹 FuncMap pour passer Go → JS
	funcMap := template.FuncMap{
		"json": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		},
	}

	tmpl := template.Must(template.New("").Funcs(funcMap).ParseFiles(
		"templates/layout.html",
		"templates/index.html",
	))

	data := PageData{
		Artists: artistsPage,
	}

	tmpl.ExecuteTemplate(w, "layout", data)
}
