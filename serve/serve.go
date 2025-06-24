package serve

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type GitHubColor struct {
    Color string `json:"color"`
}

type ViewData struct {
    Stats  []Stat
    Colors map[string]string
}

func funcMap() template.FuncMap {
    return template.FuncMap{
        "add":      func(a, b float64) float64 { return a + b },
		"sub": func(a, b float64) float64 { return a - b },
        "mul":      func(a, b float64) float64 { return a * b },
        "deg2rad":  func(deg float64) float64 { return deg * math.Pi / 180 },
        "cos":      func(x float64) float64 { return math.Cos(x) },
        "sin":      func(x float64) float64 { return math.Sin(x) },
        "mod":      func(a, b int) int { return a % b },
        "div":      func(a, b int) int { return a / b },
        "cond":     func(b bool, t, f int) int { if b { return t }; return f },
        "default":  func(val, fallback string) string {
            if val == "" {
                return fallback
            }
            return val
        },
		"float64": func(x int) float64 { return float64(x) },
		"trim": strings.TrimSpace,
    }
}

func LoadColorMap(filename string) (map[string]string, error) {
    bytes, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    // Intermediate map to hold the full JSON structure
    var raw map[string]GitHubColor
    if err := json.Unmarshal(bytes, &raw); err != nil {
        return nil, err
    }

    // Convert to a simpler map[string]string
    colorMap := make(map[string]string)
    for lang, entry := range raw {
        colorMap[lang] = entry.Color
    }

    return colorMap, nil
}

func Serve() {
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")

		colors, err := LoadColorMap("colors.json")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}		

		fmt.Println("Go :", colors["Go"])

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

		// Load the HTML template from a file
		tmpl, err := template.New("card.svg").
			Funcs(funcMap()).      // Register functions here
			ParseFiles("serve/card.svg") // Or Parse(yourTemplateString)
		if err != nil {
			log.Fatal(err)
		}

		// Create a data struct for the template
		data := ViewData{
			Stats:  response, // assuming response has this field
			Colors: colors,
		}

		fmt.Println(colors[response[0].Language])

		// Execute the template with data
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
