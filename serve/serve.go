package serve

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

func Serve() {
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		urlString := r.URL.String()
		parsedURL, err := url.Parse(urlString)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return
		}

		username := parsedURL.Query().Get("username")

		// Fetch the GraphQL query from GitHub
		query := fetch(username)

		// Collect the commits counts for each language
		response := collect(query)

		// Marshal the JSON response
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Load the HTML template from a file
		tmpl, err := template.ParseFiles("serve/template.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create a JavaScript-safe version of the JSON data
		jsSafeData := template.JS(jsonResponse)

		// Create a data struct for the template
		data := struct {
			JSONData template.JS
		}{
			JSONData: jsSafeData,
		}

		// Execute the template with data
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
