package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"

	godotenv "github.com/joho/godotenv"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type LanguageCommit struct {
	Language string
	Commits  int
}

func main() {
    godotenv.Load(".env")

    accessToken := os.Getenv("TOKEN")

    src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
    httpClient := oauth2.NewClient(context.Background(), src)
    client := githubv4.NewClient(httpClient)

	targetUser := "Ruberald"

    type CommitTarget struct {
        History struct {
            TotalCount githubv4.Int
        }
    }

	var query struct {
		User struct {
			Repositories struct {
				Edges []struct {
					Node struct {
						Name           string
						PrimaryLanguage struct {
							Name string
						}
						DefaultBranchRef struct {
							Target struct {
                                CommitTarget `graphql:"... on Commit"`
							}
						}
					}
				}
			} `graphql:"repositories(first: 100)"`
		} `graphql:"user(login: $username)"`
	}

	variables := map[string]interface{}{
		"username": githubv4.String(targetUser),
	}

	err := client.Query(context.Background(), &query, variables)

	if err != nil {
		fmt.Println(err)
        os.Exit(1)
	}

    commitCounts := make(map[string]int)
	totalCommits := 0

	for _, edge := range query.User.Repositories.Edges {
		commitCount := edge.Node.DefaultBranchRef.Target.History.TotalCount
		language := edge.Node.PrimaryLanguage.Name

		// Skip repositories with an empty string for the programming language.
		if language != "" {
			// Update the commit count for the language.
			commitCounts[language] += int(commitCount)
			totalCommits += int(commitCount)
		}
	}

    // Collect language-commit pairs in a slice for sorting.
    var languageCommits []LanguageCommit
    for language, count := range commitCounts {
        languageCommits = append(languageCommits, LanguageCommit{Language: language, Commits: count})
    }

    // Sort the language-commit pairs by commit count in descending order.
    sort.SliceStable(languageCommits, func(i, j int) bool {
        return languageCommits[i].Commits > languageCommits[j].Commits
    })

    response := make(map[string]interface{})
    stats := []map[string]interface{}{}
    for _, lc := range languageCommits {
        percentage := float64(lc.Commits) / float64(totalCommits) * 100
        languageStats := map[string]interface{}{
            "language":        lc.Language,
            "commits_count":   lc.Commits,
            "commits_percent": fmt.Sprintf("%.2f%%", percentage),
        }
        stats = append(stats, languageStats)
    }
    response["stats"] = stats

    http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
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
