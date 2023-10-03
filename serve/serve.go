package serve

import (
	"encoding/json"
	"fmt"
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

        // fetch the GraphQL query from Github
        query := fetch(username)

        // Collect the commits counts for each language
        response := collect(query)

        jsonResponse, err := json.Marshal(response)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.Write(jsonResponse)
    })

    http.ListenAndServe(":8080", nil)
}
