package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var templates = template.Must(template.New("").Funcs(template.FuncMap{}).ParseGlob("web/templates/*.go*"))

func main() {
	r := mux.NewRouter()

	// Serve static files (CSS)
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := templates.ExecuteTemplate(w, "index.gohtmx", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	/* Example Handling of HTMX form submission

	r.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
	    // Process the form submission (add task logic here)
	    // For now, let's just print the task to the console
	    task := r.FormValue("task")
	    fmt.Println("New Task:", task)

	    http.Redirect(w, r, "/", http.StatusSeeOther)
	}).Methods("POST")

	*/

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
